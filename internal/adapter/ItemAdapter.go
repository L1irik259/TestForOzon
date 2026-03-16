package adapter

import (
	"errors"
	"time"

	"gorm.io/gorm"

	domain "github.com/L1irik259/TestForOzon/internal/domain"
)

// Оптимизировать метод GetAllItemsByDate
type ItemAdapter struct {
	db *gorm.DB

	staticDataAdapter  *ItemStaticDataAdapter
	dynamicDataAdapter *ItemDynamicDataAdapter
}

func NewItemAdapter(
	db *gorm.DB,
) *ItemAdapter {
	return &ItemAdapter{
		db: db,

		staticDataAdapter:  NewItemStaticDataAdapter(db),
		dynamicDataAdapter: NewItemDynamicDataAdapter(db),
	}
}

func (a *ItemAdapter) GetItemByIDByDate(
	id string,
	date time.Time,
) (*domain.Item, error) {
	var result *domain.Item

	err := a.db.Transaction(func(tx *gorm.DB) error {
		staticData, err := a.staticDataAdapter.GetItemByIDByDB(tx, id)
		if err != nil {
			return err
		}

		dynamicData, err := a.dynamicDataAdapter.GetItemByIDByDBByData(tx, id, date)
		if err != nil {
			return err
		}

		item := domain.JoinItem(staticData, *dynamicData)
		result = item

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// Надо будет оптимизировать, так как сейчас делается n+1 запрос, нежели 2 запросов
// и хранения данных используя мапы
func (a *ItemAdapter) GetAllItemsByDate(date time.Time) ([]*domain.Item, error) {
	var results []*domain.Item

	err := a.db.Transaction(func(tx *gorm.DB) error {
		staticDataList, err := a.staticDataAdapter.GetAllItemsByDB(tx)

		if err != nil {
			return err
		}

		dynamicDataMap, err := a.dynamicDataAdapter.GetItemByDateByDB(tx, date)
		if err != nil {
			return err
		}

		for _, staticData := range staticDataList {
			dynamicData, ok := dynamicDataMap[staticData.ID]
			if !ok {
				return errors.New("Нет динамических данных c нужным ID и нужным Date")
			}

			item := domain.JoinItem(staticData, *dynamicData)
			results = append(results, item)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (a *ItemAdapter) CreateItem(
	itemStaticData *domain.ItemStaticData,
	itemDynamicData *domain.ItemDynamicData,
) error {
	if err := a.db.Transaction(func(tx *gorm.DB) error {
		if err := a.staticDataAdapter.CreateItemByDB(tx, itemStaticData); err != nil {
			return err
		}

		if err := a.dynamicDataAdapter.CreateItemByDB(tx, itemDynamicData); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (a *ItemAdapter) CreateItemDynamicData(
	itemDynamicData *domain.ItemDynamicData,
) error {
	return a.dynamicDataAdapter.CreateItemByDB(a.db, itemDynamicData)
}

func (a *ItemAdapter) CreateItemStaticData(
	itemStaticData *domain.ItemStaticData,
) error {
	return a.staticDataAdapter.CreateItemByDB(a.db, itemStaticData)
}
