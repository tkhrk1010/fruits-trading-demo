package main

import (
	"fmt"
	"time"
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/tkhrk1010/fruits-trading-demo/actors/aggregator"
	"github.com/tkhrk1010/fruits-trading-demo/actors/collector"
	// "github.com/tkhrk1010/fruits-trading-demo/infra/rdb"
)

func main() {
	system := actor.NewActorSystem()

	orangeInventoryProps := actor.PropsFromProducer(func() actor.Actor { return &collector.OrangeInventoryCollectorActor{} })
	orangeInventoryCollector := system.Root.Spawn(orangeInventoryProps)

	AppleInventoryProps := actor.PropsFromProducer(func() actor.Actor { return &collector.AppleInventoryCollectorActor{} })
	AppleInventoryCollector := system.Root.Spawn(AppleInventoryProps)

	bananaInventoryProps := actor.PropsFromProducer(func() actor.Actor { return &collector.BananaInventoryCollectorActor{} })
	bananaInventoryCollector := system.Root.Spawn(bananaInventoryProps)

	aggregatorProps := actor.PropsFromProducer(func() actor.Actor {
		return &aggregator.AggregatorActor{
			Collectors: map[string]*actor.PID{
				"appleInventory":  AppleInventoryCollector,
				"orangeInventory": orangeInventoryCollector,
				"bananaInventory": bananaInventoryCollector,
			},
		}
	})
	aggrgtr := system.Root.Spawn(aggregatorProps)

	future := system.Root.RequestFuture(aggrgtr, &aggregator.AggregateRequest{
		ItemNames: []string{"orangeInventory", "appleInventory", "bananaInventory"},
	}, 10*time.Second)

	result, err := future.Result()
	if err != nil {
		fmt.Printf("Error while waiting for result: %s\n", err)
		return
	}

	response, ok := result.(*aggregator.AggregateResponse)
	if !ok {
		fmt.Println("Invalid response type")
		return
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	fmt.Println("Aggregate Results (JSON):", string(jsonResponse))

	// db, err := rdb.ConnectDB()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer db.Close()

	// tableName := "results"
	// if err := rdb.InsertData(db, tableName, response.Results); err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("Data inserted successfully.")

}
