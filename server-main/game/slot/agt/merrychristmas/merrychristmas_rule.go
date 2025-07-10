package merrychristmas

// See: https://demo.agtsoftware.com/games/agt/christmas

import (
	_ "embed"

	"github.com/winnie192/slotgame/server/game/slot"
)

//go:embed merrychristmas_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels3x](reels)

// Lined payment.
var LinePay = [8]float64{
	500, // 1 snowman
	0,   // 2 scatter
	250, // 3 ice
	100, // 4 sled
	20,  // 5 house
	10,  // 6 bell
	5,   // 7 deer
}

// Bet lines
var BetLines = slot.BetLinesHot3

type Game struct {
	slot.Slotx[slot.Screen3x3] `yaml:",inline"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slotx: slot.Slotx[slot.Screen3x3]{
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
	for li := 1; li <= g.Sel; li++ {
		var line = BetLines[li-1]

		var numw, numl slot.Pos = 0, 3
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 3; x++ {
			var sx = g.Scr.Pos(x, line)
			if g.FSR > 0 && sx == wild {
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

		if numw == 3 {
			var pay = LinePay[wild-1]
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * pay,
				Mult: 1,
				Sym:  wild,
				Num:  numw,
				Line: li,
				XY:   line.CopyL(numw),
			})
		} else if numl == 3 {
			var pay = LinePay[syml-1]
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

	if count := g.Scr.ScatNum(scat); count >= 3 {
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel) * 10,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   g.Scr.ScatPos(scat),
			Free: 20,
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	g.Scr.Spin(reels)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
