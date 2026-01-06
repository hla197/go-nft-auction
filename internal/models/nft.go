package models

import "gorm.io/gorm"

type Nft struct {
	Id                 int64  `json:"id" comment:"id"`
	OwnerAddress       string `json:"owner_address" comment:"拥有者地址"`
	ChainId            int64  `json:"chain_id" comment:"链id"`
	Contract           string `json:"contract" comment:"合约地址"`
	TokenId            string `json:"token_id" comment:"token_id"`
	MetaData           string `json:"meta_data" comment:"元数据"`
	TokenUri           string `json:"token_uri" comment:"token_uri"`
	MintTime           int64  `json:"mint_time" comment:"铸造时间"`
	MintTxHash         string `json:"mint_tx_hash" comment:"铸造交易hash"`
	MintBlock          uint64 `json:"mint_block" comment:"铸造块高"`
	CreatedAddress     string `json:"created_address" comment:"创建者地址"`
	LastTransferTime   int64  `json:"last_transfer_time" comment:"最后转账时间"`
	LastTransferTxHash string `json:"last_transfer_tx_hash" comment:"最后转账交易hash"`
	*gorm.Model
}
