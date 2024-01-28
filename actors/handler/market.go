package handler

import (
	"time"
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/tkhrk1010/fruits-trading-demo/actors/market"
)

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
