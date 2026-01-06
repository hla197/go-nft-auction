package handlers

import (
	"chain/internal/infra/db"
	"chain/internal/logger"
	"chain/internal/models"
	"chain/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuctionHandler struct {
}

type NftHistoryRequest struct {
	*utils.FieldValidate
	NftContract string `json:"nft_contract" binding:"required" label:"nft合约地址"`
	NftTokenId  string `json:"nft_token_id" binding:"required" label:"nft token_id"`
	ChainID     string `json:"chain_id" binding:"required" label:"链id"`
}

// 获取NFT历史拍卖
func (h *AuctionHandler) GetNftHistory(c *gin.Context) {
	var req NftHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var validate utils.FieldValidateIF = req
		logger.Log.Errorf("GetNftHistory err: %v", err)
		msg := validate.Validate(err, req)
		utils.Fail(c, 500, msg)
		return
	}

	logger.Log.Infof("GetNftHistory req: %v", req)

	history := []models.Auction{}
	db.DB.Where("nft_contract = ? and nft_token_id = ? and chain_id = ? and active = 0", req.NftContract, req.NftTokenId, req.ChainID).Order("end_time desc").Find(&history)
	utils.Success(c, history, "")
	return
}
