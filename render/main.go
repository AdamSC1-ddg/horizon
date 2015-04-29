package render

import (
	"github.com/stellar/go-horizon/db"
	"github.com/stellar/go-horizon/render/hal"
	"net/http"
)

type Transform func(interface{}) interface{}

func Collection(w http.ResponseWriter, r *http.Request, q db.Query, t Transform) {
	// TODO: negotiate, see if we should stream

	records, err := db.Results(q)
	if err != nil {
		panic(err)
	}

	resources := make([]interface{}, len(records))
	for i, record := range records {
		resources[i] = t(record)
	}

	page := hal.Page{
		Records: resources,
	}

	hal.RenderPage(w, page)
}

func Single(w http.ResponseWriter, r *http.Request, q db.Query, t Transform) {
	// TODO: negotiate, see if we should stream

	record, err := db.First(q)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else if record == nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		resource := t(record)
		hal.Render(w, resource)
	}
}
