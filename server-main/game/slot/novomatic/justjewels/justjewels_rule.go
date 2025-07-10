package justjewels

// See: https://www.slotsmate.com/software/novomatic/just-jewels-deluxe

import (
	_ "embed"

	"github.com/winnie192/slotgame/server/game/slot"
)

//go:embed justjewels_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [8][5]float64{
	{0, 0, 50, 500, 5000}, // 1 crown
	{0, 0, 30, 150, 500},  // 2 gold
	{0, 0, 30, 150, 500},  // 3 money
	{0, 0, 15, 50, 200},   // 4 ruby
	{0, 0, 15, 50, 200},   // 5 sapphire
	{0, 0, 10, 25, 150},   // 6 emerald
	{0, 0, 10, 25, 150},   // 7 amethyst
	{},                    // 8 euro
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 2, 10, 50} // 8 euro

// Bet lines
var BetLines = slot.BetLinesNvm10

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

const scat = 8

func (g *Game) Scanner(wins *slot.Wins) {
	g.ScanLined(wins)
	g.ScanScatters(wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	for li := 1; li <= g.Sel; li++ {
		var line = BetLines[li-1]

		var numl slot.Pos = 1
		var syml = g.Scr.Pos(3, line)
		var xy slot.Linex
		xy.Set(3, line.At(3))
		if g.Scr.Pos(2, line) == syml {
			xy.Set(2, line.At(2))
			numl++
			if g.Scr.Pos(1, line) == syml {
				xy.Set(1, line.At(1))
				numl++
			}
		}
		if g.Scr.Pos(4, line) == syml {
			xy.Set(4, line.At(4))
			numl++
			if g.Scr.Pos(5, line) == syml {
				xy.Set(5, line.At(5))
				numl++
			}
		}

		if numl >= 3 {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * LinePay[syml-1][numl-1],
				Mult: 1,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   xy,
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
