//go:build !prod || full || igt

package richgirl

import (
	"github.com/winnie192/slotgame/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "IGT", Name: "Rich Girl"}, // see: https://www.slotsmate.com/software/igt/rich-girl
	},
	GP: game.GPsel |
		game.GPretrig |
		game.GPfgreel |
		game.GPscat |
		game.GPwild |
		game.GPwmult,
	SX:  5,
	SY:  3,
	SN:  len(LinePayReg),
	LN:  len(BetLines),
	BN:  0,
	RTP: game.MakeRtpList(ReelsMap),
}

func init() {
	Info.SetupFactory(func() any { return NewGame() }, CalcStatReg)
}
