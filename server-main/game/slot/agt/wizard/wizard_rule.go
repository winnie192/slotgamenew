package wizard

// See: https://demo.agtsoftware.com/games/agt/wizard

import (
	_ "embed"

	"github.com/winnie192/slotgame/server/game/slot"
)

//go:embed wizard_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [11][5]float64{
	{0, 10, 50, 250, 1000}, //  1 wild
	{},                     //  2 scatter
	{0, 5, 10, 20, 50},     //  3 owl
	{0, 4, 10, 20, 45},     //  4 cat
	{0, 0, 8, 15, 40},      //  5 cauldron
	{0, 0, 8, 15, 30},      //  6 emerald
	{0, 0, 6, 15, 25},      //  7 ruby
	{0, 0, 4, 10, 20},      //  8 ace
	{0, 0, 4, 10, 20},      //  9 king
	{0, 0, 2, 8, 15},       // 10 queen
	{0, 0, 2, 8, 15},       // 11 jack
}

// Scatter freespins table
var ScatFreespin = [12]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 15, 30} // 2 scatter

// Bet lines
var BetLines = slot.BetLinesAgt5x4[:50]

type Game struct {
	slot.Slotx[slot.Screen5x4] `yaml:",inline"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slotx: slot.Slotx[slot.Screen5x4]{
			Sel: len(BetLines),
			Bet: 1,
		},
	}
}

func (g *Game) Clone() slot.SlotGame {
	var clone = *g
	return &clone
}

const wild, scat = 1, 2

func (g *Game) Scanner(wins *slot.Wins) {
	g.ScanLined(wins)
	g.ScanScatters(wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	for li := 1; li <= g.Sel; li++ {
		var line = BetLines[li-1]

		var numw, numl slot.Pos = 0, 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = g.Scr.Pos(x, line)
			if sx == wild {
				if syml == 0 {
					numw = x
				}
			} else if syml == 0 {
				syml = sx
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

		var payw, payl float64
		if numw >= 2 {
			payw = LinePay[wild-1][numw-1]
		}
		if numl >= 2 && syml > 0 {
			payl = LinePay[syml-1][numl-1]
		}
		if payl > payw {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payl,
				Mult: 1,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   line.CopyL(numl),
			})
		} else if payw > 0 {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payw,
				Mult: 1,
				Sym:  wild,
				Num:  numw,
				Line: li,
				XY:   line.CopyL(numw),
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.Scr.SymNum(scat); count >= 10 {
		const pay float64 = 2
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel) * pay,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			Free: ScatFreespin[count-1],
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	g.Scr.Spin(reels)
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
