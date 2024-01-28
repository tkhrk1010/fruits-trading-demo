package domain

type BananaInventory struct {}

func (bananaInventory *BananaInventory) GetAveragePurchasingPrice() float32 {
	return 1.0
}

func (bananaInventory *BananaInventory) GetCount() int {
	return 1
}