package reddcoin

import (
	"encoding/json"

	"github.com/golang/glog"
	"github.com/juju/errors"
	"github.com/trezor/blockbook/bchain"
	"github.com/trezor/blockbook/bchain/coins/btc"
	"github.com/trezor/blockbook/common"
)

// ReddcoinRPC is an interface to JSON-RPC bitcoind service.
type ReddcoinRPC struct {
	*btc.BitcoinRPC
}

// ResGetBlockChainInfo is a response to GetChainInfo request
type ResGetBlockChainInfo struct {
	Error  *bchain.RPCError `json:"error"`
	Result struct {
		Chain         string            `json:"chain"`
		Blocks        int               `json:"blocks"`
		Headers       int               `json:"headers"`
		Bestblockhash string            `json:"bestblockhash"`
		Difficulty    common.JSONNumber `json:"difficulty"`
		MoneySupply   common.JSONNumber `json:"moneysupply"`
		SizeOnDisk    int64             `json:"size_on_disk"`
		Warnings      string            `json:"warnings"`
	} `json:"result"`
}

// NewReddcoinRPC is an interface to JSON-RPC bitcoind service.
func NewReddcoinRPC(config json.RawMessage, pushHandler func(bchain.NotificationType)) (bchain.BlockChain, error) {
	b, err := btc.NewBitcoinRPC(config, pushHandler)
	if err != nil {
		return nil, err
	}

	s := &ReddcoinRPC{
		b.(*btc.BitcoinRPC),
	}
	s.RPCMarshaler = btc.JSONMarshalerV2{}
	s.ChainConfig.SupportsEstimateFee = false

	return s, nil
}

// Initialize initializes ReddcoinRPC instance.
func (b *ReddcoinRPC) Initialize() error {
	ci, err := b.GetChainInfo()
	if err != nil {
		return err
	}
	chainName := ci.Chain

	glog.Info("Chain name ", chainName)
	params := GetChainParams(chainName)

	b.Parser = NewReddcoinParser(params, b.ChainConfig)

	// parameters for getInfo request
	if params.Net == MainReddNet {
		b.Testnet = false
		b.Network = "livenet"
	} else {
		b.Testnet = true
		b.Network = "testnet"
	}

	glog.Info("rpc: block chain ", params.Name)

	return nil
}

// GetChainInfo return info about the blockchain
func (b *ReddcoinRPC) GetChainInfo() (*bchain.ChainInfo, error) {
	chainInfo := ResGetBlockChainInfo{}
	err := b.Call(&btc.CmdGetBlockChainInfo{Method: "getblockchaininfo"}, &chainInfo)
	if err != nil {
		return nil, err
	}
	if chainInfo.Error != nil {
		return nil, chainInfo.Error
	}

	networkInfo := btc.ResGetNetworkInfo{}
	err = b.Call(&btc.CmdGetNetworkInfo{Method: "getnetworkinfo"}, &networkInfo)
	if err != nil {
		return nil, err
	}
	if networkInfo.Error != nil {
		return nil, networkInfo.Error
	}

	return &bchain.ChainInfo{
		Bestblockhash:   chainInfo.Result.Bestblockhash,
		Blocks:          chainInfo.Result.Blocks,
		Chain:           chainInfo.Result.Chain,
		Difficulty:      string(chainInfo.Result.Difficulty),
		MoneySupply:     string(chainInfo.Result.MoneySupply),
		Headers:         chainInfo.Result.Headers,
		SizeOnDisk:      chainInfo.Result.SizeOnDisk,
		Version:         string(networkInfo.Result.Version),
		Subversion:      string(networkInfo.Result.Subversion),
		ProtocolVersion: string(networkInfo.Result.ProtocolVersion),
		Timeoffset:      networkInfo.Result.Timeoffset,
		Warnings:        networkInfo.Result.Warnings,
	}, nil
}

// GetPeerInfo returns info about connected peers
func (b *ReddcoinRPC) GetPeerInfo() ([]bchain.PeerInfo, error) {
	return b.BitcoinRPC.GetPeerInfo()
}

// GetBlock returns block with given hash.
func (s *ReddcoinRPC) GetBlock(hash string, height uint32) (*bchain.Block, error) {
	var err error
	if hash == "" && height > 0 {
		hash, err = s.GetBlockHash(height)
		if err != nil {
			return nil, err
		}
	}

	glog.V(1).Info("rpc: getblock (verbosity=1) ", hash)

	res := btc.ResGetBlockThin{}
	req := btc.CmdGetBlock{Method: "getblock"}
	req.Params.BlockHash = hash
	req.Params.Verbosity = 1
	err = s.Call(&req, &res)

	if err != nil {
		return nil, errors.Annotatef(err, "hash %v", hash)
	}
	if res.Error != nil {
		return nil, errors.Annotatef(res.Error, "hash %v", hash)
	}

	txs := make([]bchain.Tx, 0, len(res.Result.Txids))
	for _, txid := range res.Result.Txids {
		tx, err := s.GetTransaction(txid)
		if err != nil {
			if err == bchain.ErrTxNotFound {
				glog.Errorf("rpc: getblock: skipping transanction in block %s due error: %s", hash, err)
				continue
			}
			return nil, err
		}
		txs = append(txs, *tx)
	}
	block := &bchain.Block{
		BlockHeader: res.Result.BlockHeader,
		Txs:         txs,
	}
	return block, nil
}

// GetTransactionForMempool returns a transaction by the transaction ID.
// It could be optimized for mempool, i.e. without block time and confirmations
func (s *ReddcoinRPC) GetTransactionForMempool(txid string) (*bchain.Tx, error) {
	return s.GetTransaction(txid)
}
