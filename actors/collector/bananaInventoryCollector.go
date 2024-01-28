package collector

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/tkhrk1010/fruits-trading-demo/domain"
)

type BananaInventoryCollectorActor struct {
	InventoryCollectorActor
}

func (state *BananaInventoryCollectorActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *CollectInventoryRequest:
		fmt.Println("BananaInventoryCollectorActor: Received CollectInventoryRequest")

		// collect banana inventory info
		bananas := domain.BananaInventory{}
		averagePurchasingPrice := bananas.GetAveragePurchasingPrice()
		count := bananas.GetCount()

		ctx.Respond(&InventoryResponse{ItemName: "banana", AveragePurchasingPrice: averagePurchasingPrice, Count: count})

	case *actor.ReceiveTimeout:
		fmt.Println("BananaInventoryCollectorActor: Received timeout")
		ctx.Respond(&InventoryResponse{ItemName: "banana", AveragePurchasingPrice: 0.0, Count: 0})
	}
}
