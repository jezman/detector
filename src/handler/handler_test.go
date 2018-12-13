package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAlertDetector(t *testing.T) {
	req, err := http.NewRequest("GET", os.Getenv("ALERT_URL"), nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", os.Getenv("ALERT_KEY"))

	// first request
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(alertDetector)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var expected = `{"triggerStatus": "wait confirmation"}`

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// second confirmation request
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(alertDetector)

	handler.ServeHTTP(rr, req)

	expected = `{"triggerStatus": "run"}`

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}
