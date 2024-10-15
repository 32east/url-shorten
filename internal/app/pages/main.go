package pages

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"
	"url-short/internal/app/database"
	"url-short/internal/app/templates"
)

var ctx = context.Background()

func Main(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		var split = strings.Split(r.URL.Path, "/")

		if len(split) > 0 {
			var id int
			var fullUrl string
			var scanErr = database.Postgres.QueryRow(ctx, `select id, full_url from urls where small_url = $1;`, split[1]).Scan(&id, &fullUrl)

			if scanErr != nil {
				if !errors.Is(scanErr, sql.ErrNoRows) {
					log.Println("r.URL.Path: ", r.URL.Path, ": ", scanErr)
				}

				goto skip
			}

			http.Redirect(w, r, fullUrl, http.StatusFound)

			go func() {
				var _, execErr = database.Postgres.Exec(ctx, `update urls set clicks = clicks + 1 where id = $1;`, id)

				if execErr != nil {
					log.Println(execErr)
				}
			}()

			return
		}
	skip:
		http.NotFound(w, r)
		return
	}

	templates.Main.Execute(w, nil)
}
