package cmd

type VPAResponse struct {
	VPA          string `json:"vpa"`
	Success      bool   `json:"success"`
	CustomerName string `json:"customer_name"`
	Error        error  `json:"error"`
}
