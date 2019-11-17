package business

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;not null;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (base *Base) BeforeCreate(scope *gorm.Scope) error {
	id := uuid.NewV4()
	return scope.SetColumn("ID", id)
}

type Employee struct {
	Base
	Name       string `gorm:"not null;"`
	Categories []*Category
}

type Category struct {
	Base
	Name string `json:"name,omitempty" gorm:"unique;not null;"`
}

type Gift struct {
	Base
	Name       string `json:"name,omitempty" gorm:"unique;not null;"`
	Categories []*Category
}
