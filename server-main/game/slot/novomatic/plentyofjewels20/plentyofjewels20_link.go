//go:build !prod || full || novomatic

package plentyofjewels20

import (
	"github.com/winnie192/slotgame/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Plenty of Jewels 20 hot"}, // see: https://www.slotsmate.com/software/novomatic/plenty-of-jewels-20-hot
		{Prov: "Novomatic", Name: "Plenty of Fruit 20 hot"},  // see: https://www.slotsmate.com/software/novomatic/plenty-of-fruit-20-hot
	},
	GP: game.GPsel |
		game.GPfgno |
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
