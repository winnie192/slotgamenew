//go:build !prod || full || novomatic

package powerstars

import (
	"github.com/winnie192/slotgame/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Power Stars"},
	},
	GP: game.GPsel |
		game.GPrline |
		game.GPfgno |
		game.GPwild,
	SX:  5,
	SY:  3,
	SN:  len(LinePay),
	LN:  len(BetLines),
	BN:  0,
	RTP: game.MakeRtpList(ChanceMap),
}

func init() {
	Info.SetupFactory(func() any { return NewGame() }, CalcStat)
}
