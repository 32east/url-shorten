package timer

import (
	"time"
	"url-short/internal/app/middleware"
)

func Initialize() {
	var newTicker = time.NewTicker(time.Minute * 5)

	for {
		<-newTicker.C

		for k, v := range middleware.APIAntiSpam {
			if !time.Now().After((*v).TimeLimit) {
				continue
			}

			delete(middleware.APIAntiSpam, k)
		}
	}
}
