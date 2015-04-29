package horizon

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stellar/go-horizon/test"
	"testing"
)

func TestLedgerActions(t *testing.T) {

	Convey("Ledger Actions:", t, func() {
		test.LoadScenario("base")
		app := NewTestApp()
		rh := NewRequestHelper(app)

		Convey("GET /ledgers/1", func() {
			w := rh.Get("/ledgers/1", test.RequestHelperNoop)

			So(w.Code, ShouldEqual, 200)

			var result ledgerResource
			err := json.Unmarshal(w.Body.Bytes(), &result)
			So(err, ShouldBeNil)
			So(result.Sequence, ShouldEqual, 1)
		})

		Convey("GET /ledgers/100", func() {
			w := rh.Get("/ledgers/100", test.RequestHelperNoop)

			So(w.Code, ShouldEqual, 404)
		})

		Convey("GET /ledgers", func() {

			w := rh.Get("/ledgers", test.RequestHelperNoop)

			var result map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &result)
			So(err, ShouldBeNil)
			So(w.Code, ShouldEqual, 200)

			embedded := result["_embedded"].(map[string]interface{})
			records := embedded["records"].([]interface{})

			So(len(records), ShouldEqual, 4)
		})
	})
}
