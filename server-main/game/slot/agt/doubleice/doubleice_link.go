//go:build !prod || full || agt

package doubleice

import (
	"github.com/winnie192/slotgame/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Double Ice"},
		{Prov: "AGT", Name: "Double Hot"}, // see: https://demo.agtsoftware.com/games/agt/double
	},
	GP: game.GPfgno |
		game.GPfill,
	SX:  3,
	SY:  3,
	SN:  len(LinePay),
	LN:  len(BetLines),
	BN:  0,
	RTP: game.MakeRtpList(ReelsMap),
}

func init() {
	Info.SetupFactory(func() any { return NewGame() }, CalcStat)
}
