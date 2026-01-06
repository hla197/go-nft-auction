package models

import "gorm.io/gorm"

type Auction struct {
	Id             uint64 `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	AuctionId      uint64 `json:"auction_id" gorm:"uniqueIndex:uk_chain_auction"`
	NftContract    string `json:"nft_contract" comment:"nft合约地址" `
	NftTokenId     uint64 `json:"nft_token_id" gorm:"type:char(42)" `
	ChainId        int64  `json:"chain_id" comment:"链id"`
	OwnerAddress   string `json:"owner_address" comment:"拥有者地址" `
	StartPrice     uint64 `json:"start_price" comment:"开始价格" `
	StartPriceUsd  uint64 `json:"start_price_usd" comment:"开始价格usd" `
	BidAddress     string `json:"bid_address" comment:"竞拍者地址" `
	BidPrice       uint64 `json:"bid_price" comment:"竞拍价格" `
	BidPriceUsd    uint64 `json:"bid_price_usd" comment:"竞拍价格usd" `
	BidTokenId     string `json:"bid_token_id" comment:"竞拍token_id 0x0为ETH" `
	BidTime        int64  `json:"bid_time" comment:"竞拍时间" `
	StartTime      int64  `json:"start_time" comment:"开始时间" `
	EndTime        int64  `json:"end_time" comment:"结束时间" `
	Active         int64  `json:"active" comment:"状态" `
	TxHash         string `json:"tx_hash" comment:"交易hash" `
	CreatorAddress string `json:"creator_address" comment:"创建者地址" `
	BlockNumber    uint64 `json:"block_number" comment:"块高" `
	BlockTime      uint64 `json:"block_time" comment:"块时间" `
	BlockHash      string `json:"block_hash" comment:"块hash" `
	CreatedAt      int64  `json:"created_at" comment:"创建时间" `
	UpdatedAt      int64  `json:"updated_at" comment:"更新时间" `
	*gorm.Model
}
