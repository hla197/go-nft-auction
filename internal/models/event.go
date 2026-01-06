package models

import "gorm.io/gorm"

type Events struct {
	Id          uint64 `json:"id" gorm:"primary_key;AUTO_INCREMENT" `
	Address     string `json:"address" gorm:"type:char(42)" `
	Data        string `json:"data" gorm:"type:longtext" `
	BlockNumber uint64 `json:"block_number"`
	TxHash      string `json:"tx_hash" gorm:"type:char(66)" `
	TxIndex     uint   `json:"tx_index" `
	BlockHash   string `json:"block_hash" gorm:"type:varchar(256)" `
	Topics      string `json:"topics" gorm:"type:longtext"` // 使用 JSON 格式存储字符串数组
	LogIndex    uint   `json:"log_index"`
	Removed     bool   `json:"removed"`
	*gorm.Model
}
