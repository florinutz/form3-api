package memory

import (
	"errors"

	"form3/business"

	"gopkg.in/mgo.v2/bson"
)

// Storage stores payments in memory
type Storage map[string]*business.Payment

// Create creates a new payment
func (c Storage) Create(payment business.Payment) (*business.Payment, error) {
	if payment.Id == "" {
		payment.Id = bson.NewObjectId()
	}

	idStr := payment.Id.Hex()
	if _, exists := c[idStr]; exists {
		return nil, errors.New("payment already exists")
	}

	c[idStr] = &payment

	return &payment, nil
}

// Retrieve retrieves a payment by its ID
func (c Storage) Retrieve(id []byte) (payment *business.Payment, err error) {
	var exists bool
	if payment, exists = c[string(id)]; !exists {
		err = errors.New("payment doesn't exist")
	}

	return
}

// Update updates a payment
func (c Storage) Update(payment business.Payment) error {
	idStr := payment.Id.Hex()

	if _, exists := c[idStr]; !exists {
		return errors.New("payment doesn't already exist")
	}

	c[idStr] = &payment

	return nil
}

// Delete deletes a payment
func (c Storage) Delete(id []byte) error {
	idStr := string(id)

	var exists bool
	if _, exists = c[idStr]; !exists {
		return errors.New("payment doesn't exist")
	}

	delete(c, idStr)

	return nil
}

// List lists all payments
func (c Storage) List() (payments []*business.Payment, err error) {
	for _, payment := range c {
		payments = append(payments, payment)
	}

	return
}
