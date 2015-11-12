// Package resource contains the type definitions for all of horizons
// response resources.
package resource

import (
	"github.com/stellar/horizon/render/hal"
)

// AccountResource is the summary of an account
type Account struct {
	Links struct {
		Self         hal.Link `json:"self"`
		Transactions hal.Link `json:"transactions"`
		Operations   hal.Link `json:"operations"`
		Payments     hal.Link `json:"payments"`
		Effects      hal.Link `json:"effects"`
		Offers       hal.Link `json:"offers"`
	} `json:"_links"`

	HistoryAccount
	Sequence             int64             `json:"sequence"`
	SubentryCount        int32             `json:"subentry_count"`
	InflationDestination string            `json:"inflation_destination,omitempty"`
	HomeDomain           string            `json:"home_domain,omitempty"`
	Thresholds           AccountThresholds `json:"thresholds"`
	Flags                AccountFlags      `json:"flags"`
	Balances             []Balance         `json:"balances"`
	Signers              []Signer          `json:"signers"`
}

// AccountFlags represents the state of an account's flags
type AccountFlags struct {
	AuthRequired  bool `json:"auth_required"`
	AuthRevocable bool `json:"auth_revocable"`
}

// AccountThresholds represents an accounts "thresholds", the numerical values
// needed to satisfy the authorization of a given operation.
type AccountThresholds struct {
	LowThreshold  byte `json:"low_threshold"`
	MedThreshold  byte `json:"med_threshold"`
	HighThreshold byte `json:"high_threshold"`
}

// Balance represents an account's holdings for a single currency type
type Balance struct {
	Type    string `json:"asset_type"`
	Balance string `json:"balance"`
	// additional trustline data
	Code   string `json:"asset_code,omitempty"`
	Issuer string `json:"issuer,omitempty"`
	Limit  string `json:"limit,omitempty"`
}

// HistoryAccount is a simple resource, used for the account collection
// actions.  It provides only the TotalOrderId of the account and its address.
type HistoryAccount struct {
	ID      string `json:"id"`
	PT      string `json:"paging_token"`
	Address string `json:"address"`
}

// Signer represents one of an account's signers.
type Signer struct {
	Address string `json:"address"`
	Weight  int32  `json:"weight"`
}
