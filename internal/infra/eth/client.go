package eth

import (
	"chain/internal/logger"
	"context"
	"log"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

var EthMgr *EthClientManager

// EthClientManager 管理 RPC / WSS 连接
type EthClientManager struct {
	rpcURL string
	wssURL string

	rpc atomic.Value // *ethclient.Client
	wss atomic.Value // *ethclient.Client

	wssReconnect atomic.Value // chan struct{}

	ctx    context.Context
	cancel context.CancelFunc
}

// NewEthClientManager 创建 Manager
func NewEthClientManager(rpcURL, wssURL string) {
	ctx, cancel := context.WithCancel(context.Background())
	logger.Log.Infoln("[eth] init manager", rpcURL, wssURL)
	EthMgr = &EthClientManager{
		rpcURL: rpcURL,
		wssURL: wssURL,
		ctx:    ctx,
		cancel: cancel,
	}

	// 初始化 WSS 重连信号
	EthMgr.wssReconnect.Store(make(chan struct{}))

	// 初始化 RPC
	if err := EthMgr.reconnectRPC(); err != nil {
		logger.Log.Panicf("[eth] init rpc failed:", err)
		return
	}

	// 初始化 WSS（可选）
	if wssURL != "" {
		if err := EthMgr.reconnectWSS(); err != nil {
			logger.Log.Panicf("[eth] init wss failed:", err)
			return
		}
	}

}

// CallRPC 封装 client 调用
func CallRPC(ctx context.Context, mgr *EthClientManager, f func(client *ethclient.Client) error) error {
	var lastErr error
	const maxRetry = 3
	for i := 0; i < maxRetry; i++ {
		client := mgr.RPC()
		if client == nil {
			// client 还在重连
			time.Sleep(time.Second)
			continue
		}

		lastErr = f(client)
		if lastErr == nil {
			return nil
		}

		// 出错通知 manager，触发重连
		mgr.OnRpcError(lastErr)

		time.Sleep(500 * time.Millisecond)
	}
	return lastErr
}

//
// ===== Client 获取（无锁）=====
//

func (m *EthClientManager) RPC() *ethclient.Client {
	logger.Log.Infoln("[eth] get rpc", m.rpc)

	if v := m.rpc.Load(); v != nil {
		return v.(*ethclient.Client)
	}
	return nil
}

func (m *EthClientManager) WSS() *ethclient.Client {
	if v := m.wss.Load(); v != nil {
		return v.(*ethclient.Client)
	}
	return nil
}

//
// ===== 错误驱动重连 =====
//

// OnRpcError RPC 调用报错时调用
func (m *EthClientManager) OnRpcError(err error) {
	if err == nil {
		return
	}
	log.Println("[eth] rpc error:", err)
	go m.reconnectRPC()
}

// OnWssError WSS 订阅报错时调用
func (m *EthClientManager) OnWssError(err error) {
	if err == nil {
		return
	}
	logger.Log.Errorln("[eth] wss error:", err)
	go m.reconnectWSS()
}

//
// ===== 重连实现 =====
//

func (m *EthClientManager) reconnectRPC() error {
	client, err := ethclient.Dial(m.rpcURL)
	if err != nil {
		return err
	}
	old := m.rpc.Swap(client)
	if old != nil {
		old.(*ethclient.Client).Close()
	}

	logger.Log.Infoln("[eth] rpc reconnected")
	return nil
}

func (m *EthClientManager) reconnectWSS() error {
	client, err := ethclient.Dial(m.wssURL)
	if err != nil {
		return err
	}

	old := m.wss.Swap(client)
	if old != nil {
		old.(*ethclient.Client).Close()
	}

	// 通知所有订阅者重建
	oldCh := m.wssReconnect.Load().(chan struct{})
	close(oldCh)
	m.wssReconnect.Store(make(chan struct{}))

	logger.Log.Infoln("[eth] wss reconnected")
	return nil
}

//
// ===== 订阅恢复信号 =====
//

func (m *EthClientManager) WssReconnect() <-chan struct{} {
	return m.wssReconnect.Load().(chan struct{})
}

//
// ===== 关闭 =====
//

func (m *EthClientManager) Close() {
	m.cancel()

	if c := m.rpc.Load(); c != nil {
		c.(*ethclient.Client).Close()
	}

	if c := m.wss.Load(); c != nil {
		c.(*ethclient.Client).Close()
	}
}
