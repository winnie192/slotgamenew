//go:build !prod || full || novomatic

package beetlemania

import (
	"github.com/winnie192/slotgame/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Beetle Mania"},
		{Prov: "Novomatic", Name: "Beetle Mania Deluxe"},
		{Prov: "Novomatic", Name: "Hot Target"}, // see: https://www.slotsmate.com/software/novomatic/hot-target
	},
	GP: game.GPsel |
		game.GPfghas |
		game.GPfgreel |
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
	Info.SetupFactory(func() any { return NewGame() }, CalcStatReg)
}
