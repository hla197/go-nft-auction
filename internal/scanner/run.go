package scanner

import (
	"chain/internal/global"
	"chain/internal/infra/eth"
)

func Run() {
	// 启动监听区块
	// scanNft := NewScanNft(eth.EthMgr.RPC(), global.BlockChainConfig.NftContract, global.BlockChainConfig.NftStartBlock)
	// scanNft.Start()

	ScanAuction := NewScanAuction(eth.EthMgr.RPC(), global.BlockChainConfig.NftAuctionContract, global.BlockChainConfig.NftAuctionStartBlock)
	ScanAuction.Start()
}
