# fruits-trading-demo
trade fruits with actor model and scatter-gather pattern

## Usecases
### Register Purchasing
```mermaid
graph TD
  Apple("Apple register API")
  Orange("Orange register API")
  Banana("Banana register API")
  InventoryDB("Inventory Management DB")

  Apple -->|"register"| InventoryDB
  Orange -->|"register"| InventoryDB
  Banana -->|"register"| InventoryDB
```

### Trade judgement

```mermaid
sequenceDiagram
  participant Trader
  participant InfoCollector 
  participant InventoryAggregator
  participant AppleInventoryCollector
  participant OrangeInventoryCollector
  participant BananaInventoryCollector
  participant MarketPriceCollector
  participant InventoryDB
  participant TradeResultDB

  Trader ->> InfoCollector : Start Transaction
  InfoCollector ->> MarketPriceCollector : Request Market Price
  InfoCollector ->> InventoryAggregator : Collect Information
  InventoryAggregator ->> AppleInventoryCollector : Get Inventory Information
  InventoryAggregator ->> OrangeInventoryCollector : Get Inventory Information
  InventoryAggregator ->> BananaInventoryCollector : Get Inventory Information
  AppleInventoryCollector ->> InventoryDB : Get Inventory Information
  InventoryDB -->> AppleInventoryCollector :  Inventory Information
  AppleInventoryCollector -->> InventoryAggregator : Inventory Information
  OrangeInventoryCollector ->> InventoryDB : Get Inventory Information
  InventoryDB -->> OrangeInventoryCollector :  Inventory Information
  OrangeInventoryCollector -->> InventoryAggregator : Inventory Information
  BananaInventoryCollector ->> InventoryDB : Get Inventory Information
  InventoryDB -->> BananaInventoryCollector :  Inventory Information
  BananaInventoryCollector -->> InventoryAggregator : Inventory Information
  InventoryAggregator -->> InfoCollector : Inventory Information
  MarketPriceCollector -->> InfoCollector : Market Price Information
  InfoCollector -->> Trader : Information
  Trader ->> Trader : judge sell/hold
  Trader ->> TradeResultDB : Save Result

```

## Actor structure
```mermaid
graph TD
  InfoCollector --> PriceCollector
  InfoCollector --> InventoryAggregator
  InventoryAggregator --> AppleInventoryCollector
  InventoryAggregator --> OrangeInventoryCollector
  InventoryAggregator --> BananaInventoryCollector
```
## QuickStart
```
$ docker-compose up -d
$ go run .
```