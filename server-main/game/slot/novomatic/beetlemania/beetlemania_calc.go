package beetlemania

import (
	"context"
	"fmt"
	"io"
	"math"
	"time"

	"github.com/winnie192/slotgame/server/game/slot"
)

// Attention! On freespins can be calculated median only, not expectation.

func CalcStatBon(ctx context.Context) float64 {
	var reels = ReelsBon
	var g = NewGame()
	g.Sel = 1
	g.FSR = 10 // set free spins mode
	var s slot.Stat

	var calc = func(w io.Writer) float64 {
		var reshuf = float64(s.Count())
		var lrtp = s.LineRTP(g.Sel)
		var qjazz = float64(s.BonusCount(jbonus)) / reshuf
		var jpow = math.Pow(2, 10*qjazz) // jazz power
		var rtpjazz = lrtp*jpow - lrtp
		var rtp = lrtp * jpow
		fmt.Fprintf(w, "symbols: %.5g(lined) + 0(scatter) = %.6f%%\n", lrtp, lrtp)
		fmt.Fprintf(w, "jazzbee bonuses: frequency 1/%.5g, pow = %.5g, rtp = %.6f%%\n", reshuf/float64(s.BonusCount(jbonus)), jpow, rtpjazz)
		fmt.Fprintf(w, "RTP = rtp(sym) + rtp(jazz) = %.5g + %.5g = %.6f%%\n", lrtp, rtpjazz, rtp)
		return rtp
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc,
		time.Tick(2*time.Second), time.Tick(2*time.Second))
}

func CalcStatReg(ctx context.Context, mrtp float64) float64 {
	fmt.Printf("*bonus reels calculations*\n")
	var rtpfs = CalcStatBon(ctx)
	if ctx.Err() != nil {
		return 0
	}
	fmt.Printf("*regular reels calculations*\n")
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	var g = NewGame()
	g.Sel = 1
	var s slot.Stat

	var calc = func(w io.Writer) float64 {
		var reshuf = float64(s.Count())
		var lrtp, srtp = s.LineRTP(g.Sel), s.ScatRTP(g.Sel)
		var rtpsym = lrtp + srtp
		var q = float64(s.FreeCount()) / reshuf
		var rtp = rtpsym + q*rtpfs
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Fprintf(w, "free spins %d, q = %.5g\n", s.FreeCount(), q)
		fmt.Fprintf(w, "free games frequency: 1/%.5g\n", reshuf/float64(s.FreeHits()))
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%\n", rtpsym, q, rtpfs, rtp)
		return rtp
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc,
		time.Tick(2*time.Second), time.Tick(2*time.Second))
}
