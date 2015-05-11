package db

import (
	_ "github.com/lib/pq"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stellar/go-horizon/test"
	"testing"
)

func TestTransactionPageQuery(t *testing.T) {
	Convey("TransactionPageQuery", t, func() {
		test.LoadScenario("base")
		db := OpenTestDatabase()
		defer db.Close()

		makeQuery := func(c string, o string, l int32) TransactionPageQuery {
			pq, err := NewPageQuery(c, o, l)

			So(err, ShouldBeNil)

			return TransactionPageQuery{
				GormQuery: GormQuery{&db},
				PageQuery: pq,
			}
		}

		Convey("orders properly", func() {
			// asc orders ascending by id
			records := MustResults(makeQuery("", "asc", 0))

			So(records, ShouldBeOrderedAscending, func(r interface{}) int64 {
				So(r, ShouldHaveSameTypeAs, TransactionRecord{})
				return r.(TransactionRecord).Id
			})

			// desc orders descending by id
			records = MustResults(makeQuery("", "desc", 0))

			So(records, ShouldBeOrderedDescending, func(r interface{}) int64 {
				So(r, ShouldHaveSameTypeAs, TransactionRecord{})
				return r.(TransactionRecord).Id
			})
		})

		Convey("limits properly", func() {
			// returns number specified
			records := MustResults(makeQuery("", "asc", 3))
			So(len(records), ShouldEqual, 3)

			// returns all rows if limit is higher
			records = MustResults(makeQuery("", "asc", 10))
			So(len(records), ShouldEqual, 4)
		})

		Convey("cursor works properly", func() {
			// lowest id if ordered ascending and no cursor
			record := MustFirst(makeQuery("", "asc", 0))
			So(record.(TransactionRecord).Id, ShouldEqual, 12884905984)

			// highest id if ordered descending and no cursor
			record = MustFirst(makeQuery("", "desc", 0))
			So(record.(TransactionRecord).Id, ShouldEqual, 17179873280)

			// starts after the cursor if ordered ascending
			record = MustFirst(makeQuery("12884905984", "asc", 0))
			So(record.(TransactionRecord).Id, ShouldEqual, 12884910080)

			// starts before the cursor if ordered descending
			record = MustFirst(makeQuery("17179873280", "desc", 0))
			So(record.(TransactionRecord).Id, ShouldEqual, 12884914176)
		})

		Convey("restricts to address properly", func() {
			address := "gspbxqXqEUZkiCCEFFCN9Vu4FLucdjLLdLcsV6E82Qc1T7ehsTC"
			q := makeQuery("", "asc", 0)
			q.AccountAddress = address
			records := MustResults(q)

			So(len(records), ShouldEqual, 3)

			for _, r := range records {
				So(r.(TransactionRecord).Account, ShouldEqual, address)
			}
		})

		Convey("restricts to ledger properly", func() {
			q := makeQuery("", "asc", 0)
			q.LedgerSequence = 4
			records := MustResults(q)

			So(len(records), ShouldEqual, 1)

			for _, r := range records {
				So(r.(TransactionRecord).LedgerSequence, ShouldEqual, q.LedgerSequence)
			}
		})
	})
}
