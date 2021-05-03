package active

import (
	"net"
	"time"

	"github.com/p9c/p9/pkg/amt"
	"github.com/p9c/p9/pkg/btcaddr"
	"github.com/p9c/p9/pkg/connmgr"

	"github.com/p9c/p9/pkg/chaincfg"
)

// Config stores current state of the node
type Config struct {
	Lookup              connmgr.LookupFunc
	Oniondial           func(string, string, time.Duration) (net.Conn, error)
	Dial                func(string, string, time.Duration) (net.Conn, error)
	AddedCheckpoints    []chaincfg.Checkpoint
	ActiveMiningAddrs   []btcaddr.Address
	ActiveMinerKey      []byte
	ActiveMinRelayTxFee amt.Amount
	ActiveWhitelists    []*net.IPNet
	DropAddrIndex       bool
	DropTxIndex         bool
	DropCfIndex         bool
	Save                bool
}
