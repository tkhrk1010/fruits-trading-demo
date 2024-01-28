package collector

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/tkhrk1010/fruits-trading-demo/domain"
)

type AppleInventoryCollectorActor struct {
	InventoryCollectorActor
}

func (state *AppleInventoryCollectorActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *CollectInventoryRequest:
		fmt.Println("AppleInventoryCollectorActor: Received CollectInventoryRequest")

		// collect apple inventory info
		apples := domain.AppleInventory{}
		averagePurchasingPrice := apples.GetAveragePurchasingPrice()
		count := apples.GetCount()

		ctx.Respond(&InventoryResponse{ItemName: "apple", AveragePurchasingPrice: averagePurchasingPrice, Count: count})

	case *actor.ReceiveTimeout:
		fmt.Println("AppleInventoryCollectorActor: Received timeout")
		ctx.Respond(&InventoryResponse{ItemName: "apple", AveragePurchasingPrice: 0.0, Count: 0})
	}
}
