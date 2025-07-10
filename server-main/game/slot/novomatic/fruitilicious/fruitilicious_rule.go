package fruitilicious

// See: https://www.slotsmate.com/software/novomatic/fruitilicious

import (
	_ "embed"

	"github.com/winnie192/slotgame/server/game/slot"
)

//go:embed fruitilicious_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [7][5]float64{
	{0, 0, 100, 500, 5000}, // 1 seven
	{0, 0, 25, 100, 500},   // 2 melon
	{0, 0, 25, 100, 500},   // 3 grapes
	{0, 0, 10, 30, 125},    // 4 plum
	{0, 0, 10, 30, 125},    // 5 orange
	{0, 0, 10, 30, 125},    // 6 lemon
	{0, 0, 10, 30, 125},    // 7 cherry
}

// Bet lines
var BetLines = slot.BetLinesHot5

type Game struct {
	slot.Slotx[slot.Screen5x3] `yaml:",inline"`
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

func (g *Game) Scanner(wins *slot.Wins) {
	for li := 1; li <= g.Sel; li++ {
		var line = BetLines[li-1]
		var x slot.Pos

		var numl slot.Pos = 5
		var syml = g.Scr.Pos(1, line)
		for x = 2; x <= 5; x++ {
			var sx = g.Scr.Pos(x, line)
			if sx != syml {
				numl = x - 1
				break
			}
		}

		if numl >= 3 {
			var pay = LinePay[syml-1][numl-1]
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * pay,
				Mult: 1,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   line.CopyL(numl),
			})
			continue
		}

		var numr slot.Pos = 5
		var symr = g.Scr.Pos(5, line)
		for x = 4; x >= 1; x-- {
			var sx = g.Scr.Pos(x, line)
			if sx != symr {
				numr = 5 - x
				break
			}
		}

		if numr >= 3 {
			var pay = LinePay[symr-1][numr-1]
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * pay,
				Mult: 1,
				Sym:  symr,
				Num:  numr,
				Line: li,
				XY:   line.CopyR5(numr),
			})
			continue
		}
	}
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	g.Scr.Spin(reels)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
