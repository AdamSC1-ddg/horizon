package horizon

import (
	"fmt"
	"github.com/jagregory/halgo"
	"github.com/stellar/go-horizon/db"
	"github.com/stellar/go-horizon/render"
	"github.com/stellar/go-horizon/render/problem"
	"github.com/zenazn/goji/web"
	"golang.org/x/net/context"
	"math"
	"net/http"
	"time"
)

type LedgerResource struct {
	halgo.Links
	Id               string    `json:"id"`
	Hash             string    `json:"hash"`
	PrevHash         string    `json:"prev_hash"`
	Sequence         int32     `json:"sequence"`
	TransactionCount int32     `json:"transaction_count"`
	OperationCount   int32     `json:"operation_count"`
	ClosedAt         time.Time `json:"closed_at"`
}

func (l LedgerResource) SseData() interface{} { return l }
func (l LedgerResource) Err() error           { return nil }

//TODO: return the paging token for the ledger, not the id
func (l LedgerResource) SseId() string { return l.Id }

func NewLedgerResource(in db.LedgerRecord) LedgerResource {
	self := fmt.Sprintf("/ledgers/%d", in.Sequence)
	return LedgerResource{
		Links: halgo.Links{}.
			Self(self).
			Link("transactions", self+"/transactions{?after}{?limit}{?order}").
			Link("operations", self+"/operations{?after}{?limit}{?order}").
			Link("effects", self+"/effects{?after}{?limit}{?order}"),
		Id:       in.LedgerHash,
		Hash:     in.LedgerHash,
		PrevHash: in.PreviousLedgerHash,
		Sequence: in.Sequence,
	}
}

func ledgerIndexAction(c web.C, w http.ResponseWriter, r *http.Request) {
	ah := &ActionHelper{c: c, r: r}
	app := ah.App()
	_, order, limit := ah.GetPagingParams()
	after := ah.GetInt32("after")

	if ah.Err() != nil {
		http.Error(w, ah.Err().Error(), http.StatusBadRequest)
		return
	}

	if after == 0 && order == "desc" {
		after = math.MaxInt32
	}

	query := db.LedgerPageQuery{app.HistoryQuery(), after, order, limit}

	render.Collection(w, r, query, ledgerRecordToResource)
}

func ledgerShowAction(c web.C, w http.ResponseWriter, r *http.Request) {
	ah := &ActionHelper{c: c, r: r}
	app := ah.App()
	sequence := ah.GetInt32("id")

	if ah.Err() != nil {
		problem.Render(context.TODO(), w, problem.NotFound)
		return
	}

	query := db.LedgerBySequenceQuery{app.HistoryQuery(), sequence}

	render.Single(w, r, query, ledgerRecordToResource)
}

func ledgerRecordToResource(record db.Record) (render.Resource, error) {
	return NewLedgerResource(record.(db.LedgerRecord)), nil
}
