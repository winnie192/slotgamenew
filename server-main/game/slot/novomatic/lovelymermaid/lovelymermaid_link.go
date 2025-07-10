//go:build !prod || full || novomatic

package lovelymermaid

import (
	"github.com/winnie192/slotgame/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Lovely Mermaid"}, // see: https://www.slotsmate.com/software/novomatic/lovely-mermaid
	},
	GP: game.GPsel |
		game.GPretrig |
		game.GPjack |
		game.GPscat |
		game.GPwild,
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
