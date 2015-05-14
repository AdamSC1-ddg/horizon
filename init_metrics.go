package horizon

import (
	"github.com/rcrowley/go-metrics"
	"github.com/stellar/go-horizon/db"
)

func initMetrics(app *App) {
	app.metrics = metrics.NewRegistry()
}

func initDbMetrics(app *App) {
	app.metrics.Register("db.active_query_count", db.QueryGauge())
	app.metrics.Register("db.active_query_count", db.QueryGauge())
	app.metrics.Register("history.latest_ledger", db.HorizonLedgerGauge())
	app.metrics.Register("stellar_core.latest_ledger", db.StellarCoreLedgerGauge())
}
