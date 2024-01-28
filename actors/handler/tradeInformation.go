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

func NewTradeInformationHandler() *TradeInformationHandler {
	return &TradeInformationHandler{
		system: actor.NewActorSystem(),
	}
}

func (handler *TradeInformationHandler) spawnActor(producer func() actor.Actor) *actor.PID {
	props := actor.PropsFromProducer(producer)
	return handler.system.Root.Spawn(props)
}

func (handler *TradeInformationHandler) requestActorResponse(props *actor.Props, request interface{}, timeout time.Duration) (interface{}, error) {
	actorRef := handler.system.Root.Spawn(props)
	future := handler.system.Root.RequestFuture(actorRef, request, timeout)
	return future.Result()
}

func (handler *TradeInformationHandler) collectInventoryInfo() (*aggregator.AggregateResponse, error) {
	orangeInventoryCollector := handler.spawnActor(func() actor.Actor { return &collector.OrangeActor{} })
	appleInventoryCollector := handler.spawnActor(func() actor.Actor { return &collector.AppleActor{} })
	bananaInventoryCollector := handler.spawnActor(func() actor.Actor { return &collector.BananaActor{} })

	aggregatorProps := actor.PropsFromProducer(func() actor.Actor {
		return &aggregator.AggregatorActor{
			Collectors: map[string]*actor.PID{
				"appleInventory":  appleInventoryCollector,
				"orangeInventory": orangeInventoryCollector,
				"bananaInventory": bananaInventoryCollector,
			},
		}
	})

	result, err := handler.requestActorResponse(aggregatorProps, &aggregator.AggregateRequest{
		ItemNames: []string{"orangeInventory", "appleInventory", "bananaInventory"},
	}, 10*time.Second)
	if err != nil {
		return nil, err
	}

	aggResp, ok := result.(*aggregator.AggregateResponse)
	if !ok {
		return nil, fmt.Errorf("invalid response type")
	}
	return aggResp, nil
}

func (handler *TradeInformationHandler) collectMarketInfo() (*market.Response, error) {
	marketProps := actor.PropsFromProducer(func() actor.Actor { return &market.Actor{} })

	result, err := handler.requestActorResponse(marketProps, &market.Request{
		ItemNames: []string{"apple", "orange", "banana"},
	}, 10*time.Second)
	if err != nil {
		return nil, err
	}

	mktResp, ok := result.(*market.Response)
	if !ok {
		return nil, fmt.Errorf("invalid response type")
	}
	return mktResp, nil
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

func (handler *TradeInformationHandler) GetTradeInformation() ([]byte, error) {
	aggregatorResponse, err := handler.collectInventoryInfo()
	if err != nil {
		return nil, err
	}

	marketResponse, err := handler.collectMarketInfo()
	if err != nil {
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

type CombinedItemDetails struct {
	AveragePurchasingPrice float32 `json:"averagePurchasingPrice"`
	Count                  int     `json:"count"`
	MarketPrice            float32 `json:"marketPrice"`
}

type CombinedResponse struct {
	Items map[string]CombinedItemDetails `json:"items"`
}
