package main

import "net/http"

type LoginRequest struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	RecaptchaResponse string `json:"h-captcha-response"`
}

type MyEvent struct {
	SearchTerm string `json:"search_term"`
	SecretCode string `json:"secret_code"`
}

type CashF struct {
	VPA       string `json:"vpa"`
	OrderHash string `json:"orderHash"`
}

type CashFResponse struct {
	Status  string `json:"status"`
	Message struct {
		Status       int    `json:"status"`
		Message      string `json:"message"`
		CustomerName string `json:"customerName"`
		VpaStatus    string `json:"vpaStatus"`
	} `json:"message"`
}

type GoIbResponse struct {
	Error     interface{} `json:"error"`
	ErrorCode string      `json:"error_code"`
	Msg       string      `json:"msg"`
	Name      string      `json:"name"`
	Status    bool        `json:"status"`
}

type LambdaResponse struct {
	Results []string `json:"results"`
	Errors  []error  `json:"errors"`
}

type HTTPResponse struct {
	VPA    string
	Result *http.Response
	Errors error
}
