package main

import (
	"fmt"

	"github.com/tkhrk1010/fruits-trading-demo/actors/handler"
)

func main() {

	tradeInformationHandler := handler.NewTradeInformationHandler()
	jsonResponse, err := tradeInformationHandler.GetTradeInformation()
	if err != nil {
		fmt.Println("Error while getting trade information:", err)
		return
	}

	fmt.Println("Aggregate Results (JSON):", string(jsonResponse))
}