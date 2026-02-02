package rdb

import (
	"context"

	"gorm.io/gorm"

	"github.com/cylonchau/hermes/pkg/model"
)

// ZoneDAO Zone的关系型数据访问层 - 目标支持MySQL/PostgreSQL/SQLite等关系型数据库
type ZoneDAO struct {
	db *gorm.DB
}

// NewZoneDAO 创建ZoneDAO实例
func NewZoneDAO(db *gorm.DB) *ZoneDAO {
	db.AutoMigrate(&model.Zone{})
	return &ZoneDAO{db: db}
}

// Create 创建Zone
func (dao *ZoneDAO) Create(ctx context.Context, zone *model.Zone) error {
	return dao.db.WithContext(ctx).Create(zone).Error
}

// GetByID 根据ID查询Zone
func (dao *ZoneDAO) GetByID(ctx context.Context, id int64) (*model.Zone, error) {
	var zone model.Zone
	err := dao.db.WithContext(ctx).Where("id = ? AND is_active = ?", id, true).First(&zone).Error
	if err != nil {
		return nil, err
	}
	return &zone, nil
}

// GetByName 根据名称查询Zone
func (dao *ZoneDAO) GetByName(ctx context.Context, name string) (*model.Zone, error) {
	var zone model.Zone
	err := dao.db.WithContext(ctx).Where("name = ? AND is_active = ?", name, true).First(&zone).Error
	if err != nil {
		return nil, err
	}
	return &zone, nil
}

// GetAll 查询所有Zone
func (dao *ZoneDAO) GetAll(ctx context.Context, limit, offset int) ([]*model.Zone, error) {
	var zones []*model.Zone
	query := dao.db.WithContext(ctx).Where("is_active = ?", true)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&zones).Error
	return zones, err
}

// Update 更新Zone
func (dao *ZoneDAO) Update(ctx context.Context, zone *model.Zone) error {
	return dao.db.WithContext(ctx).Save(zone).Error
}

// Delete 物理删除Zone
func (dao *ZoneDAO) Delete(ctx context.Context, id int64) error {
	return dao.db.WithContext(ctx).Delete(&model.Zone{}, id).Error
}

// SoftDelete 软删除Zone
func (dao *ZoneDAO) SoftDelete(ctx context.Context, id int64) error {
	return dao.db.WithContext(ctx).Model(&model.Zone{}).Where("id = ?", id).Update("is_active", false).Error
}

// Count 统计Zone数量
func (dao *ZoneDAO) Count(ctx context.Context) (int64, error) {
	var count int64
	err := dao.db.WithContext(ctx).Model(&model.Zone{}).Where("is_active = ?", true).Count(&count).Error
	return count, err
}

// ExistsByName 检查Zone名称是否存在
func (dao *ZoneDAO) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	err := dao.db.WithContext(ctx).Model(&model.Zone{}).Where("name = ? AND is_active = ?", name, true).Count(&count).Error
	return count > 0, err
}

// GetActiveZones 查询活跃的Zone
func (dao *ZoneDAO) GetActiveZones(ctx context.Context) ([]*model.Zone, error) {
	var zones []*model.Zone
	err := dao.db.WithContext(ctx).Where("is_active = ?", true).Find(&zones).Error
	return zones, err
}

// GetByContact 根据联系人查询Zone
func (dao *ZoneDAO) GetByContact(ctx context.Context, contact string) ([]*model.Zone, error) {
	var zones []*model.Zone
	err := dao.db.WithContext(ctx).Where("contact = ? AND is_active = ?", contact, true).Find(&zones).Error
	return zones, err
}

// GetByEmail 根据邮箱查询Zone
func (dao *ZoneDAO) GetByEmail(ctx context.Context, email string) ([]*model.Zone, error) {
	var zones []*model.Zone
	err := dao.db.WithContext(ctx).Where("email = ? AND is_active = ?", email, true).Find(&zones).Error
	return zones, err
}

// Search 搜索Zone
func (dao *ZoneDAO) Search(ctx context.Context, keyword string, limit, offset int) ([]*model.Zone, error) {
	var zones []*model.Zone
	query := dao.db.WithContext(ctx).Where("is_active = ?", true)

	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ? OR contact LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&zones).Error
	return zones, err
}

// BatchCreate 批量创建Zone
func (dao *ZoneDAO) BatchCreate(ctx context.Context, zones []*model.Zone) error {
	return dao.db.WithContext(ctx).CreateInBatches(zones, 100).Error
}

// BatchUpdate 批量更新Zone
func (dao *ZoneDAO) BatchUpdate(ctx context.Context, zones []*model.Zone) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, zone := range zones {
			if err := tx.Save(zone).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// BatchDelete 批量删除Zone
func (dao *ZoneDAO) BatchDelete(ctx context.Context, ids []int64) error {
	return dao.db.WithContext(ctx).Delete(&model.Zone{}, ids).Error
}
