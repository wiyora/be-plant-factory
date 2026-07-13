package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
)

type Tenant struct {
	ID        uuid.UUID           `db:"id"`
	Name      string              `db:"name"`
	Logo      string              `db:"logo"`
	Status    entity.TenantStatus `db:"status"`
	CreatedAt time.Time           `db:"created_at"`
	UpdatedAt *time.Time          `db:"updated_at"`
}

func (t Tenant) ToEntity() entity.Tenant {
	return entity.Tenant{
		ID:        t.ID,
		Name:      t.Name,
		Logo:      t.Logo,
		Status:    t.Status,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

type ListTenant struct {
	ID        uuid.UUID           `db:"id"`
	Name      string              `db:"name"`
	Logo      string              `db:"logo"`
	Status    entity.TenantStatus `db:"status"`
	CreatedAt time.Time           `db:"created_at"`
}

func (t ListTenant) ToEntity() entity.Tenant {
	return entity.Tenant{
		ID:        t.ID,
		Name:      t.Name,
		Logo:      t.Logo,
		Status:    t.Status,
		CreatedAt: t.CreatedAt,
	}
}
