package txsub

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stellar/go-stellar-base/build"
	"github.com/stellar/horizon/test"
)

func TestTxsub(t *testing.T) {
	Convey("txsub.System", t, func() {
		ctx := test.Context()
		submitter := &MockSubmitter{}
		results := &MockResultProvider{}

		system := &System{
			pending:           NewDefaultSubmissionList(),
			submitter:         submitter,
			results:           results,
			networkPassphrase: build.TestNetwork.Passphrase,
		}

		successTx := Result{
			Hash:           "c492d87c4642815dfb3c7dcce01af4effd162b031064098a0d786b6e0a00fd74",
			LedgerSequence: 2,
			EnvelopeXDR:    "AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAACgAAAAAAAAABAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAAAO5rKAAAAAAAAAAABVvwF9wAAAEAKZ7IPj/46PuWU6ZOtyMosctNAkXRNX9WCAI5RnfRk+AyxDLoDZP/9l3NvsxQtWj9juQOuoBlFLnWu8intgxQA",
			ResultXDR:      "xJLYfEZCgV37PH3M4Br07/0WKwMQZAmKDXhrbgoA/XQAAAAAAAAACgAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAA==",
		}

		badSeq := SubmissionResult{
			Err: &FailedTransactionError{"AAAAAAAAAAD////7AAAAAA=="},
		}

		Convey("Submit", func() {
			Convey("returns the result provided by the ResultProvider", func() {
				results.ResultForHash = &successTx
				r := <-system.Submit(ctx, successTx.EnvelopeXDR)

				So(r.Err, ShouldBeNil)
				So(r.Hash, ShouldEqual, successTx.Hash)
				So(submitter.WasSubmittedTo, ShouldBeFalse)
			})

			Convey("returns the error from submission if no result is found by hash and the submitter returns an error", func() {
				submitter.R.Err = errors.New("busted for some reason")
				r := <-system.Submit(ctx, successTx.EnvelopeXDR)

				So(r.Err, ShouldNotBeNil)
				So(submitter.WasSubmittedTo, ShouldBeTrue)
			})

			Convey("if the error is bad_seq and the result at the transaction's sequence number is for the same hash, return result", func() {
				submitter.R = badSeq
				results.ResultForAddressAndSequence = &successTx

				r := <-system.Submit(ctx, successTx.EnvelopeXDR)

				So(r.Err, ShouldBeNil)
				So(r.Hash, ShouldEqual, successTx.Hash)
				So(submitter.WasSubmittedTo, ShouldBeTrue)
			})

			Convey("if the error is bad_seq and the result isn't for the same hash, return error", func() {
				submitter.R = badSeq
				results.ResultForAddressAndSequence = &successTx
				results.ResultForAddressAndSequence.Hash = "some_other_hash"
				r := <-system.Submit(ctx, successTx.EnvelopeXDR)

				So(r.Err, ShouldNotBeNil)
				So(submitter.WasSubmittedTo, ShouldBeTrue)
			})

			Convey("if error is bad_seq and no result is found, return error", func() {
				submitter.R = badSeq
				r := <-system.Submit(ctx, successTx.EnvelopeXDR)

				So(r.Err, ShouldNotBeNil)
				So(submitter.WasSubmittedTo, ShouldBeTrue)
			})

			Convey("if no result found and no error submitting, add to open transaction list", func() {
				_ = system.Submit(ctx, successTx.EnvelopeXDR)
				pending := system.pending.Pending()
				So(len(pending), ShouldEqual, 1)
				t.Logf("passphrase: %s", system.networkPassphrase)
				So(pending[0], ShouldEqual, successTx.Hash)
			})
		})

		Convey("Tick", func() {
			// At each tick
			// TODO:   if no open transactions, don't error out
			// TODO:   checks for results, finishing any available
			// TODO:   times-out and removes old submissions
			// TODO:   if open transactions, but result provider has no new results, keep transactions in open list
		})

	})
}

type MockSubmitter struct {
	R              SubmissionResult
	WasSubmittedTo bool
}

func (sub *MockSubmitter) Submit(env string) SubmissionResult {
	sub.WasSubmittedTo = true
	return sub.R
}

type MockResultProvider struct {
	ResultForHash               *Result
	ResultForAddressAndSequence *Result
}

func (results *MockResultProvider) ResultByHash(hash string) (Result, bool) {
	if results.ResultForHash == nil {
		return Result{}, false
	}

	r := *results.ResultForHash
	return r, true
}

func (results *MockResultProvider) ResultByAddressAndSequence(address string, sequence uint64) (Result, bool) {
	if results.ResultForAddressAndSequence == nil {
		return Result{}, false
	}

	r := *results.ResultForAddressAndSequence
	return r, true
}
