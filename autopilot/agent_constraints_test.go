package autopilot

import (
	"testing"
	"time"

	prand "math/rand"

	"github.com/Katano-Sukune/xpcutil"
	"github.com/lightningnetwork/lnd/lnwire"
)

func TestConstraintsChannelBudget(t *testing.T) {
	t.Parallel()

	prand.Seed(time.Now().Unix())

	const (
		minChanSize = 0
		maxChanSize = xpcutil.Amount(xpcutil.SatoshiPerBitcoin)

		chanLimit = 3

		threshold = 0.5
	)

	constraints := NewConstraints(
		minChanSize,
		maxChanSize,
		chanLimit,
		0,
		threshold,
	)

	randChanID := func() lnwire.ShortChannelID {
		return lnwire.NewShortChanIDFromInt(uint64(prand.Int63()))
	}

	testCases := []struct {
		channels  []Channel
		walletAmt xpcutil.Amount

		needMore     bool
		amtAvailable xpcutil.Amount
		numMore      uint32
	}{
		// Many available funds, but already have too many active open
		// channels.
		{
			[]Channel{
				{
					ChanID:   randChanID(),
					Capacity: xpcutil.Amount(prand.Int31()),
				},
				{
					ChanID:   randChanID(),
					Capacity: xpcutil.Amount(prand.Int31()),
				},
				{
					ChanID:   randChanID(),
					Capacity: xpcutil.Amount(prand.Int31()),
				},
			},
			xpcutil.Amount(xpcutil.SatoshiPerBitcoin * 10),
			false,
			0,
			0,
		},

		// Ratio of funds in channels and total funds meets the
		// threshold.
		{
			[]Channel{
				{
					ChanID:   randChanID(),
					Capacity: xpcutil.Amount(xpcutil.SatoshiPerBitcoin),
				},
				{
					ChanID:   randChanID(),
					Capacity: xpcutil.Amount(xpcutil.SatoshiPerBitcoin),
				},
			},
			xpcutil.Amount(xpcutil.SatoshiPerBitcoin * 2),
			false,
			0,
			0,
		},

		// Ratio of funds in channels and total funds is below the
		// threshold. We have 10 BTC allocated amongst channels and
		// funds, atm. We're targeting 50%, so 5 BTC should be
		// allocated. Only 1 BTC is atm, so 4 BTC should be
		// recommended. We should also request 2 more channels as the
		// limit is 3.
		{
			[]Channel{
				{
					ChanID:   randChanID(),
					Capacity: xpcutil.Amount(xpcutil.SatoshiPerBitcoin),
				},
			},
			xpcutil.Amount(xpcutil.SatoshiPerBitcoin * 9),
			true,
			xpcutil.Amount(xpcutil.SatoshiPerBitcoin * 4),
			2,
		},

		// Ratio of funds in channels and total funds is below the
		// threshold. We have 14 BTC total amongst the wallet's
		// balance, and our currently opened channels. Since we're
		// targeting a 50% allocation, we should commit 7 BTC. The
		// current channels commit 4 BTC, so we should expected 3 BTC
		// to be committed. We should only request a single additional
		// channel as the limit is 3.
		{
			[]Channel{
				{
					ChanID:   randChanID(),
					Capacity: xpcutil.Amount(xpcutil.SatoshiPerBitcoin),
				},
				{
					ChanID:   randChanID(),
					Capacity: xpcutil.Amount(xpcutil.SatoshiPerBitcoin * 3),
				},
			},
			xpcutil.Amount(xpcutil.SatoshiPerBitcoin * 10),
			true,
			xpcutil.Amount(xpcutil.SatoshiPerBitcoin * 3),
			1,
		},

		// Ratio of funds in channels and total funds is above the
		// threshold.
		{
			[]Channel{
				{
					ChanID:   randChanID(),
					Capacity: xpcutil.Amount(xpcutil.SatoshiPerBitcoin),
				},
				{
					ChanID:   randChanID(),
					Capacity: xpcutil.Amount(xpcutil.SatoshiPerBitcoin),
				},
			},
			xpcutil.Amount(xpcutil.SatoshiPerBitcoin),
			false,
			0,
			0,
		},
	}

	for i, testCase := range testCases {
		amtToAllocate, numMore := constraints.ChannelBudget(
			testCase.channels, testCase.walletAmt,
		)

		if amtToAllocate != testCase.amtAvailable {
			t.Fatalf("test #%v: expected %v, got %v",
				i, testCase.amtAvailable, amtToAllocate)
		}
		if numMore != testCase.numMore {
			t.Fatalf("test #%v: expected %v, got %v",
				i, testCase.numMore, numMore)
		}
	}
}
