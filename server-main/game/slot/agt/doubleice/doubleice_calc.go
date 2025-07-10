package doubleice

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/winnie192/slotgame/server/game/slot"
)

func CalcStat(ctx context.Context, mrtp float64) float64 {
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	var g = NewGame()
	g.Sel = 1
	var s slot.Stat

	var calc = func(w io.Writer) float64 {
		var lrtp = s.LineRTP(g.Sel)
		fmt.Fprintf(w, "RTP = %.6f%%\n", lrtp)
		return lrtp
	}

	return slot.ScanReels3x(ctx, &s, g, reels, calc,
		time.Tick(2*time.Second), time.Tick(2*time.Second))
}
