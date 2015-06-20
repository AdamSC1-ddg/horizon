package db

import (
	"database/sql"
	"time"

	"github.com/stellar/go-horizon/log"

	"golang.org/x/net/context"
)

// LedgerClosePump starts a background proc that continually watches the
// history database provided.  The watch is stopped after the provided context
// is cancelled.
//
// Every second, the proc spawned by calling this func will check to see
// if a new ledger has been imported (by ruby-horizon as of 2015-04-30, but
// should eventually end up being in this project).  If a new ledger is seen
// the proc triggers the streaming system to run all watched queries and
// update connected clients
func LedgerClosePump(ctx context.Context, db *sql.DB) {
	go func() {
		var lastSeenLedger int32
		for {
			select {
			case <-time.After(1 * time.Second):
				var latestLedger int32
				row := db.QueryRow("SELECT MAX(sequence) FROM history_ledgers")
				err := row.Scan(&latestLedger)

				if err != nil {
					log.Warn(ctx, "Failed to check latest ledger", err)
					break
				}

				if latestLedger > lastSeenLedger {
					log.Debugf(ctx, "saw new ledger: %d, prev: %d", latestLedger, lastSeenLedger)
					lastSeenLedger = latestLedger
					// TODO: emit, which will trigger the streaming
				}

			case <-ctx.Done():
				log.Info(ctx, "canceling ledger pump")
				return
			}
		}
	}()
}
