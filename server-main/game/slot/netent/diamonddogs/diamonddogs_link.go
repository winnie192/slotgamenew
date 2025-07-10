//go:build !prod || full || netent

package diamonddogs

import (
	"github.com/winnie192/slotgame/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "NetEnt", Name: "Diamond Dogs"},
		{Prov: "NetEnt", Name: "Voodoo Vibes"},
	},
	GP: game.GPsel |
		game.GPretrig |
		game.GPfgreel |
		game.GPfgmult |
		game.GPscat |
		game.GPwild,
	SX:  5,
	SY:  3,
	SN:  len(LinePay),
	LN:  len(BetLines),
	BN:  1,
	RTP: game.MakeRtpList(ReelsMap),
}

func init() {
	Info.SetupFactory(func() any { return NewGame() }, CalcStatReg)
}
