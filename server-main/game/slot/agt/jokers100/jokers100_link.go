//go:build !prod || full || agt

package jokers100

import (
	"github.com/winnie192/slotgame/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "100 Jokers"},
		{Prov: "AGT", Name: "50 Happy Santa"}, // see: https://demo.agtsoftware.com/games/agt/happysanta50
		{Prov: "AGT", Name: "40 Bigfoot"},     // see: https://demo.agtsoftware.com/games/agt/bigfoot40
	},
	GP: game.GPsel |
		game.GPfgno |
		game.GPscat |
		game.GPwild,
	SX:  5,
	SY:  4,
	SN:  len(LinePay),
	LN:  len(BetLines),
	BN:  0,
	RTP: game.MakeRtpList(ReelsMap),
}

func init() {
	Info.SetupFactory(func() any { return NewGame() }, CalcStat)
}
