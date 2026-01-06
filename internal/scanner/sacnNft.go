package scanner

import (
	"chain/internal/infra/redis"
	"chain/internal/logger"
	"context"
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
)

// ===== ERC721 Transfer ABI（只要 event 就够）=====
const Erc721ABIJSON = `
[
  {
    "anonymous": false,
    "inputs": [
      { "indexed": true, "internalType": "address", "name": "from", "type": "address" },
      { "indexed": true, "internalType": "address", "name": "to", "type": "address" },
      { "indexed": true, "internalType": "uint256", "name": "tokenId", "type": "uint256" }
    ],
    "name": "Transfer",
    "type": "event"
  }
]
`

type ScanNft struct {
	client *ethclient.Client
	abi    abi.ABI
	ctx    context.Context

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
func NewScanNft(client *ethclient.Client, contractAddress string, startBlock uint64) *ScanNft {
	obj := &ScanNft{}
	abi, err := abi.JSON(strings.NewReader(Erc721ABIJSON))
	if err != nil {
		log.Fatal(err)
	}
	obj.ctx = context.Background()
	obj.client = client
	obj.contract = common.HexToAddress(contractAddress)
	obj.batchSize = 2000
	obj.pollInterval = 5 * time.Second
	obj.startBlock = startBlock
	obj.lastBlockKey = fmt.Sprintf("scanner:last_block:%s", contractAddress)
	obj.lockScanKey = fmt.Sprintf("scanner:lock:%s", contractAddress)
	obj.abi = abi
	return obj
}

func (s *ScanNft) Start() {
	go func() {
		logger.Log.Infoln("ScanNft Start")
		ticker := time.NewTicker(s.pollInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				logger.Log.Infoln("ScanNft Start")
				s.tickerScan() // 每5秒扫描一次
			case <-s.ctx.Done():
				logger.Log.Infoln("[scanner] stopped")
				return
			}
		}
	}()
}

func (s *ScanNft) HandleEvents(logs []types.Log) error {
	for _, log := range logs {
		if log.Removed {
			continue
		}
		from := common.HexToAddress(log.Topics[1].Hex())
		to := common.HexToAddress(log.Topics[2].Hex())
		tokenId := new(big.Int).SetBytes(log.Topics[3].Bytes())
		logger.Log.Infof("NFT Transfer %s -> %s tokenId=%s", from.Hex(), to.Hex(), tokenId.String())
	}
	return nil
}

func (s *ScanNft) tickerScan() {
	if !s.lock() {
		logger.Log.Infoln("[scanner] locked")
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
	if err != nil {
		logger.Log.Errorln(err)
		return
	}
	logs, err := s.findEvents(startBlock, lastBlock)
	if err != nil {
		logger.Log.Errorln(err)
		return
	}
	logger.Log.Infoln("ScanNft handleEvents", logs)
	if err := s.HandleEvents(logs); err != nil {
		logger.Log.Errorln(err)
		return
	}

}

func (s *ScanNft) findEvents(startBlock uint64, lastBlock uint64) ([]types.Log, error) {
	logger.Log.Infoln("ScanNft findEvents start")
	endBlock := min(startBlock+s.batchSize-1, lastBlock)
	logger.Log.Infof("ScanNft findEvents %d -> %d", startBlock, endBlock)

	transferSig := []byte("Transfer(address,address,uint256)")
	transferTopic := crypto.Keccak256Hash(transferSig)

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(startBlock)),
		ToBlock:   big.NewInt(int64(endBlock)),
		Addresses: []common.Address{s.contract},
		Topics:    [][]common.Hash{{transferTopic}},
	}

	queryLogs, err := s.client.FilterLogs(s.ctx, query)
	if err != nil {
		return nil, err
	}
	s.SaveLastBlock(endBlock)
	return queryLogs, nil
}

/**
 * 获取开始块
 */
func (s *ScanNft) ResolvestartBlock(lastBlock uint64) (uint64, error) {
	last, err := s.GetLastBlock()
	if err == nil {
		return last, nil
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

/**
 * 获取最后处理的块高
 */
func (s *ScanNft) GetLastBlock() (uint64, error) {
	// 获取Redis中的记录
	val, _ := redis.Client.Get(redis.Ctx, s.lastBlockKey).Result()
	if val != "" {
		return strconv.ParseUint(val, 10, 64)
	}

	latest, err := s.LatestBlock()
	if err != nil {
		return 0, err
	}
	return latest, nil
}

func (s *ScanNft) LatestBlock() (uint64, error) {
	return s.client.BlockNumber(s.ctx)
}

func (s *ScanNft) SaveLastBlock(b uint64) error {
	return redis.Client.Set(redis.Ctx, s.lastBlockKey, b, 0).Err()
}

func (s *ScanNft) lock() bool {
	ok, err := redis.Client.SetNX(redis.Ctx, s.lockScanKey, true, 30*time.Second).Result()
	if err != nil {
		logger.Log.Errorln("Lock error:", err)
		return false
	}
	return ok
}

func (s *ScanNft) unlock() error {
	return redis.Client.Del(redis.Ctx, s.lockScanKey).Err()
}
