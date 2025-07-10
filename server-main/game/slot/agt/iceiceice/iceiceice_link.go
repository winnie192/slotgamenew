//go:build !prod || full || agt

package iceiceice

import (
	"github.com/winnie192/slotgame/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Ice Ice Ice"},
		{Prov: "AGT", Name: "5 Hot Hot Hot"}, // see: https://demo.agtsoftware.com/games/agt/hothothot5
	},
	GP: game.GPretrig |
		game.GPscat |
		game.GPwild,
	SX:  3,
	SY:  3,
	SN:  len(LinePay),
	LN:  len(BetLines),
	BN:  0,
	RTP: game.MakeRtpList(ReelsMap),
}

func init() {
	Info.SetupFactory(func() any { return NewGame() }, CalcStatReg)
}
