package api

import (
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

const PaymentsPath = "/payments"

// GetMux generates the app router
func (s *Service) GetMux() *chi.Mux {
	r := chi.NewRouter()

	r.Use(
		middleware.RequestID, // this can be passed further downstream through microservices and then (open)traced
		middleware.RealIP,
		middleware.Recoverer,
		middleware.RedirectSlashes,
		// timeout for processing a request:
		middleware.Timeout(20*time.Second),
		middleware.AllowContentType("application/json"),          // set json the only valid request content type
		middleware.SetHeader("Content-Type", "application/json"), // respond with json
		middleware.SetHeader("Access-Control-Allow-Origin", "*"), // CORS
		render.SetContentType(render.ContentTypeJSON),            // default renderer
		// use logrus for logging
		middleware.RequestLogger(&StructuredLogger{s.logger}),
	)

	r.Mount(PaymentsPath, getPaymentsRouter(s))

	return r
}

func getPaymentsRouter(s *Service) *chi.Mux {
	r := chi.NewRouter()
	accessCORSLocationMiddleware := middleware.SetHeader("Access-Control-Expose-Headers", "Location")

	r.Use(accessCORSLocationMiddleware)
	r.Post("/", s.Create)
	r.Get("/", s.List) // todo add a paginator middleware for this route

	r.Route("/{paymentID}", func(r chi.Router) {
		r.Use(s.singlePaymentCtx) // Load the *Payment on the request context or return 404 if paymentID is not found
		r.Get("/", s.Retrieve)
		r.Put("/", s.Update)
		r.Delete("/", s.Delete)
	})

	return r
}
