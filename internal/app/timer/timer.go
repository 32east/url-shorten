package timer

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"url-short/internal/app/middleware"
)

func maxmindDownload() {
	const errorFunction = "maxmindDownload"
	var res *http.Response
	var resErr error

	for {
		res, resErr = http.Get("https://github.com/P3TERX/GeoLite.mmdb/raw/download/GeoLite2-Country.mmdb")

		if resErr != nil {
			res.Body.Close()
			log.Println(errorFunction, resErr)
			time.Sleep(time.Second * 3)
			continue
		}

		var r, rErr = io.ReadAll(res.Body)

		if rErr != nil {
			res.Body.Close()
			log.Println(errorFunction, rErr)
			continue
		}

		os.Remove("GeoLite2-Country.mmdb")
		var f, fErr = os.Create("GeoLite2-Country.mmdb")

		if fErr != nil {
			res.Body.Close()
			log.Println(errorFunction, fErr)
			continue
		}

		f.Write(r)
		f.Close()
		res.Body.Close()

		break
	}
}

func Initialize() {
	var newTicker = time.NewTicker(time.Minute * 5)

	go func() {
		for {
			<-newTicker.C

			for k, v := range middleware.APIAntiSpam {
				if !time.Now().After((*v).TimeLimit) {
					continue
				}

				delete(middleware.APIAntiSpam, k)
			}
		}
	}()

	maxmindDownload()

	var tickerDay = time.NewTicker(time.Hour * 24)

	go func() {
		for {
			<-tickerDay.C
			maxmindDownload()
		}
	}()
}
