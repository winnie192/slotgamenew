//go:build !prod || full || novomatic

package secretelixir

import (
	"github.com/winnie192/slotgame/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Secret Elixir"}, // see: https://www.slotsmate.com/software/novomatic/secret-elixir
	},
	GP: game.GPsel |
		game.GPretrig |
		game.GPfgreel |
		game.GPfgmult |
		game.GPrmult |
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
