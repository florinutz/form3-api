package business

import uuid "github.com/satori/go.uuid"

// Storage is the persistence provider interface. It needs to be implemented for every persistence provider.
type Storage interface {
	// AttachGift attaches a gift to an employee
	AttachGift(id uuid.UUID) (bool, error)
	// Retrieve returns nil if the object was not found
	Retrieve(id []byte) (*Employee, error)
	// ImportData migrates data into the storage
	ImportData() error
}
