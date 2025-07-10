//go:build !prod || full || novomatic

package ultrasevens

import (
	"github.com/winnie192/slotgame/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Ultra Sevens"}, // see: https://www.slotsmate.com/software/novomatic/ultra-sevens
	},
	GP: game.GPfgno |
		game.GPjack |
		game.GPscat,
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
