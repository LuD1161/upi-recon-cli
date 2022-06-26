package main

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func MakeRequest(payloadsChan <-chan CashF, resultsChan chan<- HTTPResponse) {
	client := http.Client{Timeout: time.Duration(3) * time.Second}
	url := "https://payments.cashfree.com/pgbillpayuiapi/upi/validate"

	for payload := range payloadsChan {

		response := HTTPResponse{
			Result: &http.Response{},
			Errors: nil,
			VPA:    payload.VPA,
		}

		// log.Printf("Payload : %+v", payload)
		jsonPayload, _ := json.Marshal(payload)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
		if err != nil {
			response.Errors = err
			resultsChan <- response
			continue
		}

		req.Header.Add("Connection", "close")
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		response.Result = resp
		response.Errors = err
		resultsChan <- response

	}
}
