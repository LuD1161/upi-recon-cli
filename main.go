package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
)

var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

// Regex matches for user input check
// https://stackoverflow.com/a/38554480
var isMobileNumber = regexp.MustCompile(`^[6-9]\d{9}$`).MatchString //"^[a-zA-Z0-9+_.-]+@[a-zA-Z0-9.-]+$"

func GetUPI(searchTerm string) LambdaResponse {
	threads := 8
	lResponse := LambdaResponse{}
	Results := []string{}
	ErrArr := []error{}
	lResponse.Errors = ErrArr
	lResponse.Results = Results

	if !isMobileNumber(searchTerm) {
		lResponse.Errors = append(ErrArr, fmt.Errorf("not a mobile number : %s", searchTerm))
		return lResponse
	}

	vpasChan := make(chan CashF, threads)
	resultsChan := make(chan HTTPResponse)

	for i := 0; i < threads; i++ {
		go MakeRequest(vpasChan, resultsChan)
	}

	go func() {
		for _, upiHandle := range UPIHandles {
			orderHash := fmt.Sprintf("QVsj6YQPg3xq1sLm%s", RandStringBytes(4))
			data := CashF{
				VPA:       fmt.Sprintf("%s@%s", searchTerm, upiHandle),
				OrderHash: orderHash,
			}
			vpasChan <- data
		}
	}()

	for i := 0; i < len(UPIHandles); i++ {
		cashFR := CashFResponse{}
		result := <-resultsChan

		if result.Errors != nil {
			continue
		}

		body, err := ioutil.ReadAll(result.Result.Body)

		if err != nil {
			log.Printf("Error : %s", err.Error())
			continue
		}

		result.Result.Body.Close()
		err = json.Unmarshal(body, &cashFR)

		if err != nil {
			log.Printf("Error : %s", err.Error())
			continue
		}

		if cashFR.Message.VpaStatus == "AVAILABLE" {
			// log.Printf("Response : %+v", cashFR)
			lResponse.Results = append(lResponse.Results, result.VPA)
		}
	}
	return lResponse
}

func HandleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// secret := os.Getenv("HCAPTCHA_SECRET")
	var data MyEvent
	err := json.Unmarshal([]byte(req.Body), &data)
	if err != nil {
		return serverError(err)
	}

	results := GetUPI(data.SearchTerm)

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
		Body:       fmt.Sprintf("%s : %s", http.StatusText(http.StatusInternalServerError), err.Error()),
	}, nil
}

func main() {
	results := GetUPI("9882539413")
	js, _ := json.Marshal(results)
	fmt.Printf("%s", js)
	// lambda.Start(HandleRequest)
}
