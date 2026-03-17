package service

import (
	"errors"
	"time"

	adapter "github.com/L1irik259/TestForOzon/internal/adapter"
	domain "github.com/L1irik259/TestForOzon/internal/domain"
)

type ItemService struct {
	itemAdapter *adapter.ItemAdapter
}

func NewItemService(
	itemAdapter *adapter.ItemAdapter,
) *ItemService {
	return &ItemService{
		itemAdapter: itemAdapter,
	}
}

func (s *ItemService) FindItemByIDByDate(
	id string,
	date time.Time,
) (*domain.Item, error) {

	d := date.UTC()

	if d.Hour() != 0 ||
		d.Minute() != 0 ||
		d.Second() != 0 ||
		d.Nanosecond() != 0 {
		return nil, errors.New("Неверный формат данных, нужно указать только нужный день.")
	}

	return s.itemAdapter.GetItemByIDByDate(id, date)
}

func (s *ItemService) NewItem(
	id, numCode, charCode, name string,
	nominal int,
	value, vunitRate float64,
	date time.Time,
) error {
	itemStatic := domain.NewItemStaticData(id, numCode, charCode, name)
	itemDynamic := domain.NewItemDynamicData(id, nominal, value, date)

	return s.itemAdapter.CreateItem(itemStatic, itemDynamic)
}

func (s *ItemService) NewItemDynamicData(
	id string,
	nominal int,
	value float64,
	date time.Time,
) error {
	d := date.UTC()

	if d.Hour() != 0 ||
		d.Minute() != 0 ||
		d.Second() != 0 ||
		d.Nanosecond() != 0 {
		return errors.New("Неверный формат данных, нужно указать только нужный день.")
	}

	itemDynamic := domain.NewItemDynamicData(id, nominal, value, date)

	return s.itemAdapter.CreateItemDynamicData(itemDynamic)
}

func (s *ItemService) NewItemStaticData(
	id, numCode, charCode, name string,
) error {
	if id == "" || numCode == "" || charCode == "" || name == "" {
		return errors.New("Неверный формат данных, все поля должны быть заполнены.")
	}

	itemStaticData := domain.NewItemStaticData(id, numCode, charCode, name)
	return s.itemAdapter.CreateItemStaticData(itemStaticData)
}

func (s *ItemService) FindAllItemsByDate(date time.Time) ([]*domain.Item, error) {
	d := date.UTC()

	if d.Hour() != 0 ||
		d.Minute() != 0 ||
		d.Second() != 0 ||
		d.Nanosecond() != 0 {
		return nil, errors.New("Неверный формат данных, нужно указать только нужный день.")
	}

	return s.itemAdapter.GetAllItemsByDate(date)
}
