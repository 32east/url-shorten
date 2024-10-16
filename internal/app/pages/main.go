package pages

import (
	"context"
	"database/sql"
	"errors"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
	"net/http"
	"strings"
	"url-short/internal/app/database"
	"url-short/internal/app/templates"
)

var ctx = context.Background()

func getISOCountry(addr string) (country *geoip2.Country) {
	var db, err = geoip2.Open("GeoLite2-Country.mmdb")
	if err != nil {
		return
	}

	defer db.Close()

	var split = strings.Split(addr, ":")
	var ip = net.ParseIP(split[0])
	var record, recordErr = db.Country(ip)
	if recordErr != nil {
		return
	}

	return record
}

func Main(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		var split = strings.Split(r.URL.Path, "/")

		if len(split) > 0 {
			var smallUrl string
			var fullUrl string
			var scanErr = database.Postgres.QueryRow(ctx, `select small_url, full_url from urls where small_url = $1;`, split[1]).Scan(&smallUrl, &fullUrl)

			if scanErr != nil {
				if !errors.Is(scanErr, sql.ErrNoRows) {
					log.Println("r.URL.Path: ", r.URL.Path, ": ", scanErr)
				}

				goto skip
			}

			http.Redirect(w, r, fullUrl, http.StatusFound)

			var addr = r.RemoteAddr
			var userAgent = r.Header.Get("User-Agent")
			var record = getISOCountry(addr)

			go func() {
				if record != nil {
					var _, execErr = database.Postgres.Exec(ctx, `insert into clicks(small_url, date, ip, country, user_agent)
values ($1, now(), $2, $3, $4)`, smallUrl, addr, record.Country.IsoCode, userAgent)

					if execErr != nil {
						log.Println(execErr)
					}
				} else {
					var _, execErr = database.Postgres.Exec(ctx, `insert into clicks(small_url, date, ip, user_agent)
values ($1, now(), $2, $3)`, smallUrl, addr, userAgent)

					if execErr != nil {
						log.Println(execErr)
					}
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
