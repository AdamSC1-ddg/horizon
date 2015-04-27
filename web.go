package horizon

import (
	"github.com/rcrowley/go-metrics"
	"github.com/rs/cors"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

type Web struct {
	router *web.Mux

	requestTimer metrics.Timer
	failureMeter metrics.Meter
	successMeter metrics.Meter
}

func NewWeb(app *App) (*Web, error) {

	api := web.New()
	result := Web{
		router:       api,
		requestTimer: metrics.NewTimer(),
		failureMeter: metrics.NewMeter(),
		successMeter: metrics.NewMeter(),
	}

	app.metrics.Register("requests.total", result.requestTimer)
	app.metrics.Register("requests.succeeded", result.successMeter)
	app.metrics.Register("requests.failed", result.failureMeter)

	api.Use(middleware.RequestID)
	api.Use(middleware.Logger)
	api.Use(middleware.Recoverer)
	api.Use(middleware.AutomaticOptions)
	api.Use(appMiddleware(app))
	api.Use(requestMetricsMiddleware)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})
	api.Use(c.Handler)

	// define routes
	api.Get("/", rootAction)
	api.Get("/metrics", metricsAction)
	api.Get("/ledgers", ledgerIndexAction)
	api.Get("/ledgers/:id", ledgerShowAction)

	return &result, nil
}
