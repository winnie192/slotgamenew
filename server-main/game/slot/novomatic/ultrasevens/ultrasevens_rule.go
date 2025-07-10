package ultrasevens

// See: https://www.slotsmate.com/software/novomatic/ultra-sevens

import (
	_ "embed"

	"github.com/winnie192/slotgame/server/game/slot"
)

//go:embed ultrasevens_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

//go:embed ultrasevens_jack.yaml
var jack []byte

var JackMap = slot.ReadMap[[3]float64](jack)

// Lined payment.
var LinePay = [7][5]float64{
	{0, 10, 100, 1000, 10000}, // 1 seven
	{0, 0, 40, 200, 500},      // 2 melon
	{0, 0, 40, 200, 500},      // 3 grapes
	{0, 0, 10, 50, 200},       // 4 plum
	{0, 0, 10, 50, 200},       // 5 orange
	{0, 0, 10, 50, 200},       // 6 lemon
	{0, 5, 10, 50, 200},       // 7 cherry
}

// Bet lines
var BetLines = slot.BetLinesNvm5x4[:]

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

const (
	ssj1 = 1
	ssj2 = 2
	ssj3 = 3
)

func Filled(screen *slot.Screen5x4) slot.Sym {
	var sym = screen[4][3]
	for x := range 5 {
		for y := range 4 {
			if screen[x][y] != sym {
				return 0
			}
		}
	}
	return sym
}

func (g *Game) Scanner(wins *slot.Wins) {
	switch sym := Filled(&g.Scr); sym {
	case 1:
		*wins = append(*wins, slot.WinItem{
			Sym: sym,
			JID: ssj1,
		})
		return
	case 2, 3:
		*wins = append(*wins, slot.WinItem{
			Sym: sym,
			JID: ssj2,
		})
		return
	case 4, 5, 6, 7:
		*wins = append(*wins, slot.WinItem{
			Sym: sym,
			JID: ssj3,
		})
		return
	}
	for li := 1; li <= g.Sel; li++ {
		var line = BetLines[li-1]

		var numl slot.Pos = 5
		var syml = g.Scr.Pos(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = g.Scr.Pos(x, line)
			if sx != syml {
				numl = x - 1
				break
			}
		}

		if pay := LinePay[syml-1][numl-1]; pay > 0 {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * pay,
				Mult: 1,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   line.CopyL(numl),
			})
		}
	}
}

func (g *Game) Cost() (float64, bool) {
	return g.Bet * float64(g.Sel), true
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	g.Scr.Spin(reels)
}

func (g *Game) Spawn(wins slot.Wins, fund, mrtp float64) {
	for i, wi := range wins {
		if wi.JID != 0 {
			var ji = wi.JID - 1
			var bulk, _ = slot.FindClosest(JackMap, mrtp)
			var jf = bulk[ji] * g.Bet / slot.JackBasis
			if jf > 1 {
				jf = 1
			}
			wins[i].Jack = jf * fund
		}
	}
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
