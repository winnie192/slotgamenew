//go:build !prod || full || playtech

package greatblue

import (
	"github.com/winnie192/slotgame/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "Playtech", Name: "Great Blue"}, // see: https://freeslotshub.com/playtech/great-blue/
		{Prov: "Playtech", Name: "Irish Luck"}, // see: https://freeslotshub.com/playtech/irish-luck/
	},
	GP: game.GPsel |
		game.GPretrig |
		game.GPfgmult |
		game.GPscat |
		game.GPwild |
		game.GPwmult,
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
