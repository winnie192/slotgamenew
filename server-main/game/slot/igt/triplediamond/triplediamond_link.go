//go:build !prod || full || igt

package triplediamond

import (
	"github.com/winnie192/slotgame/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "IGT", Name: "Triple Diamond"}, // see: https://www.slotsmate.com/software/igt/triple-diamond
	},
	GP: game.GPsel |
		game.GPfgno |
		game.GPwild |
		game.GPwmult,
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
