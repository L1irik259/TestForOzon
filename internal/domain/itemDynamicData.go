package domain

import (
	"time"
)

type ItemDynamicData struct {
	ID        string
	Nominal   int
	Value     float64
	VunitRate float64
	Date      time.Time
}

func NewItemDynamicData(id string, nominal int, value float64, date time.Time) *ItemDynamicData {
	return &ItemDynamicData{
		ID:        id,
		Nominal:   nominal,
		Value:     value,
		VunitRate: value / float64(nominal),
		Date:      date,
	}
}

func NewItemDynamicDataByIdByDate(id string, date time.Time) *ItemDynamicData {
	return &ItemDynamicData{
		ID:        id,
		Nominal:   0,
		Value:     0,
		VunitRate: 0,
		Date:      date,
	}
}
