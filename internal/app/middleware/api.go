package middleware

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
	"url-short/internal/app/models"
)

var APIAntiSpam = make(map[string]*models.Spammer)

func API(path string, method string, exec func(w http.ResponseWriter, r *http.Request, response *models.Response, query *models.Response)) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		var response = make(models.Response)

		defer func() {
			var code = response["code"]

			if code != nil {
				w.WriteHeader(code.(int))

				delete(response, "code")
			}

			var encoderErr = json.NewEncoder(w).Encode(response)

			if encoderErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal server error"))
			}
		}()

		var addr = r.RemoteAddr
		var split = strings.Split(addr, ":")
		var eAddr = split[0]
		var spammer, ok = APIAntiSpam[eAddr]

		if ok {
			if time.Now().After((*spammer).TimeLimit) {
				(*spammer).Count = 0
			}

			if (*spammer).Count >= 5 {
				response["success"], response["reason"], response["code"] = false, "try again later", http.StatusTooEarly
				return
			}
		} else if !ok {
			APIAntiSpam[eAddr] = &models.Spammer{
				Addr:      eAddr,
				Count:     0,
				TimeLimit: time.Now().Add(time.Second * 30),
			}

			spammer, ok = APIAntiSpam[eAddr]
		}

		(*spammer).Count += 1
		(*spammer).TimeLimit = time.Now().Add(time.Second * 30)

		if r.Method != method {
			response["success"], response["reason"], response["code"] = false, "method not allowed", 403
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			response["success"], response["reason"], response["code"] = false, "content-type not allowed", 403
			return
		}

		var query *models.Response
		if method == "POST" {
			var bRead, bReadErr = io.ReadAll(r.Body)

			if bReadErr != nil {
				response["success"], response["reason"], response["code"] = false, "internal server error", 500
				return
			}

			var sRead = make(models.Response)
			var dErr = json.Unmarshal(bRead, &sRead)

			if dErr != nil {
				response["success"], response["reason"], response["code"] = false, "internal server error", 500
				return
			}

			query = &sRead
		}

		exec(w, r, &response, query)
	})
}
