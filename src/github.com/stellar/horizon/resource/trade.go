package resource

import (
	"errors"

	"github.com/stellar/horizon/db"
	"github.com/stellar/horizon/render/hal"
)

// Populate fills out the details
func (res *Trade) Populate(row db.EffectRecord) (err error) {
	if row.Type != db.EffectTrade {
		err = errors.New("invalid effect; not a trade")
		return
	}
	row.UnmarshalDetails(res)
	res.ID = row.PagingToken()
	res.PT = row.PagingToken()
	res.Buyer = row.Account

	lb := hal.LinkBuilder{}
	res.Links.Self = lb.Link("/accounts", res.Seller)
	res.Links.Seller = lb.Link("/accounts", res.Seller)
	res.Links.Buyer = lb.Link("/accounts", res.Buyer)
	return
}

// PagingToken implementation for hal.Pageable
func (res Trade) PagingToken() string {
	return res.PT
}
