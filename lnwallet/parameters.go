package lnwallet

import (
	"github.com/Katano-Sukune/xpcutil"
	"github.com/Katano-Sukune/xpcwallet/wallet/txrules"
	"github.com/lightningnetwork/lnd/input"
)

// DefaultDustLimit is used to calculate the dust HTLC amount which will be
// send to other node during funding process.
func DefaultDustLimit() xpcutil.Amount {
	return txrules.GetDustThreshold(input.P2WSHSize, txrules.DefaultRelayFeePerKb)
}
