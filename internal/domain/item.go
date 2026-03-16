package domain

type Item struct {
	itemStaticData  *ItemStaticData
	itemDynamicData ItemDynamicData
}

func NewItem(staticData *ItemStaticData, dynamicData ItemDynamicData) *Item {
	return &Item{
		itemStaticData:  staticData,
		itemDynamicData: dynamicData,
	}
}

func JoinItem(staticData *ItemStaticData, dynamicData ItemDynamicData) *Item {
	return &Item{
		itemStaticData:  staticData,
		itemDynamicData: dynamicData,
	}
}
