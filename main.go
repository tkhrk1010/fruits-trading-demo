package main

import (
	"fmt"

	"github.com/tkhrk1010/fruits-trading-demo/actors/handler"
)

func main() {

	tradeSupportInformationHandler := handler.NewTradeSupportInformationHandler()
	jsonResponse, err := tradeSupportInformationHandler.GetTradeSupportInformation()
	if err != nil {
		fmt.Println("Error while getting trade information:", err)
		return
	}

	fmt.Println("Aggregate Results (JSON):", string(jsonResponse))
}
