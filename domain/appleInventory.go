package domain

type AppleInventory struct {}

func (appleInventory *AppleInventory) GetAveragePurchasingPrice() float32 {
	return 1.0
}

func (appleInventory *AppleInventory) GetCount() int {
	return 1
}