package api

import (
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

// NewService instantiates the employee Storage
// Requires a logger.
func NewService(persistence business.Storage, logger logrus.Logger) *Service {
	return &Service{persistence, &logger}
}

// Create is the handler responsible with attaching a gift to an employee.
// Employee is referred by the uuid.
// The gift is attached only once.
func (s *Service) Create(w http.ResponseWriter, r *http.Request) {
	var (
		createReq CreateRequest
		ok        bool
		err       error
	)

	if err = render.Bind(r, &createReq); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if ok, err = s.persistence.AttachGift(createReq.Id); err != nil {
		render.Render(w, r, ErrInternalError(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &CreatedResponse{ok})
}

// Retrieve is the handler responsible with retrieving an object. It requires an id as input.
func (s *Service) Retrieve(w http.ResponseWriter, r *http.Request) {
	var (
		obj *business.Employee
		err error
	)

	var id string
	if id = chi.URLParam(r, "id"); id == "" {
		render.Render(w, r, ErrNotFound)
		return
	}

	if obj, err = s.persistence.Retrieve([]byte(id)); err != nil {
		render.Render(w, r, ErrNotFound)
		return
	}

	responseRepresentation := &RetrieveResponse{*obj}
	if err := render.Render(w, r, responseRepresentation); err != nil {
		render.Status(r, http.StatusUnprocessableEntity)
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Status(r, http.StatusOK)
}
