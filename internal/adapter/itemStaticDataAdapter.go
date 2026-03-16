package adapter

import (
	domain "github.com/L1irik259/TestForOzon/internal/domain"

	"gorm.io/gorm"
)

type ItemStaticDataAdapter struct {
	db *gorm.DB
}

func NewItemStaticDataAdapter(db *gorm.DB) *ItemStaticDataAdapter {
	return &ItemStaticDataAdapter{db: db}
}

func (a *ItemStaticDataAdapter) GetAllItems() ([]*domain.ItemStaticData, error) {
	var items []*domain.ItemStaticData

	if err := a.db.Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (a *ItemStaticDataAdapter) GetAllItemsByDB(tx *gorm.DB) ([]*domain.ItemStaticData, error) {
	var items []*domain.ItemStaticData

	if err := tx.Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (a *ItemStaticDataAdapter) GetItemByIDByDB(tx *gorm.DB, id string) (*domain.ItemStaticData, error) {
	var item domain.ItemStaticData

	if err := tx.First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (a *ItemStaticDataAdapter) CreateItemByDB(tx *gorm.DB, item *domain.ItemStaticData) error {
	return tx.Create(item).Error
}

func (a *ItemStaticDataAdapter) UpdateItemByDB(tx *gorm.DB, item *domain.ItemStaticData) error {
	return tx.Save(item).Error
}

func (a *ItemStaticDataAdapter) DeleteItemByDB(tx *gorm.DB, id string) error {
	return tx.Delete(&domain.ItemStaticData{}, "id = ?", id).Error
}

func (a *ItemStaticDataAdapter) GetItemByID(id string) (*domain.ItemStaticData, error) {
	return a.GetItemByIDByDB(a.db, id)
}

func (a *ItemStaticDataAdapter) CreateItem(item *domain.ItemStaticData) error {
	return a.CreateItemByDB(a.db, item)
}

func (a *ItemStaticDataAdapter) UpdateItem(item *domain.ItemStaticData) error {
	return a.UpdateItemByDB(a.db, item)
}

func (a *ItemStaticDataAdapter) DeleteItem(id string) error {
	return a.DeleteItemByDB(a.db, id)
}
