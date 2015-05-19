package db

import (
	"strconv"
	"testing"

	_ "github.com/lib/pq"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stellar/go-horizon/test"
)

func TestLedgerPageQuery(t *testing.T) {
	test.LoadScenario("base")
	ctx := test.Context()
	db := OpenTestDatabase()
	defer db.Close()

	Convey("LedgerPageQuery", t, func() {
		pq, err := NewPageQuery("0", "asc", 3)
		So(err, ShouldBeNil)

		q := LedgerPageQuery{SqlQuery{db}, pq}
		ledgers, err := Results(ctx, q)

		So(err, ShouldBeNil)
		So(len(ledgers), ShouldEqual, 3)

		// ensure each record is after the previous
		current := q.Cursor

		for _, ledger := range ledgers {
			ledger := ledger.(LedgerRecord)
			So(ledger.Id, ShouldBeGreaterThan, current)
			current = ledger.Id
		}

		lastLedger := ledgers[len(ledgers)-1].(Pageable)
		cursor, _ := strconv.ParseInt(lastLedger.PagingToken(), 10, 64)
		q.Cursor = cursor

		ledgers, err = Results(ctx, q)

		So(err, ShouldBeNil)
		So(len(ledgers), ShouldEqual, 1)

		current = q.Cursor

		for _, ledger := range ledgers {
			ledger := ledger.(LedgerRecord)
			So(ledger.Id, ShouldBeGreaterThan, current)
			current = ledger.Id
		}

	})
}
