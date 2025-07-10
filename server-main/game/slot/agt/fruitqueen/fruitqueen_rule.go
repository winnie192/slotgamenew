package fruitqueen

import (
	_ "embed"

	"github.com/winnie192/slotgame/server/game/slot"
)

//go:embed fruitqueen_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [8][5]float64{
	{},                    // 1 scatter
	{0, 0, 40, 400, 1200}, // 2 wild
	{0, 0, 20, 100, 400},  // 3 grapes
	{0, 0, 20, 40, 200},   // 4 strawberry
	{0, 0, 20, 40, 200},   // 5 plum
	{0, 0, 10, 24, 120},   // 6 pear
	{0, 0, 10, 20, 100},   // 7 orange
	{0, 0, 8, 16, 80},     // 8 cherry
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 5, 25, 500} // 1 scatter

// Bet lines
var BetLines = []slot.Linex{
	{1, 1, 1, 1, 1}, // 1
	{2, 2, 2, 2, 2}, // 2
	{3, 3, 3, 3, 3}, // 3
	{4, 4, 4, 4, 4}, // 4
	{5, 5, 5, 5, 5}, // 5
	{6, 6, 6, 6, 6}, // 6
	{1, 2, 3, 4, 5}, // 7
	{2, 3, 4, 5, 6}, // 8
	{5, 4, 3, 2, 1}, // 9
	{6, 5, 4, 3, 2}, // 10
	{1, 2, 3, 2, 1}, // 11
	{2, 3, 4, 3, 2}, // 12
	{3, 4, 5, 4, 3}, // 13
	{4, 5, 6, 5, 4}, // 14
	{3, 2, 1, 2, 3}, // 15
	{4, 3, 2, 3, 4}, // 16
	{5, 4, 3, 4, 5}, // 17
	{6, 5, 4, 5, 6}, // 18
}

type Game struct {
	slot.Slotx[slot.Screenx] `yaml:",inline"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame() *Game {
	var g = &Game{
		Slotx: slot.Slotx[slot.Screenx]{
			Sel: len(BetLines),
			Bet: 1,
		},
	}
	g.Scr.SetDim(5, 6)
	return g
}

func (g *Game) Clone() slot.SlotGame {
	var clone = *g
	return &clone
}

const wild, scat = 2, 1

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
		if numw >= 3 {
			payw = LinePay[wild-1][numw-1]
		}
		if numl >= 3 && syml > 0 {
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
