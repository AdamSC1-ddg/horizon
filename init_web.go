package horizon

import (
	"github.com/PuerkitoBio/throttled"
	"github.com/PuerkitoBio/throttled/store"
	"github.com/rcrowley/go-metrics"
	"github.com/rs/cors"
	"github.com/sebest/xff"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Web struct {
	router      *web.Mux
	rateLimiter *throttled.Throttler

	requestTimer metrics.Timer
	failureMeter metrics.Meter
	successMeter metrics.Meter
}

func initWeb(app *App) {
	app.web = &Web{
		router:       web.New(),
		requestTimer: metrics.NewTimer(),
		failureMeter: metrics.NewMeter(),
		successMeter: metrics.NewMeter(),
	}
}

func initWebMetrics(app *App) {
	app.metrics.Register("requests.total", app.web.requestTimer)
	app.metrics.Register("requests.succeeded", app.web.successMeter)
	app.metrics.Register("requests.failed", app.web.failureMeter)
}

func initWebMiddleware(app *App) {
	app.web.router.Use(middleware.EnvInit)
	app.web.router.Use(middleware.RequestID)
	app.web.router.Use(xff.XFF)
	app.web.router.Use(app.Middleware)
	app.web.router.Use(middleware.Logger)
	app.web.router.Use(RecoverMiddleware)
	app.web.router.Use(middleware.AutomaticOptions)
	app.web.router.Use(requestMetricsMiddleware)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})
	app.web.router.Use(c.Handler)

	app.web.router.Use(app.web.RateLimitMiddleware)
}

func initWebActions(app *App) {
	app.web.router.Get("/", rootAction)
	app.web.router.Get("/metrics", metricsAction)

	// ledger actions
	app.web.router.Get("/ledgers", ledgerIndexAction)
	app.web.router.Get("/ledgers/:id", ledgerShowAction)
	app.web.router.Get("/ledgers/:ledger_id/transactions", transactionIndexAction)
	app.web.router.Get("/ledgers/:ledger_id/operations", operationIndexAction)
	app.web.router.Get("/ledgers/:ledger_id/payments", notImplementedAction)
	app.web.router.Get("/ledgers/:ledger_id/effects", notImplementedAction)

	// account actions
	app.web.router.Get("/accounts", notImplementedAction)
	app.web.router.Get("/accounts/:id", accountShowAction)
	app.web.router.Get("/accounts/:account_id/transactions", transactionIndexAction)
	app.web.router.Get("/accounts/:account_id/operations", operationIndexAction)
	app.web.router.Get("/accounts/:account_id/payments", notImplementedAction)
	app.web.router.Get("/accounts/:account_id/effects", notImplementedAction)

	// transaction actions
	app.web.router.Get("/transactions", transactionIndexAction)
	app.web.router.Get("/transactions/:id", transactionShowAction)
	app.web.router.Get("/transactions/:tx_id/operations", operationIndexAction)
	app.web.router.Get("/transactions/:tx_id/payments", notImplementedAction)
	app.web.router.Get("/transactions/:tx_id/effects", notImplementedAction)

	// operation actions
	app.web.router.Get("/operations", operationIndexAction)
	app.web.router.Get("/operations/:id", notImplementedAction)
	app.web.router.Get("/operations/:op_id/effects", notImplementedAction)

	app.web.router.Get("/payments", notImplementedAction)

	// go-horizon doesn't implement everything horizon did,
	// so we reverse proxy if we can
	if app.config.RubyHorizonUrl != "" {

		u, err := url.Parse(app.config.RubyHorizonUrl)
		if err != nil {
			panic("cannot parse ruby-horizon-url")
		}

		rp := httputil.NewSingleHostReverseProxy(u)
		app.web.router.Post("/transactions", rp)
		app.web.router.Post("/friendbot", rp)
		app.web.router.Get("/friendbot", rp)
	} else {
		app.web.router.Post("/transactions", notImplementedAction)
		app.web.router.Post("/friendbot", notImplementedAction)
		app.web.router.Get("/friendbot", notImplementedAction)
	}

	app.web.router.NotFound(notFoundAction)
}

func initWebRateLimiter(app *App) {
	rateLimitStore := store.NewMemStore(1000)

	if app.redis != nil {
		rateLimitStore = store.NewRedisStore(app.redis, "throttle:", 0)
	}

	rateLimiter := throttled.RateLimit(
		app.config.RateLimit,
		&throttled.VaryBy{Custom: remoteAddrIp},
		rateLimitStore,
	)
	rateLimiter.DeniedHandler = http.HandlerFunc(rateLimitExceededAction)
	app.web.rateLimiter = rateLimiter
}

func remoteAddrIp(r *http.Request) string {
	ip := strings.SplitN(r.RemoteAddr, ":", 2)[0]
	return ip
}
