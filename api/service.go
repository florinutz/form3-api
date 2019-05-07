package api

import (
	"context"
	"net/http"

	"form3/business"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

// Service exposes http handlers for the api service's operations
type Service struct {
	persistence business.Storage
	logger      *logrus.Logger
}

// NewService instantiates the payment Storage
func NewService(persistence business.Storage, logger logrus.Logger) *Service {
	return &Service{persistence, &logger}
}

// List is the handler for the list request
func (s *Service) List(w http.ResponseWriter, r *http.Request) {
	payments, err := s.persistence.List()
	if err != nil {
		render.Render(w, r, ErrInternalError(err))
		return
	}

	if err := render.RenderList(w, r, NewPaymentListResponse(payments)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// Create is the handler responsible with creating a payment
func (s *Service) Create(w http.ResponseWriter, r *http.Request) {
	var (
		paymentCreateReq PaymentCreateRequest
		payment          *business.Payment
		err              error
	)

	if err = render.Bind(r, &paymentCreateReq); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// the nil pointer check is done in the Bind() call above
	if payment, err = s.persistence.Create(*paymentCreateReq.Payment); err != nil {
		render.Render(w, r, ErrInternalError(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &PaymentCreatedResponse{*payment})
}

// Retrieve is the handler responsible with retrieving a payment. It requires an id as input.
func (s *Service) Retrieve(w http.ResponseWriter, r *http.Request) {
	payment := r.Context().Value(contextKeyPayment("payment")).(*business.Payment)

	responseRepresentation := &PaymentRetrieveResponse{*payment}
	if err := render.Render(w, r, responseRepresentation); err != nil {
		render.Status(r, http.StatusUnprocessableEntity)
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Status(r, http.StatusOK)
}

// Update is the handler responsible with updating a payment.
// It requires a payment as input.
func (s *Service) Update(w http.ResponseWriter, r *http.Request) {
	var (
		paymentUpdateReq PaymentCreateRequest
		err              error
	)

	if err = render.Bind(r, &paymentUpdateReq); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// the nil pointer check is done in the Bind() call above
	if err = s.persistence.Update(*paymentUpdateReq.Payment); err != nil {
		render.Render(w, r, ErrInternalError(err))
		return
	}

	render.Status(r, http.StatusAccepted)
	render.Render(w, r, &PaymentCreatedResponse{*paymentUpdateReq.Payment})
}

// Delete is the handler responsible with deleting a payment.
// It requires the paymentID as input.
func (s *Service) Delete(w http.ResponseWriter, r *http.Request) {
	payment := r.Context().Value(contextKeyPayment("payment")).(*business.Payment)

	if err := s.persistence.Delete([]byte(payment.Id.Hex())); err != nil {
		render.Render(w, r, ErrInternalError(err))
		return
	}

	render.Status(r, http.StatusAccepted)
}

type contextKeyPayment string

// singlePaymentCtx is a middleware used by CRUD's R, U and D.
// It's used for loading a Payment object from the paymentID URL parameter passed along with the Request.
// In case the Payment could not be found, we stop here and return a 404.
func (s *Service) singlePaymentCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			payment *business.Payment
			err     error
		)

		paymentID := chi.URLParam(r, "paymentID")
		if paymentID == "" {
			render.Render(w, r, ErrNotFound)
			return
		}

		if payment, err = s.persistence.Retrieve([]byte(paymentID)); err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyPayment("payment"), payment)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
