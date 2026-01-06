package scanner

import (
	"chain/internal/infra/db"
	"chain/internal/infra/redis"
	"chain/internal/logger"
	"chain/internal/models"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm/clause"
)

const AuctionABIJSON = `
[
  {
    "inputs": [
      { "internalType": "uint256", "name": "", "type": "uint256" }
    ],
    "name": "auctions",
    "outputs": [
      { "internalType": "address", "name": "seller", "type": "address" },
      { "internalType": "uint256", "name": "startPrice", "type": "uint256" },
      { "internalType": "uint256", "name": "startPriceUsd", "type": "uint256" },
      { "internalType": "uint256", "name": "highestBid", "type": "uint256" },
      { "internalType": "address", "name": "highestBidder", "type": "address" },
      { "internalType": "uint256", "name": "highestBidUsd", "type": "uint256" },
      { "internalType": "address", "name": "highestBidderToken", "type": "address" },
      { "internalType": "uint256", "name": "endTime", "type": "uint256" },
      { "internalType": "uint256", "name": "startTime", "type": "uint256" },
      { "internalType": "address", "name": "nftAddress", "type": "address" },
      { "internalType": "uint256", "name": "tokenId", "type": "uint256" },
      { "internalType": "address", "name": "token", "type": "address" },
      { "internalType": "bool", "name": "active", "type": "bool" }
    ],
    "stateMutability": "view",
    "type": "function"
  }
]
`

type ScanAuction struct {
	client       *ethclient.Client
	abi          abi.ABI
	ctx          context.Context
	chainID      *big.Int
	contract     common.Address
	batchSize    uint64        // 一次扫 200 个 block
	pollInterval time.Duration // 轮询间隔
	startBlock   uint64
	lastBlockKey string // 最后处理的块高
	lockScanKey  string // 锁
}

/**
 * 扫描 NFT
 * @param client  RPC 客户端
 * @param contractAddress  NFT 合约地址
 * @param startBlock  首次扫描的开始块高，一般为创建合约的区块
 */
func NewScanAuction(client *ethclient.Client, contractAddress string, startBlock uint64) *ScanAuction {
	obj := &ScanAuction{}
	abi, err := abi.JSON(strings.NewReader(AuctionABIJSON))
	if err != nil {
		log.Fatal(err)
	}
	logger.Log.Infoln("[auction scanner] start:", contractAddress)
	obj.abi = abi
	obj.ctx = context.Background()
	obj.client = client
	obj.contract = common.HexToAddress(contractAddress)
	obj.batchSize = 2000
	obj.pollInterval = 5 * time.Second
	obj.startBlock = startBlock
	obj.lastBlockKey = fmt.Sprintf("auction scanner:last_block:%s", contractAddress)
	obj.lockScanKey = fmt.Sprintf("auction scanner:lock:%s", contractAddress)
	chainId, _ := client.ChainID(obj.ctx)
	obj.chainID = chainId
	return obj
}

func (s *ScanAuction) Start() {
	go func() {
		logger.Log.Infoln("ScanAuction Start")
		ticker := time.NewTicker(s.pollInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				s.tickerScan() // 每5秒扫描一次
			case <-s.ctx.Done():
				logger.Log.Infoln("[auction scanner] stopped")
				return
			}
		}
	}()
}

func (s *ScanAuction) HandleEvents(logs []types.Log) error {
	topics := getTopics()
	createdTopic := topics["AuctionCreated"]
	bidTopic := topics["BidPlaced"]
	endedTopic := topics["AuctionEnded"]
	var events []models.Events
	for _, vLog := range logs {
		event := s.getEvent(vLog)
		// 获取交易哈希、区块哈希、块高
		txHash := vLog.TxHash.Hex()
		blockHash := vLog.BlockHash.Hex()
		blockNumber := vLog.BlockNumber
		blockTime := vLog.BlockTimestamp
		events = append(events, *event)

		switch vLog.Topics[0] {
		case createdTopic:
			auctionId := new(big.Int).SetBytes(vLog.Topics[2].Bytes())
			logger.Log.Infoln("Auction created:", auctionId.String())
			s.saveAuctionTopic(auctionId, blockNumber, blockHash, txHash, blockTime)
		case bidTopic:
			auctionId := new(big.Int).SetBytes(vLog.Topics[2].Bytes())
			bidder := vLog.Topics[1].Hex()
			amount := new(big.Int).SetBytes(vLog.Data)
			s.saveAuctionLogTopic(auctionId, bidder, amount, blockNumber, blockHash, txHash, blockTime)
		case endedTopic:
			auctionId := new(big.Int).SetBytes(vLog.Topics[2].Bytes())
			logger.Log.Infoln("Auction created:", auctionId.String())
			s.saveAuctionTopic(auctionId, blockNumber, blockHash, txHash, blockTime)
		}
	}
	if len(events) > 0 {
		if err := db.DB.CreateInBatches(events, 100).Error; err != nil {
			// 如果插入过程中发生错误，则返回错误信息
			logger.Log.Errorf("Error saving events: %v", err)
			return err
		}
	}

	return nil
}

func (s *ScanAuction) saveAuctionLogTopic(auctionId *big.Int, bidder string, amount *big.Int, blockNumber uint64, blockHash string, txHash string, blockTime uint64) {
	// 调用 auctions getter
	auction, err := s.getAuctionFromCache(auctionId)
	if err != nil {
		auction = s.saveAuctionTopic(auctionId, blockNumber, blockHash, txHash, blockTime)
	}

	bidLog := &models.AuctionLog{
		AuctionId:   auctionId.Uint64(),
		NftContract: s.contract.Hex(),
		NftTokenId:  auction.NftTokenId,
		ChainId:     s.chainID.Uint64(),
		BidAddress:  bidder,
		BidPriceUsd: amount.Uint64(),
		BidTime:     time.Now().Unix(),
		TxHash:      txHash,
		BlockNumber: blockNumber,
		BlockHash:   blockHash,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	if err := db.DB.Create(bidLog).Error; err != nil {
		logger.Log.Errorf(err.Error())
	}

	// 如果 txHash 存在，则不做插入，避免重复
	err = db.DB.Where("tx_hash = ?", txHash).FirstOrCreate(&bidLog).Error
	if err != nil {
		logger.Log.Errorf("Error saving auction log: %v", err)
	} else {
		logger.Log.Infoln("Auction log saved:", bidLog)
	}
}

func (s *ScanAuction) saveAuctionTopic(
	auctionId *big.Int,
	blockNumber uint64,
	blockHash string,
	txHash string,
	blockTime uint64,
) *models.Auction {

	data, err := s.abi.Pack("auctions", auctionId)
	if err != nil {
		logger.Log.Errorln("Pack auctions error:", err)
		return nil
	}

	res, err := s.client.CallContract(
		s.ctx,
		ethereum.CallMsg{
			To:   &s.contract,
			Data: data,
		},
		nil,
	)
	if err != nil {
		logger.Log.Errorln("CallContract error:", err)
		return nil
	}

	var a struct {
		Seller             common.Address
		StartPrice         *big.Int
		StartPriceUsd      *big.Int
		HighestBid         *big.Int
		HighestBidder      common.Address
		HighestBidUsd      *big.Int
		HighestBidderToken common.Address
		EndTime            *big.Int
		StartTime          *big.Int
		NftAddress         common.Address
		TokenId            *big.Int
		Token              common.Address
		Active             bool
	}

	if err := s.abi.UnpackIntoInterface(&a, "auctions", res); err != nil {
		logger.Log.Errorln("Unpack auctions error:", err)
		return nil
	}

	now := time.Now().Unix()

	dbAuction := &models.Auction{
		AuctionId:     auctionId.Uint64(),
		NftContract:   a.NftAddress.Hex(),
		NftTokenId:    a.TokenId.Uint64(),
		ChainId:       s.chainID.Int64(),
		OwnerAddress:  a.Seller.Hex(),
		StartPrice:    a.StartPrice.Uint64(),
		StartPriceUsd: a.StartPriceUsd.Uint64(),
		BidAddress:    a.HighestBidder.Hex(),
		BidPrice:      a.HighestBid.Uint64(),
		BidPriceUsd:   a.HighestBidUsd.Uint64(),
		BidTokenId:    a.HighestBidderToken.Hex(),
		StartTime:     a.StartTime.Int64(),
		EndTime:       a.EndTime.Int64(),
		BlockTime:     blockTime,
		Active: func() int64 {
			if a.Active {
				return 1
			}
			return 0
		}(),
		TxHash:         txHash,
		BlockNumber:    blockNumber,
		BlockHash:      blockHash,
		CreatorAddress: a.Seller.Hex(),
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := db.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "auction_id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"owner_address",
			"start_price",
			"start_price_usd",
			"bid_address",
			"bid_price",
			"bid_price_usd",
			"bid_token_id",
			"start_time",
			"end_time",
			"active",
			"tx_hash",
			"block_number",
			"block_hash",
			"updated_at",
		}),
	}).Create(dbAuction).Error; err != nil {
		logger.Log.Errorln("DB save error:", err)
		return nil
	}

	// Redis 缓存
	s.cacheAuction(auctionId, dbAuction)

	return dbAuction
}

func (s *ScanAuction) cacheAuction(auctionId *big.Int, dbAuction *models.Auction) error {
	key := s.getAcutionKey(auctionId) // 返回 Redis key
	data, err := json.Marshal(dbAuction)
	if err != nil {
		return err
	}

	// SetNX：不存在才写入
	ok, err := redis.Client.SetNX(context.Background(), key, data, 7*24*time.Hour).Result()
	if err != nil {
		return err
	}
	if !ok {
		// key 已存在，不覆盖
		return nil
	}
	return nil
}

func (s *ScanAuction) getAuctionFromCache(auctionId *big.Int) (*models.Auction, error) {
	key := s.getAcutionKey(auctionId)
	data, err := redis.Client.Get(s.ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var auction models.Auction
	if err := json.Unmarshal([]byte(data), &auction); err != nil {
		return nil, err
	}
	return &auction, nil
}

func (s *ScanAuction) getAcutionKey(auctionId *big.Int) string {
	return fmt.Sprintf("auction:%s:%d:%d", s.contract.Hex(), auctionId, s.chainID)
}

func (s *ScanAuction) tickerScan() {
	if !s.lock() {
		logger.Log.Infoln("[auction scanner] locked")
		return
	}
	defer s.unlock()
	// 获取最新块高
	lastBlock, err := s.LatestBlock()
	if err != nil {
		logger.Log.Errorln(err)
		return
	}
	// 获取开始块
	startBlock, err := s.ResolvestartBlock(lastBlock)
	// 如果当前区块已经扫描过了，跳过
	if lastBlock <= startBlock {
		return
	}
	if err != nil {
		logger.Log.Errorln(err)
		return
	}
	logs, err := s.findEvents(startBlock, lastBlock)
	if err != nil {
		logger.Log.Errorln(err)
		return
	}

	if err := s.HandleEvents(logs); err != nil {
		logger.Log.Errorln(err)
		return
	}

}

func getTopics() map[string]common.Hash {
	createdSig := []byte("AuctionCreated(address,uint256,uint256,uint256)")
	createdTopic := crypto.Keccak256Hash(createdSig)

	bidSig := []byte("BidPlaced(address,uint256,uint256)")
	bidTopic := crypto.Keccak256Hash(bidSig)

	endedSig := []byte("AuctionEnded(address,uint256,uint256)")
	endedTopic := crypto.Keccak256Hash(endedSig)

	cancelSig := []byte("AuctionCancelled(address,uint256)")
	cancelTopic := crypto.Keccak256Hash(cancelSig)

	return map[string]common.Hash{
		"AuctionCreated":   createdTopic,
		"BidPlaced":        bidTopic,
		"AuctionEnded":     endedTopic,
		"AuctionCancelled": cancelTopic,
	}
}

func (s *ScanAuction) findEvents(startBlock uint64, lastBlock uint64) ([]types.Log, error) {
	endBlock := min(startBlock+s.batchSize-1, lastBlock)
	logger.Log.Infof("ScanAuction findEvents %d -> %d", startBlock, endBlock)

	topics := getTopics()
	createdTopic := topics["AuctionCreated"]
	bidTopic := topics["BidPlaced"]
	endedTopic := topics["AuctionEnded"]

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(startBlock)),
		ToBlock:   big.NewInt(int64(endBlock)),
		Addresses: []common.Address{s.contract},
		Topics: [][]common.Hash{{
			createdTopic,
			bidTopic,
			endedTopic,
		}},
	}

	queryLogs, err := s.client.FilterLogs(s.ctx, query)
	if err != nil {
		return nil, err
	}
	s.SaveLastBlock(endBlock)
	return queryLogs, nil
}

func (s *ScanAuction) getEvent(log types.Log) *models.Events {
	topicBytes, _ := json.Marshal(log.Topics)
	event := &models.Events{
		TxHash:      log.TxHash.Hex(),
		Address:     log.Address.Hex(),
		BlockNumber: log.BlockNumber,
		BlockHash:   log.BlockHash.Hex(),
		TxIndex:     log.TxIndex,
		LogIndex:    log.Index,
		Data:        hex.EncodeToString(log.Data),
		Topics:      string(topicBytes),
		Removed:     log.Removed,
	}
	return event
}

/**
 * 获取开始块
 */
func (s *ScanAuction) ResolvestartBlock(lastBlock uint64) (uint64, error) {
	// 获取Redis中的记录
	val, _ := redis.Client.Get(redis.Ctx, s.lastBlockKey).Result()
	if val != "" {
		return strconv.ParseUint(val, 10, 64)
	}
	if s.startBlock > 0 {
		return s.startBlock, nil
	}
	// 兜底：只扫最近 5000 块
	if lastBlock > 5000 {
		return lastBlock - 5000, nil
	}

	return 0, nil
}

func (s *ScanAuction) LatestBlock() (uint64, error) {
	return s.client.BlockNumber(s.ctx)
}

func (s *ScanAuction) SaveLastBlock(b uint64) error {
	return redis.Client.Set(redis.Ctx, s.lastBlockKey, b, 0).Err()
}

func (s *ScanAuction) lock() bool {
	ok, err := redis.Client.SetNX(redis.Ctx, s.lockScanKey, true, 30*time.Second).Result()
	if err != nil {
		logger.Log.Errorln("Lock error:", err)
		return false
	}
	return ok
}

func (s *ScanAuction) unlock() error {
	return redis.Client.Del(redis.Ctx, s.lockScanKey).Err()
}
