package cmd

type VPAResponse struct {
	VPA          string `json:"vpa"`
	Status       string `json:"status"`
	CustomerName string `json:"customer_name"`
	Error        error  `json:"error"`
}
