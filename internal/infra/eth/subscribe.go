package eth

import (
	"context"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

// ListenNewBlocks WSS 订阅
func ListenNewBlocks(ctx context.Context, mgr *EthClientManager, handler func(*types.Header)) {
	for {
		client := mgr.WSS()
		if client == nil {
			time.Sleep(time.Second)
			continue
		}

		headers := make(chan *types.Header)
		sub, err := client.SubscribeNewHead(ctx, headers)
		if err != nil {
			mgr.OnWssError(err)
			time.Sleep(time.Second)
			continue
		}

		log.Println("[eth] subscribed new blocks")

		for {
			select {
			case <-sub.Err():
				log.Println("[eth] subscription error")
				mgr.OnWssError(err)
				goto RESUB

			case <-mgr.WssReconnect():
				log.Println("[eth] wss reconnected, resubscribe")
				goto RESUB

			case h := <-headers:
				handler(h)
			}
		}

	RESUB:
		continue
	}
}
