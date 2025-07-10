package powerstars

// See: https://freeslotshub.com/novomatic/power-stars/

import (
	_ "embed"
	"math/rand/v2"

	"github.com/winnie192/slotgame/server/game/slot"
)

//go:embed powerstars_reel.yaml
var reels []byte

var Reels = slot.ReadObj[*slot.Reels5x](reels)

//go:embed powerstars_chance.yaml
var chance []byte

var ChanceMap = slot.ReadMap[float64](chance)

// Lined payment.
var LinePay = [9][5]float64{
	{0, 0, 100, 500, 1000}, // 1 seven
	{0, 0, 50, 200, 500},   // 2 bell
	{0, 0, 20, 50, 200},    // 3 melon
	{0, 0, 20, 50, 200},    // 4 grapes
	{0, 0, 10, 30, 150},    // 5 plum
	{0, 0, 10, 30, 150},    // 6 orange
	{0, 0, 10, 20, 100},    // 7 lemon
	{0, 0, 10, 20, 100},    // 8 cherry
	{},                     // 9 star
}

// Bet lines
var BetLines = slot.BetLinesNvm10

type Game struct {
	slot.Slotx[slot.Screen5x3] `yaml:",inline"`
	// Pinned reel wild
	PRW [5]int `json:"prw" yaml:"prw" xml:"prw"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slotx: slot.Slotx[slot.Screen5x3]{
			Sel: len(BetLines),
			Bet: 1,
		},
	}
}

func (g *Game) Clone() slot.SlotGame {
	var clone = *g
	return &clone
}

const wild = 9

func (g *Game) Scanner(wins *slot.Wins) {
	g.ScanLined(wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	var reelwild [5]bool
	var fs int
	var x, y slot.Pos
	for x = 2; x <= 4; x++ {
		if g.PRW[x-1] > 0 {
			reelwild[x-1] = true
		} else {
			for y = 1; y <= 3; y++ {
				if g.Scr.At(x, y) == wild {
					reelwild[x-1] = true
					fs = 1
					break
				}
			}
		}
	}

	for li := 1; li <= g.Sel; li++ {
		var line = BetLines[li-1]
		var syml, symr slot.Sym
		var numl, numr slot.Pos
		var payl, payr float64

		syml, numl = g.Scr.Pos(1, line), 1
		for x = 2; x <= 5; x++ {
			var sx = g.Scr.Pos(x, line)
			if sx != syml && !reelwild[x-1] {
				break
			}
			numl++
		}
		payl = LinePay[syml-1][numl-1]

		if numl < 4 {
			symr, numr = g.Scr.Pos(5, line), 1
			for x = 4; x >= 2; x-- {
				var sx = g.Scr.Pos(x, line)
				if sx != symr && !reelwild[x-1] {
					break
				}
				numr++
			}
			payr = LinePay[symr-1][numr-1]
		}

		if payl > payr {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payl,
				Mult: 1,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   line.CopyL(numl),
			})
		} else if payr > 0 {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payr,
				Mult: 1,
				Sym:  symr,
				Num:  numr,
				Line: li,
				XY:   line.CopyL(numr),
			})
		}
		if fs > 0 {
			*wins = append(*wins, slot.WinItem{
				Sym:  wild,
				Free: fs,
			})
		}
	}
}

func (g *Game) Spin(mrtp float64) {
	g.Scr.Spin(Reels)
	if !g.FreeSpins() {
		var _, wc = slot.FindClosest(ChanceMap, mrtp) // wild chance
		var x, y slot.Pos
		for x = 2; x <= 4; x++ {
			if rand.Float64() < wc {
				y = rand.N[slot.Pos](3) + 1
				g.Scr.Set(x, y, wild)
			}
		}
	}
}

func (g *Game) Apply(wins slot.Wins) {
	if g.FSR != 0 {
		g.Gain += wins.Gain()
		g.FSN++
	} else {
		g.Gain = wins.Gain()
		g.FSN = 0
	}

	var x, y slot.Pos
	for x = 2; x <= 4; x++ {
		if g.PRW[x-1] > 0 {
			g.PRW[x-1]--
		} else {
			for y = 1; y <= 3; y++ {
				if g.Scr.At(x, y) == wild {
					g.PRW[x-1] = 1
					break
				}
			}
		}
	}
}

func (g *Game) FreeSpins() bool {
	return max(g.PRW[1], g.PRW[2], g.PRW[3]) > 0
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
