package adapter

import (
	"time"

	"gorm.io/gorm"

	domain "github.com/L1irik259/TestForOzon/internal/domain"
)

type ItemDynamicDataAdapter struct {
	db *gorm.DB
}

func NewItemDynamicDataAdapter(db *gorm.DB) *ItemDynamicDataAdapter {
	adapter := &ItemDynamicDataAdapter{
		db: db,
	}

	adapter.db.AutoMigrate(&domain.ItemDynamicData{})

	return adapter
}

func (a *ItemDynamicDataAdapter) GetAllItemsByDB(tx *gorm.DB) ([]*domain.ItemDynamicData, error) {
	var items []*domain.ItemDynamicData

	if err := tx.Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (a *ItemDynamicDataAdapter) GetItemByIDByDB(tx *gorm.DB, id string) (*domain.ItemDynamicData, error) {
	var item domain.ItemDynamicData

	if err := tx.First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (a *ItemDynamicDataAdapter) GetItemByIDByDBByData(tx *gorm.DB, id string, date time.Time) (*domain.ItemDynamicData, error) {
	var item domain.ItemDynamicData

	if err := tx.Where("id = ? AND date = ?", id, date).First(&item).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (a *ItemDynamicDataAdapter) GetItemByDateByDB(
	tx *gorm.DB,
	date time.Time,
) (map[string]*domain.ItemDynamicData, error) {
	var items []*domain.ItemDynamicData

	if err := tx.Where("date = ?", date).Find(&items).Error; err != nil {
		return nil, err
	}

	result := make(map[string]*domain.ItemDynamicData, len(items))

	for _, item := range items {
		result[item.ID] = item
	}

	return result, nil
}

func (a *ItemDynamicDataAdapter) CreateItemByDB(tx *gorm.DB, item *domain.ItemDynamicData) error {
	return tx.Create(item).Error
}

func (a *ItemDynamicDataAdapter) UpdateItemByDB(tx *gorm.DB, item *domain.ItemDynamicData) error {
	return tx.Save(item).Error
}

func (a *ItemDynamicDataAdapter) GetAllItems() ([]*domain.ItemDynamicData, error) {
	return a.GetAllItemsByDB(a.db)
}

func (a *ItemDynamicDataAdapter) GetItemByID(id string) (*domain.ItemDynamicData, error) {
	return a.GetItemByIDByDB(a.db, id)
}

func (a *ItemDynamicDataAdapter) GetItemByIDByDate(id string, date time.Time) (*domain.ItemDynamicData, error) {
	return a.GetItemByIDByDBByData(a.db, id, date)
}

func (a *ItemDynamicDataAdapter) CreateItem(item *domain.ItemDynamicData) error {
	return a.CreateItemByDB(a.db, item)
}

func (a *ItemDynamicDataAdapter) UpdateItem(item *domain.ItemDynamicData) error {
	return a.UpdateItemByDB(a.db, item)
}
