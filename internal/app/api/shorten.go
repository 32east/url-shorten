package api

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"net/url"
	"sync"
	"url-short/internal/app/database"
	"url-short/internal/app/models"
)

var alphabet = "QWERTYUIOPASDFGHJKLZXCVBNM0123456789qwertyuiopasdfghjklxcvbnm"
var count = 0
var mut = &sync.Mutex{}
var ctx = context.Background()

const urlLen = 7

func InitializeShorten() {
	var val sql.NullInt64
	var scanErr = database.Postgres.QueryRow(ctx, `select max(id) from urls;`).Scan(&val)
	if scanErr != nil {
		log.Fatal("не удалось выявить максимальный id")
	}

	count = int(val.Int64)
}

func Shorten(w http.ResponseWriter, r *http.Request, response *models.Response, query *models.Response) {
	var qUrl = (*query)["url"]
	if _, err := url.ParseRequestURI(qUrl.(string)); err != nil {
		(*response)["success"], (*response)["reason"], (*response)["code"] = false, "invalid url", http.StatusBadRequest
		return
	}

	var urlStr [urlLen]byte
	var tempCount = count
	for i := 0; i < urlLen; i++ {
		urlStr[urlLen-1-i] = alphabet[tempCount%len(alphabet)]
		tempCount /= len(alphabet)
	}

	var constructed = ""
	for _, val := range urlStr {
		constructed += string(val)
	}

	var _, execErr = database.Postgres.Exec(ctx, `insert into urls(small_url, full_url, ip_creator, created_at)
values ($1, $2, $3, now())`, constructed, qUrl, r.RemoteAddr)

	if execErr != nil {
		(*response)["success"], (*response)["reason"], (*response)["code"] = false, "internal server error", http.StatusInternalServerError
		log.Println("api.Shorten: ", execErr)
		return
	}

	(*response)["success"], (*response)["url"] = true, constructed

	mut.Lock()
	count += 1
	mut.Unlock()
}
