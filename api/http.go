package api

import (
	"errors"
	"net/http"

	"form3/business"

	"github.com/go-chi/render"
)

// PaymentCreateRequest represents a request to create a payment
type PaymentCreateRequest struct {
	*business.Payment
}

// Bind implementing render's Binder
func (cpr *PaymentCreateRequest) Bind(r *http.Request) error {
	// avoid a nil pointer dereference:
	if cpr.Payment == nil {
		return errors.New("missing required Payment fields")
	}
	// post process? incoming data?
	return nil
}

// PaymentCreatedResponse represents a response for when a payment was successfully created
// todo stop using this for update
type PaymentCreatedResponse struct {
	business.Payment
}

// Render implements render's Renderer interface
func (pcr *PaymentCreatedResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// PaymentRetrieveResponse represents a retrieve response
type PaymentRetrieveResponse struct {
	business.Payment
}

// Render implements Renderer
func (pr *PaymentRetrieveResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// NewPaymentListResponse is the constructor for the list renderer
func NewPaymentListResponse(payments []*business.Payment) (list []render.Renderer) {
	for _, payment := range payments {
		list = append(list, &PaymentRetrieveResponse{*payment})
	}
	return
}
