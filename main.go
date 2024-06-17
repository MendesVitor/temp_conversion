package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
)

type ViaCEPResponse struct {
	Localidade string `json:"localidade"`
	Erro       bool   `json:"erro"`
}

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
		TempF float64 `json:"temp_f"`
	} `json:"current"`
}

type WeatherResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type HttpClient interface {
	Get(url string) (*http.Response, error)
}

var httpClient HttpClient = &http.Client{}

func main() {
	http.HandleFunc("/clima", handleWeatherRequest)
	port := ":8080"
	fmt.Printf("Server listening on port %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Panic("Server failed to start:", err)
	}
}

func handleWeatherRequest(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	if len(cep) != 8 {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	location, err := getLocationByCEP(cep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	weather, err := getWeatherByLocation(location)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := WeatherResponse{
		TempC: weather.Current.TempC,
		TempF: weather.Current.TempF,
		TempK: celsiusToKelvin(weather.Current.TempC),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getLocationByCEP(cep string) (string, error) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	resp, err := httpClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch data from ViaCEP: %v", err)
	}
	defer resp.Body.Close()

	var viaCEP ViaCEPResponse
	err = json.NewDecoder(resp.Body).Decode(&viaCEP)
	if err != nil {
		return "", fmt.Errorf("failed to decode ViaCEP response: %v", err)
	}

	if viaCEP.Erro {
		return "", fmt.Errorf("can not find zipcode")
	}

	return viaCEP.Localidade, nil
}

func getWeatherByLocation(location string) (*WeatherAPIResponse, error) {
	apiKey := "0d955ca900874ca3a08200551241606"

	sanitizedLocation := url.QueryEscape(location)
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, sanitizedLocation)

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from WeatherAPI: %v", err)
	}
	defer resp.Body.Close()

	var weather WeatherAPIResponse
	err = json.NewDecoder(resp.Body).Decode(&weather)
	if err != nil {
		return nil, fmt.Errorf("failed to decode WeatherAPI response: %v", err)
	}

	return &weather, nil
}

func celsiusToKelvin(celsius float64) float64 {
	return roundToPrecision(celsius+273.15, 1)
}

func roundToPrecision(value float64, precision int) float64 {
	p := math.Pow(10, float64(precision))
	return math.Round(value*p) / p
}
