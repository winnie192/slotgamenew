package luckyslot

// See: https://demo.agtsoftware.com/games/agt/luckyslot

import (
	_ "embed"

	"github.com/winnie192/slotgame/server/game/slot"
)

//go:embed luckyslot_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [10][5]float64{
	{},                     //  1 wild
	{},                     //  2 scatter
	{0, 10, 50, 250, 5000}, //  3 seven
	{0, 0, 40, 120, 800},   //  4 strawberry
	{0, 0, 40, 120, 600},   //  5 blueberry
	{0, 0, 20, 40, 200},    //  6 pear
	{0, 0, 12, 30, 160},    //  7 plum
	{0, 0, 12, 30, 160},    //  8 peach
	{0, 0, 8, 24, 120},     //  9 papaya
	{0, 0, 8, 24, 120},     // 10 cherry
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 5, 20, 100} // 2 scatter

// Bet lines
var BetLines = slot.BetLinesAgt5x3[:10]

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

const wild, scat = 1, 2

func (g *Game) Scanner(wins *slot.Wins) {
	g.ScanLined(wins)
	g.ScanScatters(wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	var reelwild [5]bool
	var x, y slot.Pos
	for x = 2; x <= 4; x++ {
		for y = 1; y <= 3; y++ {
			if g.Scr.At(x, y) == wild {
				reelwild[x-1] = true
				break
			}
		}
	}

	for li := 1; li <= g.Sel; li++ {
		var line = BetLines[li-1]

		var numl slot.Pos = 5
		var syml = g.Scr.Pos(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = g.Scr.Pos(x, line)
			if reelwild[x-1] {
				continue
			} else if sx != syml {
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

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.Scr.ScatNum(scat); count >= 3 {
		var pay = ScatPay[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel) * pay,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   g.Scr.ScatPos(scat),
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
