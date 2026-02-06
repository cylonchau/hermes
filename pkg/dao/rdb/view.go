package rdb

import (
	"context"

	"gorm.io/gorm"

	"github.com/cylonchau/hermes/pkg/model"
)

// ViewDAO View的关系型数据访问层
type ViewDAO struct {
	db *gorm.DB
}

// NewViewDAO 创建ViewDAO实例
func NewViewDAO(db *gorm.DB) *ViewDAO {
	return &ViewDAO{db: db}
}

// Create 创建View
func (dao *ViewDAO) Create(ctx context.Context, view *model.View) error {
	return dao.db.WithContext(ctx).Create(view).Error
}

// GetByID 根据ID获取View
func (dao *ViewDAO) GetByID(ctx context.Context, id int64) (*model.View, error) {
	var view model.View
	err := dao.db.WithContext(ctx).First(&view, id).Error
	if err != nil {
		return nil, err
	}
	return &view, nil
}

// GetByName 根据名称获取View
func (dao *ViewDAO) GetByName(ctx context.Context, name string) (*model.View, error) {
	var view model.View
	err := dao.db.WithContext(ctx).Where("name = ?", name).First(&view).Error
	if err != nil {
		return nil, err
	}
	return &view, nil
}

// GetAll 获取所有View
func (dao *ViewDAO) GetAll(ctx context.Context) ([]*model.View, error) {
	var views []*model.View
	err := dao.db.WithContext(ctx).Find(&views).Error
	return views, err
}

// Update 更新View
func (dao *ViewDAO) Update(ctx context.Context, view *model.View) error {
	return dao.db.WithContext(ctx).Save(view).Error
}

// Delete 删除View
func (dao *ViewDAO) Delete(ctx context.Context, id int64) error {
	return dao.db.WithContext(ctx).Delete(&model.View{}, id).Error
}
