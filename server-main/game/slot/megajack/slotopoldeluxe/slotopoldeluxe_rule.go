package slotgamedeluxe

import (
	_ "embed"

	"github.com/winnie192/slotgame/server/game/slot"
	"github.com/winnie192/slotgame/server/game/slot/megajack/slotgame"
)

//go:embed slotgamedeluxe_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

//go:embed slotgamedeluxe_jack.yaml
var jack []byte

var JackMap = slot.ReadMap[float64](jack)

// Lined payment.
var LinePay = [13][5]float64{
	{},                        //  1 dollar
	{0, 2, 5, 25, 100},        //  2 cherry
	{0, 2, 5, 25, 100},        //  3 plum
	{0, 0, 5, 25, 100},        //  4 wmelon
	{0, 0, 5, 25, 100},        //  5 grapes
	{0, 0, 10, 100, 250},      //  6 ananas
	{0, 0, 10, 100, 250},      //  7 lemon
	{0, 0, 10, 100, 250},      //  8 drink
	{0, 2, 10, 100, 500},      //  9 palm
	{0, 2, 10, 100, 500},      // 10 yacht
	{0, 10, 200, 2000, 10000}, // 11 eldorado
	{},                        // 12 spin
	{},                        // 13 dice
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 2, 20, 1000} // 1 dollar

const (
	mje1 = 1 // Eldorado9
	mje3 = 2 // Eldorado9
	mje6 = 3 // Eldorado9
	mje9 = 4 // Eldorado9
	mjm  = 5 // Monopoly
	mjc  = 6 // Champagne
	mjap = 7 // AztecPyramid
)

// Lined bonus games
var LineBonus = [13][5]int{
	{0, 0, 0, 0, 0},          //  1
	{0, 0, 0, 0, 0},          //  2
	{0, 0, 0, 0, 0},          //  3
	{0, 0, 0, 0, 0},          //  4
	{0, 0, 0, 0, 0},          //  5
	{0, 0, 0, 0, 0},          //  6
	{0, 0, 0, 0, 0},          //  7
	{0, 0, 0, 0, 0},          //  8
	{0, 0, 0, 0, 0},          //  9
	{0, 0, 0, 0, 0},          // 10
	{0, 0, 0, 0, 0},          // 11
	{0, 0, mje1, mje3, mje6}, // 12 Eldorado1, Eldorado3, Eldorado6
	{0, 0, 0, 0, mjm},        // 13 Monopoly
}

// Bet lines
var BetLines = slot.BetLinesMgj

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

// Not from lined paytable.
var special = [13]bool{
	true,  //  1
	false, //  2
	false, //  3
	false, //  4
	false, //  5
	false, //  6
	false, //  7
	false, //  8
	false, //  9
	false, // 10
	false, // 11
	true,  // 12
	true,  // 13
}

const (
	mjj        = 1     // jackpot ID
	wild, scat = 11, 1 // symbols
)

func (g *Game) Scanner(wins *slot.Wins) {
	g.ScanLined(wins)
	g.ScanScatters(wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	for li := 1; li <= g.Sel; li++ {
		var line = BetLines[li-1]

		var mw float64 = 1 // mult wild
		var numw, numl slot.Pos = 0, 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = g.Scr.Pos(x, line)
			if sx == wild {
				if syml == 0 {
					numw = x
				} else if special[syml-1] {
					numl = x - 1
					break
				}
				mw = 2
			} else if numw > 0 && special[sx-1] {
				numl = x - 1
				break
			} else if syml == 0 && sx != scat {
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
		if payl*mw > payw {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payl,
				Mult: mw,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   line.CopyL(numl),
			})
		} else if payw > 0 {
			var jid int
			if numw == 5 {
				jid = mjj
			}
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payw,
				Mult: 1,
				Sym:  wild,
				Num:  numw,
				Line: li,
				XY:   line.CopyL(numw),
				JID:  jid,
			})
		} else if syml > 0 && numl > 0 && LineBonus[syml-1][numl-1] > 0 {
			*wins = append(*wins, slot.WinItem{
				Mult: 1,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   line.CopyL(numl),
				BID:  LineBonus[syml-1][numl-1],
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

func (g *Game) Cost() (float64, bool) {
	return g.Bet * float64(g.Sel), true
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	g.Scr.Spin(reels)
}

func (g *Game) Spawn(wins slot.Wins, fund, mrtp float64) {
	for i, wi := range wins {
		switch wi.BID {
		case mje1:
			wins[i].Bon, wins[i].Pay = slotgame.EldoradoSpawn(g.Bet, 1)
		case mje3:
			wins[i].Bon, wins[i].Pay = slotgame.EldoradoSpawn(g.Bet, 3)
		case mje6:
			wins[i].Bon, wins[i].Pay = slotgame.EldoradoSpawn(g.Bet, 6)
		case mje9:
			wins[i].Bon, wins[i].Pay = slotgame.EldoradoSpawn(g.Bet, 9)
		case mjm:
			wins[i].Bon, wins[i].Pay = slotgame.MonopolySpawn(g.Bet)
		}
		if wi.JID != 0 {
			var bulk, _ = slot.FindClosest(JackMap, mrtp)
			var jf = bulk * g.Bet / slot.JackBasis
			if jf > 1 {
				jf = 1
			}
			wins[i].Jack = jf * fund
		}
	}
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
