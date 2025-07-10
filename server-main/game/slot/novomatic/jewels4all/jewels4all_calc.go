package jewels4all

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/winnie192/slotgame/server/game/slot"
)

func BruteForceEuro(ctx context.Context, s slot.Stater, g *Game, reels slot.Reels, x, y slot.Pos) {
	var screen = &g.Scr
	var wins slot.Wins
	var r1 = reels.Reel(1)
	var r2 = reels.Reel(2)
	var r3 = reels.Reel(3)
	var r4 = reels.Reel(4)
	var r5 = reels.Reel(5)
	for i1 := range r1 {
		screen.SetCol(1, r1, i1)
		for i2 := range r2 {
			screen.SetCol(2, r2, i2)
			for i3 := range r3 {
				screen.SetCol(3, r3, i3)
				for i4 := range r4 {
					screen.SetCol(4, r4, i4)
					for i5 := range r5 {
						screen.SetCol(5, r5, i5)
						var sym slot.Sym
						if x > 0 {
							sym = screen.At(x, y)
							screen.Set(x, y, wild)
						}
						g.Scanner(&wins)
						if x > 0 {
							screen.Set(x, y, sym)
						}
						s.Update(wins)
						wins.Reset()
						if s.Count()&100 == 0 {
							select {
							case <-ctx.Done():
								return
							default:
							}
						}
					}
				}
			}
		}
	}
}

func CalcStatEuro(ctx context.Context, x, y slot.Pos) float64 {
	var reels = Reels
	var g = NewGame()
	g.Sel = 1
	var s slot.Stat

	fmt.Printf("calculations of euro at [%d,%d]\n", x, y)

	var calc = func(w io.Writer) float64 {
		var lrtp = s.LineRTP(g.Sel)
		fmt.Fprintf(w, "RTP[%d,%d] = %.6f%%\n", x, y, lrtp)
		return lrtp
	}

	func() time.Duration {
		var t0 = time.Now()
		var ctx2, cancel2 = context.WithCancel(ctx)
		defer cancel2()
		s.SetPlan(reels.Reshuffles())
		go slot.Progress(ctx2, &s, time.Tick(2*time.Second), calc)
		BruteForceEuro(ctx2, &s, g, reels, x, y)
		return time.Since(t0)
	}()
	return calc(os.Stdout)
}

func CalcStat(ctx context.Context, mrtp float64) (rtp float64) {
	var wc, _ = slot.FindClosest(ChanceMap, mrtp) // wild chance

	var b = 1 / wc
	fmt.Printf("wild chance %.5g, b = %.5g\n", wc, b)
	var rtp00 = CalcStatEuro(ctx, 0, 0)
	var rtpeu float64
	var x, y slot.Pos
	for x = 1; x <= 5; x++ {
		for y = 1; y <= 3; y++ {
			rtpeu += CalcStatEuro(ctx, x, y)
		}
	}
	rtpeu /= 15
	rtp = rtp00 + wc*rtpeu
	fmt.Printf("euro avr: rtpeu = %.6f%%\n", rtpeu)
	fmt.Printf("RTP = %.5g(sym) + wc*%.5g(eu) = %.6f%%\n", rtp00, rtpeu, rtp)
	return
}
