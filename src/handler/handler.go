package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var alertCount uint8

func logPrint(r *http.Request) {
	log.Println(r.RemoteAddr)
	log.Println(r.URL)
	log.Println(r.UserAgent())
}

func timer() {
	time.Sleep(3 * time.Minute)

	if alertCount > 0 {
		log.Println("alert not confirmation, alert counter has reset")
		alertCount = 0
	}
}

func alertDetector(w http.ResponseWriter, r *http.Request) {
	if len(r.Header["Authorization"]) == 0 || r.Header["Authorization"][0] != os.Getenv("ALERT_KEY") {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("wrong authorization key")
		logPrint(r)
	} else {
		alertCount++
		logPrint(r)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if alertCount == 1 {
			go timer()
			io.WriteString(w, `{"triggerStatus": "wait confirmation"}`)
		} else if alertCount > 1 {
			io.WriteString(w, `{"triggerStatus": "run"}`)
			alertCount = 0
		}
	}
}

func main() {
	http.HandleFunc(os.Getenv("ALERT_URL"), alertDetector)
	http.ListenAndServe(":8000", nil)
}
