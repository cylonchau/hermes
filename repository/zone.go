package repository

import (
	"context"

	"github.com/cylonchau/hermes/pkg/model"
)

// ZoneRepository Zone仓储接口 - 定义Zone数据存储的抽象接口
type ZoneRepository interface {
	// 基本CRUD操作
	Create(ctx context.Context, zone *model.Zone) error
	GetByID(ctx context.Context, id uint) (*model.Zone, error)
	GetByName(ctx context.Context, name string) (*model.Zone, error)
	GetAll(ctx context.Context, limit, offset int) ([]*model.Zone, error)
	Update(ctx context.Context, zone *model.Zone) error
	Delete(ctx context.Context, id uint) error
	SoftDelete(ctx context.Context, id uint) error // 软删除

	// 查询操作
	Count(ctx context.Context) (int64, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
	GetActiveZones(ctx context.Context) ([]*model.Zone, error)
	GetByContact(ctx context.Context, contact string) ([]*model.Zone, error)
	GetByEmail(ctx context.Context, email string) ([]*model.Zone, error)
	Search(ctx context.Context, keyword string, limit, offset int) ([]*model.Zone, error)

	// 批量操作
	BatchCreate(ctx context.Context, zones []*model.Zone) error
	BatchUpdate(ctx context.Context, zones []*model.Zone) error
	BatchDelete(ctx context.Context, ids []uint) error
}
