package fortiter

import (
	"github.com/algorand/go-algorand/logging/logspec"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type AgreementEvent struct {
	logspec.AgreementEvent
	Message string    `json:"msg"`
	Time    time.Time `json:"time"`
}

func SaveAgreement(event AgreementEvent, db *sqlx.DB) error {
	a := event.AgreementEvent
	db.MustExec("INSERT INTO agreements (type, round, period, step, hash, sender, object_round, object_period, object_step, weight, weight_total, message, time) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)", a.Type, a.Round, a.Period, a.Step, a.Hash, a.Sender, a.ObjectRound, a.ObjectPeriod, a.ObjectStep, a.Weight, a.WeightTotal, event.Message, event.Time)
	return nil
}
