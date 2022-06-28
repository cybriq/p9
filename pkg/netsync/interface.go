package netsync

import (
	"github.com/cybriq/p9/pkg/blockchain"
	"github.com/cybriq/p9/pkg/chaincfg"
	"github.com/cybriq/p9/pkg/chainhash"
	"github.com/cybriq/p9/pkg/mempool"
	"github.com/cybriq/p9/pkg/peer"
	"github.com/cybriq/p9/pkg/util"
	"github.com/cybriq/p9/pkg/wire"
)

// PeerNotifier exposes methods to notify peers of status changes to transactions, blocks, etc. Currently server (in the
// main package) implements this interface.
type PeerNotifier interface {
	AnnounceNewTransactions(newTxs []*mempool.TxDesc)
	UpdatePeerHeights(
		latestBlkHash *chainhash.Hash, latestHeight int32,
		updateSource *peer.Peer,
	)
	RelayInventory(invVect *wire.InvVect, data interface{})
	TransactionConfirmed(tx *util.Tx)
}

// Config is a configuration struct used to initialize a new SyncManager.
type Config struct {
	PeerNotifier       PeerNotifier
	Chain              *blockchain.BlockChain
	TxMemPool          *mempool.TxPool
	ChainParams        *chaincfg.Params
	DisableCheckpoints bool
	MaxPeers           int
	FeeEstimator       *mempool.FeeEstimator
}
