package api

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"sync"
	"url-short/internal/app/database"
	"url-short/internal/app/models"
	"url-short/internal/app/shorten"
)

var alphabet = "QWERTYUIOPASDFGHJKLZXCVBNM0123456789qwertyuiopasdfghjklxcvbnm"
var mut = &sync.Mutex{}
var ctx = context.Background()

// TODO: Позже сделать нормальный блэклист.
var Blacklist = map[string]bool{
	"localhost:8080": true,
}

const urlLen = 7

func Shorten(w http.ResponseWriter, r *http.Request, response *models.Response, query *models.Response) {
	var qUrl = (*query)["url"].(string)
	var urlParsed, err = url.ParseRequestURI(qUrl)
	if err != nil {
		(*response)["success"], (*response)["reason"], (*response)["code"] = false, "invalid url", http.StatusBadRequest
		return
	}

	var urlStr [urlLen]byte
	var tempCount = shorten.Count
	for i := 0; i < urlLen; i++ {
		urlStr[urlLen-1-i] = alphabet[tempCount%len(alphabet)]
		tempCount /= len(alphabet)
	}

	var constructed = ""
	for _, val := range urlStr {
		constructed += string(val)
	}

	var _, ok = Blacklist[urlParsed.Host]

	if ok {
		(*response)["success"], (*response)["reason"], (*response)["code"] = false, "blacklisted url", http.StatusBadRequest
		return
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
	shorten.Count += 1
	mut.Unlock()
}
