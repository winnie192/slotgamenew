//go:build !prod || full || agt

package aislot

import (
	"github.com/winnie192/slotgame/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "AI"},          // see: https://demo.agtsoftware.com/games/agt/aislot
		{Prov: "AGT", Name: "Tesla"},       // see: https://demo.agtsoftware.com/games/agt/tesla
		{Prov: "AGT", Name: "Book of Set"}, // see: https://demo.agtsoftware.com/games/agt/bookofset
		{Prov: "AGT", Name: "Pharaoh II"},  // see: https://demo.agtsoftware.com/games/agt/pharaoh2
	},
	GP: game.GPsel |
		game.GPretrig |
		game.GPscat |
		game.GPwild,
	SX:  5,
	SY:  3,
	SN:  len(LinePay),
	LN:  len(BetLines),
	BN:  0,
	RTP: game.MakeRtpList(ReelsMap),
}

func init() {
	Info.SetupFactory(func() any { return NewGame() }, CalcStat)
}
