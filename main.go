package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const siteVerifyURL = "https://hcaptcha.com/siteverify"

var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

// Regex matches for user input check
// https://stackoverflow.com/a/38554480
var isGmail = regexp.MustCompile("[a-zA-Z0-9+_.-]+@gmail.com").MatchString   //"^[a-zA-Z0-9+_.-]+@[a-zA-Z0-9.-]+$"
var isMobileNumber = regexp.MustCompile("^[6-9]\\d{9}$").MatchString         //"^[a-zA-Z0-9+_.-]+@[a-zA-Z0-9.-]+$"
var isCarReg = regexp.MustCompile("[a-zA-Z]{2}[a-zA-Z0-9]{8-9}").MatchString //"^[a-zA-Z0-9+_.-]+@[a-zA-Z0-9.-]+$"

type LoginRequest struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	RecaptchaResponse string `json:"h-captcha-response"`
}

type SiteVerifyResponse struct {
	Success     bool      `json:"success"`
	Score       float64   `json:"score"`
	Action      string    `json:"action"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

type MyEvent struct {
	Type             string `json:"type"`
	SearchTerm       string `json:"searchTerm"`
	SiteKey          string `json:"siteKey"`
	HCaptchaResponse string `json:"hcaptcha-response"`
}

func run_scanner(typeOfSearch, searchTerm string) ([]byte, error) {
	var result []byte
	var err error
	app := "./upi-recon-cli"

	arg0 := ""

	log.Printf("Running scanner : %s %s", app, arg0)
	// Check user input for any malicious input and regex match
	switch typeOfSearch {
	case "Gmail":
		if !isGmail(searchTerm) {
			return result, errGmail
		}
		arg0 = "checkGpay"
	case "CarReg":
		arg0 = "checkFastag"
		if !isCarReg(searchTerm) {
			return result, errCarReg
		}
		searchTerm = fmt.Sprintf("netc.%s", searchTerm)
	default:
		if !isMobileNumber(searchTerm) {
			return result, errMobileNumber
		}
		arg0 = ""
	}

	log.Printf("Running scanner : %s %s", app, arg0)
	// https://stackoverflow.com/a/7786922
	cmd := exec.Command(app, arg0, searchTerm)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return []byte{}, err
	}
	return stdout, err
}

func CheckRecaptcha(secret, response string) error {
	req, err := http.NewRequest(http.MethodPost, siteVerifyURL, nil)
	if err != nil {
		return err
	}

	// Add necessary request parameters.
	q := req.URL.Query()
	q.Add("secret", secret)
	q.Add("response", response)
	req.URL.RawQuery = q.Encode()

	// Make request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response.
	var body SiteVerifyResponse
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return err
	}

	// Check recaptcha verification success.
	if !body.Success {
		return errors.New("unsuccessful recaptcha verify request")
	}

	// Check response score.
	// if body.Score < 0.5 {
	// 	return errors.New("lower received score than expected")
	// }

	// Check response action.
	// if body.Action != "login" {
	// 	return errors.New("mismatched recaptcha action")
	// }

	log.Printf("body.Score : %f | body.Action : %s", body.Score, body.Action)
	return nil
}

func HandleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	secret := os.Getenv("HCAPTCHA_SECRET")
	var data MyEvent
	err := json.Unmarshal([]byte(req.Body), &data)
	if err := CheckRecaptcha(secret, data.HCaptchaResponse); err != nil {
		log.Printf("req : %+v", req)
		// return serverError(err)
	}

	results, err := run_scanner(data.Type, data.SearchTerm)
	if err != nil {
		return serverError(err)
	}

	// The APIGatewayProxyResponse.Body field needs to be a string, so
	// we marshal the book record into JSON.
	js, err := json.Marshal(results)
	if err != nil {
		return serverError(err)
	}

	// Return a response with a 200 OK status and the JSON book record
	// as the body.
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(js),
	}, nil
}

// Add a helper for handling errors. This logs any error to os.Stderr
// and returns a 500 Internal Server Error response that the AWS API
// Gateway understands.
func serverError(err error) (events.APIGatewayProxyResponse, error) {
	errorLogger.Println(err.Error())

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

// Similarly add a helper for send responses relating to client errors.
func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
