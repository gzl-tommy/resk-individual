package logrus

import (
	log "github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

type DbxLogrusLogger struct {
	log *log.Logger
}

func NewUpperLogrusLogger() *DbxLogrusLogger {
	l := &DbxLogrusLogger{
		log: log.New(),
	}
	std := log.StandardLogger()
	l.log.Level = std.Level
	l.log.Hooks = std.Hooks
	l.log.Formatter = std.Formatter
	l.log.Out = std.Out
	return l
}

func (l *DbxLogrusLogger) Log(q *dbx.QueryStatus) {
	if q.Err == nil {
		l.log.Debug("\n" + q.String())

	} else {
		l.log.Error("\n" + q.String())
	}

	//l.log.WithFields(log.Fields{
	//	"SessID":       q.SessID,
	//	"TxID":         q.TxID,
	//	"LastInsertID": q.LastInsertID,
	//	"Start":        q.Start,
	//	"End":          q.End,
	//	"Err":          q.Err,
	//	"RowsAffected": q.RowsAffected,
	//	"SQL":          q.Query,
	//	"Args":         q.Args,
	//}).Info()
}
