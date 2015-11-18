package effects

import (
	"github.com/stellar/horizon/db"
	"github.com/stellar/horizon/render/hal"
	"github.com/stellar/horizon/resource/base"
	"golang.org/x/net/context"
)

var TypeNames = map[int32]string{
	db.EffectAccountCreated:           "account_created",
	db.EffectAccountRemoved:           "account_removed",
	db.EffectAccountCredited:          "account_credited",
	db.EffectAccountDebited:           "account_debited",
	db.EffectAccountThresholdsUpdated: "account_thresholds_updated",
	db.EffectAccountHomeDomainUpdated: "account_home_domain_updated",
	db.EffectAccountFlagsUpdated:      "account_flags_updated",
	db.EffectSignerCreated:            "signer_created",
	db.EffectSignerRemoved:            "signer_removed",
	db.EffectSignerUpdated:            "signer_updated",
	db.EffectTrustlineCreated:         "trustline_created",
	db.EffectTrustlineRemoved:         "trustline_removed",
	db.EffectTrustlineUpdated:         "trustline_updated",
	db.EffectTrustlineAuthorized:      "trustline_authorized",
	db.EffectTrustlineDeauthorized:    "trustline_deauthorized",
	db.EffectOfferCreated:             "offer_created",
	db.EffectOfferRemoved:             "offer_removed",
	db.EffectOfferUpdated:             "offer_updated",
	db.EffectTrade:                    "trade",
}

func New(
	ctx context.Context,
	row db.EffectRecord,
) (result hal.Pageable, err error) {

	switch row.Type {
	case db.EffectAccountCreated:
		e := AccountCreated{}
		e.Populate(ctx, row)
		err = row.UnmarshalDetails(&e)
		result = e
	case db.EffectAccountCredited:
		e := AccountCredited{}
		e.Populate(ctx, row)
		err = row.UnmarshalDetails(&e)
		result = e
	case db.EffectAccountDebited:
		e := AccountDebited{}
		e.Populate(ctx, row)
		err = row.UnmarshalDetails(&e)
		result = e
	case db.EffectAccountThresholdsUpdated:
		e := AccountThresholdsUpdated{}
		e.Populate(ctx, row)
		err = row.UnmarshalDetails(&e)
		result = e
	case db.EffectAccountHomeDomainUpdated:
		e := AccountHomeDomainUpdated{}
		e.Populate(ctx, row)
		err = row.UnmarshalDetails(&e)
		result = e
	case db.EffectAccountFlagsUpdated:
		e := AccountFlagsUpdated{}
		e.Populate(ctx, row)
		err = row.UnmarshalDetails(&e)
		result = e
	case db.EffectSignerCreated:
		e := SignerCreated{}
		e.Populate(ctx, row)
		err = row.UnmarshalDetails(&e)
		result = e
	case db.EffectSignerUpdated:
		e := SignerUpdated{}
		e.Populate(ctx, row)
		err = row.UnmarshalDetails(&e)
		result = e
	case db.EffectSignerRemoved:
		e := SignerRemoved{}
		e.Populate(ctx, row)
		err = row.UnmarshalDetails(&e)
		result = e
	case db.EffectTrustlineCreated:
		e := TrustlineCreated{}
		e.Populate(ctx, row)
		err = row.UnmarshalDetails(&e)
		result = e
	case db.EffectTrustlineUpdated:
		e := TrustlineUpdated{}
		e.Populate(ctx, row)
		err = row.UnmarshalDetails(&e)
		result = e
	case db.EffectTrustlineRemoved:
		e := TrustlineRemoved{}
		e.Populate(ctx, row)
		err = row.UnmarshalDetails(&e)
		result = e
	case db.EffectTrustlineAuthorized:
		e := TrustlineAuthorized{}
		e.Populate(ctx, row)
		err = row.UnmarshalDetails(&e)
		result = e
	case db.EffectTrustlineDeauthorized:
		e := TrustlineDeauthorized{}
		e.Populate(ctx, row)
		err = row.UnmarshalDetails(&e)
		result = e
	case db.EffectTrade:
		e := Trade{}
		e.Populate(ctx, row)
		err = row.UnmarshalDetails(&e)
		result = e
	default:
		e := Base{}
		e.Populate(ctx, row)
		result = e
	}

	return
}

type Base struct {
	Links struct {
		Operation hal.Link `json:"operation"`
		Succeeds  hal.Link `json:"succeeds"`
		Precedes  hal.Link `json:"precedes"`
	} `json:"_links"`

	ID      string `json:"id"`
	PT      string `json:"paging_token"`
	Account string `json:"account"`
	Type    string `json:"type"`
	TypeI   int32  `json:"type_i"`
}

type AccountCreated struct {
	Base
	StartingBalance string `json:"starting_balance"`
}

type AccountCredited struct {
	Base
	base.Asset
	Amount string `json:"amount"`
}

type AccountDebited struct {
	Base
	base.Asset
	Amount string `json:"amount"`
}

type AccountThresholdsUpdated struct {
	Base
	LowThreshold  int32 `json:"low_threshold"`
	MedThreshold  int32 `json:"med_threshold"`
	HighThreshold int32 `json:"high_threshold"`
}

type AccountHomeDomainUpdated struct {
	Base
	HomeDomain string `json:"home_domain"`
}

type AccountFlagsUpdated struct {
	Base
	AuthRequired  *bool `json:"auth_required_flag,omitempty"`
	AuthRevokable *bool `json:"auth_revokable_flag,omitempty"`
}

type SignerCreated struct {
	Base
	Weight    int32  `json:"weight"`
	PublicKey string `json:"public_key"`
}

type SignerRemoved struct {
	Base
	Weight    int32  `json:"weight"`
	PublicKey string `json:"public_key"`
}

type SignerUpdated struct {
	Base
	Weight    int32  `json:"weight"`
	PublicKey string `json:"public_key"`
}

type TrustlineCreated struct {
	Base
	base.Asset
	Limit string `json:"limit"`
}

type TrustlineRemoved struct {
	Base
	base.Asset
	Limit string `json:"limit"`
}

type TrustlineUpdated struct {
	Base
	base.Asset
	Limit string `json:"limit"`
}

type TrustlineAuthorized struct {
	Base
	Trustor   string `json:"trustor"`
	AssetType string `json:"asset_type"`
	AssetCode string `json:"asset_code,omitempty"`
}

type TrustlineDeauthorized struct {
	Base
	Trustor   string `json:"trustor"`
	AssetType string `json:"asset_type"`
	AssetCode string `json:"asset_code,omitempty"`
}

type Trade struct {
	Base
	Seller            string `json:"seller"`
	OfferID           int64  `json:"offer_id"`
	SoldAmount        string `json:"sold_amount"`
	SoldAssetType     string `json:"sold_asset_type"`
	SoldAssetCode     string `json:"sold_asset_code,omitempty"`
	SoldAssetIssuer   string `json:"sold_asset_issuer,omitempty"`
	BoughtAmount      string `json:"bought_amount"`
	BoughtAssetType   string `json:"bought_asset_type"`
	BoughtAssetCode   string `json:"bought_asset_code,omitempty"`
	BoughtAssetIssuer string `json:"bought_asset_issuer,omitempty"`
}
