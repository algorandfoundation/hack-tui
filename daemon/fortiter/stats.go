package fortiter

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

var StatsSchema = `
CREATE TABLE IF NOT EXISTS stats (
    address TEXT PRIMARY KEY,
    sent INTEGER,
    received INTEGER,
    failed INTEGER,
    success INTEGER
)
`

type Stats struct {
	Address  string
	Sent     uint64
	Received uint64
	Failed   uint64
	Success  uint64
}

func (s *Stats) String() string {
	return fmt.Sprintf("Address: %s\nSent: %d\nReceived: %d\nFailed: %d\nSuccess: %d\n",
		s.Address, s.Sent, s.Received, s.Failed, s.Success)
}

func (s *Stats) SaveStats(db sqlx.DB) error {
	var stats Stats
	err := db.Get(&stats, "SELECT * FROM stats WHERE address = ?", s.Address)
	if err != nil {
		db.MustExec("INSERT INTO stats (address, sent, received, failed, success) VALUES (?, ?, ?, ?, ?)", s.Address, s.Sent, s.Received, s.Failed, s.Success)
	} else {
		db.MustExec("UPDATE stats SET sent = ?, received = ?, failed = ?, success = ? WHERE address = ?",
			s.Sent, s.Received, s.Failed, s.Success, s.Address)
	}
	return nil
}
