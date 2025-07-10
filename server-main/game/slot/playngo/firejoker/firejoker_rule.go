package firejoker

// See: https://freeslotshub.com/playngo/fire-joker/

import (
	_ "embed"

	"github.com/winnie192/slotgame/server/game/slot"
)

var BonusReel = []slot.Sym{1, 2, 3, 4, 5, 6, 7}

//go:embed firejoker_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [13][5]float64{
	{0, 0, 20, 50, 100}, // 1 seven
	{0, 0, 10, 25, 50},  // 2 bell
	{0, 0, 10, 25, 50},  // 3 melon
	{0, 0, 4, 10, 20},   // 4 plum
	{0, 0, 4, 10, 20},   // 5 orange
	{0, 0, 4, 10, 20},   // 6 lemon
	{0, 0, 4, 10, 20},   // 7 cherry
	{},                  // 8 bonus
	{},                  // 9 joker
}

// Scatters payment.
var ScatPay = [5]float64{0, 0.5, 3} // 8 bonus

// Scatter freespins table
var ScatFreespin = [5]int{0, 0, 10} // 8 bonus

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

const scat, jack = 8, 9

func (g *Game) Scanner(wins *slot.Wins) {
	g.ScanLined(wins)
	g.ScanScatters(wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
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

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.Scr.ScatNum(scat); count >= 2 {
		var pay, fs = ScatPay[count-1], ScatFreespin[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel) * pay,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   g.Scr.ScatPos(scat),
			Free: fs,
		})
	}
	if count := g.Scr.ScatNum(jack); count == 5 {
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel) * 100000,
			Mult: 1,
			Sym:  jack,
			Num:  5,
			XY:   g.Scr.ScatPos(scat),
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	if g.FSR == 0 {
		g.Scr.Spin(reels)
	} else {
		g.Scr.SpinBig(reels.Reel(1), BonusReel, reels.Reel(5))
	}
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
