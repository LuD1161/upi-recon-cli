package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

// Regex matches for user input check
// https://stackoverflow.com/a/38554480
var isMobileNumber = regexp.MustCompile(`^[6-9]\d{9}$`).MatchString //"^[a-zA-Z0-9+_.-]+@[a-zA-Z0-9.-]+$"

func GetUPI(searchTerm string) LambdaResponse {
	threads := 100
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
		goIbiR := GoIbResponse{}
		result := <-resultsChan

		if result.Errors != nil {
			continue
		}

		body, err := ioutil.ReadAll(result.Result.Body)

		if err != nil {
			log.Info().Msgf("Error : %s", err.Error())
			continue
		}

		result.Result.Body.Close()
		err = json.Unmarshal(body, &goIbiR)

		if err != nil {
			log.Info().Msgf("Error : %s", err.Error())
			continue
		}

		if goIbiR.Name != "" {
			// log.Info().Msgf("Response : %+v", cashFR)
			lResponse.Results = append(lResponse.Results, result.VPA)
		}
	}
	return lResponse
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("SecretCode") != "th1s1sS3r3T" {
		http.Error(w, "Incorrect SecretCode Header", http.StatusForbidden)
		return
	}

	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"hello":"world"}`))
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error decoding request body.", http.StatusBadRequest)
			return
		}
		myEvent := MyEvent{}
		err = json.Unmarshal(body, &myEvent)
		if err != nil {
			http.Error(w, "Error decoding request body.", http.StatusBadRequest)
			return
		}
		results := GetUPI(myEvent.SearchTerm)
		// The APIGatewayProxyResponse.Body field needs to be a string, so
		// we marshal the book record into JSON.
		log.Info().Msgf("Results : %+v", results)
		js, err := json.Marshal(results)
		if err != nil {
			http.Error(w, "Error marshalling results", http.StatusBadGateway)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(js)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "I can't do that.")
	}
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	switch os.Getenv("LogLevel") {
	case "Debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.With().Timestamp().Stack().Caller().Logger()
	case "Error":
		log.Logger = log.With().Timestamp().Stack().Caller().Logger()
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	http.HandleFunc("/", mainHandler)
	port := "3133"
	done := make(chan bool)
	go http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), nil)
	log.Info().Msgf("Server listening on port 0.0.0.0:%s", port)
	<-done
}
