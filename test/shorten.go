package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var urls = make(map[string]bool)

func main() {
	var cli = &http.Client{
		Timeout: time.Second * 3,
	}

	var m, mErr = json.Marshal(map[string]interface{}{
		"url": "https://github.com/tttttt30/web-forum/blob/master/www/staticfiles/styles/styles.css",
	})

	if mErr != nil {
		log.Fatal("не удалось задекодировать json запрос", mErr)
	}

	for i := 0; i < 999999; i++ {
		var req *http.Request
		var reqErr = fmt.Errorf("-1")
		var recCount = 0

		for reqErr != nil {
			if recCount > 5 {
				log.Fatal("recursion limit exceeded")
			}

			req, reqErr = http.NewRequest("POST", "http://localhost:8080/api/v1/shorten", bytes.NewReader(m))

			if reqErr != nil {
				log.Fatal(reqErr)
			}

			req.Header.Set("Content-Type", "application/json")

			recCount += 1
		}

		var res, resErr = cli.Do(req)
		defer res.Body.Close()

		if resErr != nil {
			log.Fatal(resErr)
		}

		if res.StatusCode != 200 {
			log.Fatal("status code != 200", res)
		}

		var r, rErr = io.ReadAll(res.Body)

		if rErr != nil {
			log.Fatal(rErr)
		}

		fmt.Println("check", i, string(r))
		if _, ok := urls[string(r)]; ok {
			log.Fatal("УЖЕ ЕСТЬ", i, string(r))
		}

		urls[string(r)] = true
	}
}
