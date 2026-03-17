package domain

type Item struct {
	ItemStaticData  *ItemStaticData
	ItemDynamicData ItemDynamicData
}

func NewItem(staticData *ItemStaticData, dynamicData ItemDynamicData) *Item {
	return &Item{
		ItemStaticData:  staticData,
		ItemDynamicData: dynamicData,
	}
}

func JoinItem(staticData *ItemStaticData, dynamicData ItemDynamicData) *Item {
	return &Item{
		ItemStaticData:  staticData,
		ItemDynamicData: dynamicData,
	}
}
