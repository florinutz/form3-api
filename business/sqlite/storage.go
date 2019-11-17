package sqlite

import (
	"fmt"

	"form3/business"
	"form3/importer"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	uuid "github.com/satori/go.uuid"
)

// Storage implements the Storage interface
type Storage struct {
	db *gorm.DB
}

// New is the Storage constructor. dbLocation can be a file path or ":memory:"
func New(dbLocation string) (*Storage, error) {
	db, err := gorm.Open("sqlite3", dbLocation)
	if err != nil {
		return nil, err
	}

	db.LogMode(true)
	db.AutoMigrate(&business.Category{}, &business.Employee{}, &business.Gift{})

	return &Storage{db}, nil
}

// AttachGift implements the Storage interface's AttachGift
func (s *Storage) AttachGift(id uuid.UUID) (bool, error) {
	return true, nil
}

// Retrieve implements the Storage interface's Retrieve
func (s *Storage) Retrieve(id []byte) (*business.Employee, error) {
	return nil, nil
}

func (s *Storage) ImportData() error {
	categories, employees, gifts, err := importer.Fetch()
	if err != nil {
		return fmt.Errorf("errored while importing data: %w", err)
	}

	for _, c := range categories {
		if err := s.db.Create(c).Error; err != nil {
			return err
		}
	}

	for _, e := range employees {
		if err := s.db.Create(e).Error; err != nil {
			return err
		}
	}

	for _, g := range gifts {
		if err := s.db.Create(g).Error; err != nil {
			return err
		}
	}

	return nil
}
