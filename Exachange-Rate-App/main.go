package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ExchangeRateAPIResponse struct {
	Rates map[string]float64 `json:"rates"`
}

func main() {
	http.HandleFunc("/convert", convertHandler)
	port := ":8080"
	log.Printf("Listening on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func convertHandler(w http.ResponseWriter, r *http.Request) {
	fromCurrency := r.FormValue("from")
	toCurrency := r.FormValue("to")
	amountStr := r.FormValue("amount")
	amount := parseFloat(amountStr)

	if fromCurrency == "" || toCurrency == "" || amount == 0 {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	exchangeRate := fetchExchangeRate(fromCurrency, toCurrency)
	if exchangeRate == 0 {
		http.Error(w, "Exchange rate not found", http.StatusNotFound)
		return
	}

	convertedAmount := amount * exchangeRate

	response := map[string]float64{
		"amount": convertedAmount,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

func fetchExchangeRate(from, to string) float64 {
	apiURL := fmt.Sprintf("https://api.exchangerate-api.com/v4/latest/%s", from)
	resp, err := http.Get(apiURL)
	if err != nil {
		log.Printf("Failed to fetch exchange rate: %v\n", err)
		return 0
	}
	defer resp.Body.Close()

	var data ExchangeRateAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("Failed to decode response: %v\n", err)
		return 0
	}

	rate, ok := data.Rates[to]
	if !ok {
		log.Printf("Exchange rate for %s to %s not found\n", from, to)
		return 0
	}

	return rate
}
