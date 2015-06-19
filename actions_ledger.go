package horizon

import (
	"net/http"

	"github.com/stellar/go-horizon/db"
	"github.com/stellar/go-horizon/render"
	"github.com/stellar/go-horizon/render/hal"
	"github.com/stellar/go-horizon/render/problem"
	"github.com/stellar/go-horizon/render/sse"
	"github.com/zenazn/goji/web"
)

func ledgerIndexAction(c web.C, w http.ResponseWriter, r *http.Request) {
	ah := &ActionHelper{c: c, r: r}
	ctx := ah.Context()

	// construct query
	ah.ValidateInt64(ParamCursor)
	query := db.LedgerPageQuery{
		SqlQuery:  ah.App().HistoryQuery(),
		PageQuery: ah.GetPageQuery(),
	}

	if ah.Err() != nil {
		problem.Render(ctx, w, ah.Err())
		return
	}

	var records []db.LedgerRecord

	render := render.Renderer{}
	render.JSON = func() {
		// load records
		err := db.Select(ctx, query, &records)
		if err != nil {
			problem.Render(ctx, w, err)
		}

		page, err := NewLedgerPageResource(records, query.PageQuery)
		if err != nil {
			problem.Render(ctx, w, err)
		}

		hal.RenderPage(w, page)
	}

	render.SSE = func(stream sse.Stream) {
		err := db.Select(ctx, query, &records)
		if err != nil {
			stream.Err(err)
			return
		}

		records = records[stream.SentCount():]

		for _, record := range records {
			stream.Send(sse.Event{
				ID:   record.PagingToken(),
				Data: NewLedgerResource(record),
			})
		}

		if stream.SentCount() >= int(query.Limit) {
			stream.Done()
		}
	}

	render.Render(ctx, w, r)
}

func ledgerShowAction(c web.C, w http.ResponseWriter, r *http.Request) {
	ah := &ActionHelper{c: c, r: r}
	app := ah.App()

	// construct query
	query := db.LedgerBySequenceQuery{
		SqlQuery: app.HistoryQuery(),
		Sequence: ah.GetInt32("id"),
	}
	if ah.Err() != nil {
		problem.Render(ah.Context(), w, ah.Err())
		return
	}

	// find ledger
	found, err := db.First(ah.Context(), query)
	if err != nil {
		problem.Render(ah.Context(), w, ah.Err())
		return
	}

	ledger, ok := found.(db.LedgerRecord)
	if !ok {
		problem.Render(ah.Context(), w, problem.NotFound)
		return
	}

	render := render.Renderer{}
	render.JSON = func() {
		hal.Render(w, NewLedgerResource(ledger))
	}

	render.Render(ah.Context(), w, r)
}
