package models

import "gorm.io/gorm"

type AuctionLog struct {
	Id          uint64 `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	AuctionId   uint64 `json:"auction_id"`
	NftContract string `json:"nft_contract" comment:"nft合约地址" `
	NftTokenId  uint64 `json:"nft_token_id" gorm:"type:char(42)" `
	ChainId     uint64 `json:"chain_id" comment:"链id"`
	BidAddress  string `json:"bid_address" comment:"竞拍者地址" `
	BidPriceUsd uint64 `json:"bid_price_usd" comment:"竞拍价格usd" `
	BidTime     int64  `json:"bid_time" comment:"竞拍时间" `
	TxHash      string `json:"tx_hash" comment:"交易hash" uniqueIndex:uk_tx_hash `
	BlockNumber uint64 `json:"block_number" comment:"块高" `
	BlockHash   string `json:"block_hash" comment:"块hash" `
	CreatedAt   int64  `json:"created_at" comment:"创建时间" `
	UpdatedAt   int64  `json:"updated_at" comment:"更新时间" `
	*gorm.Model
}
