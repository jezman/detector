package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	confirmed           bool
	waitingConfirmation bool
	file                *os.File
	err                 error
)

func logPrint(r *http.Request) {
	log.Println(r.RemoteAddr)
	log.Println(r.URL)
	log.Println(r.UserAgent())
}

func timer() {
	time.Sleep(3 * time.Second)

	if waitingConfirmation && !confirmed {
		log.Println("alert not confirmation, alert counter has reset")
		waitingConfirmation = false
	}
}

func reset() {
	time.Sleep(3 * time.Second)

	confirmed = false
	waitingConfirmation = false
}


func alertsHandler(w http.ResponseWriter, r *http.Request) {
	if len(r.Header["Authorization"]) == 0 || r.Header["Authorization"][0] != os.Getenv("ALERT_KEY") {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("wrong authorization key")
		logPrint(r)
	} else {
		logPrint(r)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		detectors, ok := r.URL.Query()["detector"]

		if !ok || len(detectors[0]) < 1 {
			log.Println("Url Param 'detector' is missing")
		} else {
			log.Println("detector - " + string(detectors[0]))
		}

		if !confirmed && !waitingConfirmation {
			waitingConfirmation = true

			go timer()

			io.WriteString(w, `{"triggerStatus": "waiting for confirmation"}`)
			log.Println("alert detected, waiting for confirmation")
		} else if !confirmed && waitingConfirmation {
			log.Println("alert confirmed")
			io.WriteString(w, `{"triggerStatus": "run"}`)

			confirmed = true
			go reset()
		}
	}
}

func main() {
	confirmed = false
	waitingConfirmation = false

	file, err = os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	log.SetOutput(file)

	log.Println("handler starting...")
	http.HandleFunc(os.Getenv("ALERT_URL"), alertsHandler)
	http.ListenAndServe(":8000", nil)
}
