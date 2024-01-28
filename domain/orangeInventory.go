package domain

type OrangeInventory struct {}

func (orangeInventory *OrangeInventory) GetAveragePurchasingPrice() float32 {
	return 1.0
}

func (orangeInventory *OrangeInventory) GetCount() int {
	return 1
}