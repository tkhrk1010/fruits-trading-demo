package handler

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/tkhrk1010/fruits-trading-demo/actors/inventory/aggregator"
	"github.com/tkhrk1010/fruits-trading-demo/actors/inventory/collector"
	"github.com/tkhrk1010/fruits-trading-demo/actors/market"
)

type TradeInformationHandler struct {
	system *actor.ActorSystem
}


type CombinedItemDetails struct {
	AveragePurchasingPrice float32 `json:"averagePurchasingPrice"`
	Count                  int     `json:"count"`
	MarketPrice            float32 `json:"marketPrice"`
}

type CombinedResponse struct {
	Items map[string]CombinedItemDetails `json:"items"`
}

func NewTradeInformationHandler() *TradeInformationHandler {
	return &TradeInformationHandler{
		system: actor.NewActorSystem(),
	}
}

func (handler *TradeInformationHandler) GetTradeInformation() ([]byte, error) {
	orangeInventoryProps := actor.PropsFromProducer(func() actor.Actor { return &collector.OrangeActor{} })
	orangeInventoryCollector := handler.system.Root.Spawn(orangeInventoryProps)

	AppleInventoryProps := actor.PropsFromProducer(func() actor.Actor { return &collector.AppleActor{} })
	AppleInventoryCollector := handler.system.Root.Spawn(AppleInventoryProps)

	bananaInventoryProps := actor.PropsFromProducer(func() actor.Actor { return &collector.BananaActor{} })
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

	aggregatorResponse, ok := result.(*aggregator.AggregateResponse)
	if !ok {
		fmt.Println("Invalid response type")
		return nil, err
	}

	//
	// market
	//
	marketProps := actor.PropsFromProducer(func() actor.Actor { return &market.Actor{} })
	marketActor := handler.system.Root.Spawn(marketProps)
	marketFuture := handler.system.Root.RequestFuture(marketActor, &market.Request{
		ItemNames: []string{"apple", "orange", "banana"},
	}, 10*time.Second)

	marketResult, err := marketFuture.Result()
	if err != nil {
		fmt.Printf("Error while waiting for market result: %s\n", err)
		return nil, err
	}

	marketResponse, ok := marketResult.(*market.Response)
	if !ok {
		fmt.Println("Invalid market response type")
		return nil, err
	}

	combinedResponse := CombineResponses(aggregatorResponse, marketResponse)

	jsonResponse, err := json.Marshal(combinedResponse)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return nil, err
	}
	
	return jsonResponse, nil
	
}


func CombineResponses(aggregatorResponse *aggregator.AggregateResponse, marketResponse *market.Response) CombinedResponse {
	combined := CombinedResponse{
		Items: make(map[string]CombinedItemDetails),
	}

	for itemName, aggDetails := range aggregatorResponse.Results {
		combined.Items[itemName] = CombinedItemDetails{
			AveragePurchasingPrice: aggDetails.AveragePurchasingPrice,
			Count:                  aggDetails.Count,
		}
	}

	for itemName, marketItemDetails := range marketResponse.Results {
		if item, exists := combined.Items[itemName]; exists {
			item.MarketPrice = marketItemDetails.MarketPrice 
			combined.Items[itemName] = item
		}
	}

	return combined
}
