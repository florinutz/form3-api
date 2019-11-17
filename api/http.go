package api

import (
	"net/http"

	"form3/business"

	uuid "github.com/satori/go.uuid"
)

// CreateRequest represents a request to attach a gift to an employee
type CreateRequest struct {
	Id uuid.UUID
}

// Bind implementing render's Binder
func (cr *CreateRequest) Bind(r *http.Request) error {
	// post process? incoming data?
	return nil
}

// CreatedResponse represents a response for when a payment was successfully created
type CreatedResponse struct {
	ok bool
}

// Render implements render's Renderer interface
func (pcr *CreatedResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// RetrieveResponse represents a retrieve response
type RetrieveResponse struct {
	business.Employee
}

// Render implements Renderer
func (pr *RetrieveResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
