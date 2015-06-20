package horizon

import (
	"net/http"

	"github.com/stellar/go-horizon/actions"
	"github.com/stellar/go-horizon/db"
	"github.com/stellar/go-horizon/render/hal"
	"github.com/stellar/go-horizon/render/sse"
	"github.com/zenazn/goji/web"
)

type PaymentsIndexAction struct {
	Action
	Query   db.OperationPageQuery
	Records []db.OperationRecord
	Page    hal.Page
}

// ServeHTTPC is a method for web.Handler
func (action PaymentsIndexAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	ap.Execute(&action)
}

// LoadQuery sets action.Query from the request params
func (action *PaymentsIndexAction) LoadQuery() {
	action.ValidateInt64(actions.ParamCursor)
	action.Query = db.OperationPageQuery{
		SqlQuery:        action.App.HistoryQuery(),
		PageQuery:       action.GetPageQuery(),
		AccountAddress:  action.GetString("account_id"),
		LedgerSequence:  action.GetInt32("ledger_id"),
		TransactionHash: action.GetString("tx_id"),
		TypeFilter:      db.PaymentTypeFilter,
	}
}

// LoadRecords populates action.Records
func (action *PaymentsIndexAction) LoadRecords() {
	action.LoadQuery()
	if action.Err != nil {
		return
	}

	action.Err = db.Select(action.Ctx, action.Query, &action.Records)
}

// LoadPage populates action.Page
func (action *PaymentsIndexAction) LoadPage() {
	action.LoadRecords()
	if action.Err != nil {
		return
	}

	action.Page, action.Err = NewOperationResourcePage(action.Records, action.Query.PageQuery, "")
}

// JSON is a method for actions.JSON
func (action *PaymentsIndexAction) JSON() {
	action.LoadPage()
	if action.Err != nil {
		return
	}
	hal.Render(action.W, action.Page)
}

// SSE is a method for actions.SSE
func (action *PaymentsIndexAction) SSE(stream sse.Stream) {
	action.LoadRecords()

	if action.Err != nil {
		stream.Err(action.Err)
		return
	}

	records := action.Records[stream.SentCount():]

	for _, record := range records {
		r, err := NewOperationResource(record)

		if err != nil {
			stream.Err(action.Err)
			return
		}

		stream.Send(sse.Event{
			ID:   record.PagingToken(),
			Data: r,
		})
	}

	if stream.SentCount() >= int(action.Query.Limit) {
		stream.Done()
	}
}
