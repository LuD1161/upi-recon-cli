package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const URL = "https://www.goibibo.com/v2payments/upi/vpa/validate/?vpa="

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func MakeRequest(payloadsChan <-chan CashF, resultsChan chan<- HTTPResponse) {
	client := http.Client{Timeout: time.Duration(3) * time.Second}

	for payload := range payloadsChan {
		url := fmt.Sprintf("%s%s", URL, payload.VPA)
		response := HTTPResponse{
			Result: &http.Response{},
			Errors: nil,
			VPA:    payload.VPA,
		}

		req, err := http.NewRequest("GET", url, nil)
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
