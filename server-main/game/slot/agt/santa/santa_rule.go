package santa

// See: https://demo.agtsoftware.com/games/agt/santa

import (
	_ "embed"

	"github.com/winnie192/slotgame/server/game/slot"
)

//go:embed santa_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels4x](reels)

// Lined payment.
var LinePay = [10][5]float64{
	{},                 //  1 scatter
	{0, 20, 200, 1000}, //  2 wild
	{0, 10, 100, 500},  //  3 gnomes
	{0, 0, 50, 100},    //  4 snowman
	{0, 0, 40, 80},     //  5 christmas tree
	{0, 0, 30, 60},     //  6 socks
	{0, 0, 30, 60},     //  7 balls
	{0, 0, 20, 40},     //  8 sweets
	{0, 0, 10, 20},     //  9 present
	{0, 0, 10, 20},     // 10 bells
}

// Bet lines
var BetLines = slot.BetLinesAgt4x4[:10]

type Game struct {
	slot.Slotx[slot.Screen4x4] `yaml:",inline"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slotx: slot.Slotx[slot.Screen4x4]{
			Sel: len(BetLines),
			Bet: 1,
		},
	}
}

func (g *Game) Clone() slot.SlotGame {
	var clone = *g
	return &clone
}

const scat, wild = 1, 2

func (g *Game) Scanner(wins *slot.Wins) {
	g.ScanLined(wins)
	g.ScanScatters(wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	for li := 1; li <= g.Sel; li++ {
		var line = BetLines[li-1]

		var numw, numl slot.Pos = 0, 4
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 4; x++ {
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
	if count := g.Scr.ScatNum(scat); count > 0 {
		*wins = append(*wins, slot.WinItem{
			Sym:  scat,
			Num:  1,
			XY:   g.Scr.ScatPos(scat),
			Free: 3,
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
