package collector

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/tkhrk1010/fruits-trading-demo/domain/inventory"
)

type OrangeInventoryCollectorActor struct {
	InventoryCollectorActor
}

func (state *OrangeInventoryCollectorActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *CollectInventoryRequest:
		fmt.Println("OrangeInventoryCollectorActor: Received CollectInventoryRequest")

		// collect orange inventory info
		oranges := inventory.OrangeInventory{}
		averagePurchasingPrice := oranges.GetAveragePurchasingPrice()
		count := oranges.GetCount()

		ctx.Respond(&InventoryResponse{ItemName: "orange", AveragePurchasingPrice: averagePurchasingPrice, Count: count})

	case *actor.ReceiveTimeout:
		fmt.Println("OrangeInventoryCollectorActor: Received timeout")
		ctx.Respond(&InventoryResponse{ItemName: "orange", AveragePurchasingPrice: 0.0, Count: 0})
	}
}
