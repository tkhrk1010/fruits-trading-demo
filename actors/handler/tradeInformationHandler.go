package handler

import (
	"fmt"
	"time"
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/tkhrk1010/fruits-trading-demo/actors/inventory/aggregator"
	"github.com/tkhrk1010/fruits-trading-demo/actors/inventory/collector"
)

type TradeInformationHandler struct {
	system *actor.ActorSystem
}

func NewTradeInformationHandler() *TradeInformationHandler {
	return &TradeInformationHandler{
		system: actor.NewActorSystem(),
	}
}

func (handler *TradeInformationHandler) GetTradeInformation() ([]byte, error) {
	orangeInventoryProps := actor.PropsFromProducer(func() actor.Actor { return &collector.OrangeInventoryCollectorActor{} })
	orangeInventoryCollector := handler.system.Root.Spawn(orangeInventoryProps)

	AppleInventoryProps := actor.PropsFromProducer(func() actor.Actor { return &collector.AppleInventoryCollectorActor{} })
	AppleInventoryCollector := handler.system.Root.Spawn(AppleInventoryProps)

	bananaInventoryProps := actor.PropsFromProducer(func() actor.Actor { return &collector.BananaInventoryCollectorActor{} })
	bananaInventoryCollector := handler.system.Root.Spawn(bananaInventoryProps)

	aggregatorProps := actor.PropsFromProducer(func() actor.Actor {
		return &aggregator.AggregatorActor{
			Collectors: map[string]*actor.PID{
				"appleInventory":  AppleInventoryCollector,
				"orangeInventory": orangeInventoryCollector,
				"bananaInventory": bananaInventoryCollector,
			},
		}
	})
	aggrgtr := handler.system.Root.Spawn(aggregatorProps)

	future := handler.system.Root.RequestFuture(aggrgtr, &aggregator.AggregateRequest{
		ItemNames: []string{"orangeInventory", "appleInventory", "bananaInventory"},
	}, 10*time.Second)

	result, err := future.Result()
	if err != nil {
		fmt.Printf("Error while waiting for result: %s\n", err)
		return nil, err
	}

	response, ok := result.(*aggregator.AggregateResponse)
	if !ok {
		fmt.Println("Invalid response type")
		return nil, err
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return nil, err
	}

	return jsonResponse, nil

}
