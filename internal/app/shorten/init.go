package shorten

import (
	"context"
	"database/sql"
	"log"
	"url-short/internal/app/database"
)

var Count = 0
var ctx = context.Background()

func Initialize() {
	var val sql.NullInt64
	var scanErr = database.Postgres.QueryRow(ctx, `select max(id) from urls;`).Scan(&val)
	if scanErr != nil {
		log.Fatal("не удалось выявить максимальный id")
	}

	Count = int(val.Int64)
}
