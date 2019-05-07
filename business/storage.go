package business

// Storage is the persistence provider interface. It needs to be implemented for every persistence provider.
type Storage interface {
	// Create returns the created object id as a slice of bytes
	Create(payment Payment) (*Payment, error)
	// Retrieve returns nil if the object was not found
	Retrieve(id []byte) (*Payment, error)
	// Update updates the payment based on the object id and its associated data
	Update(payment Payment) error
	// Delete deletes the object by its id
	Delete(id []byte) error
	// List returns the list of payments (todo ordering, sorting, pagination)
	List() ([]*Payment, error)
}
