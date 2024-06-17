package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var viaCEPMockResponse = `{"localidade": "São Paulo", "erro": false}`
var weatherAPIMockResponse = `{
	"location": {
		"name": "São Paulo",
		"region": "",
		"country": "Brazil",
		"lat": -23.55,
		"lon": -46.63,
		"tz_id": "America/Sao_Paulo",
		"localtime_epoch": 1622168162,
		"localtime": "2024-06-16 17:49"
	},
	"current": {
		"temp_c": 28.5,
		"temp_f": 83.3,
		"condition": {
			"text": "Partly cloudy",
			"icon": "//cdn.weatherapi.com/weather/64x64/day/116.png",
			"code": 1003
		}
	}
}`

type mockClient struct{}

func (m *mockClient) Get(url string) (*http.Response, error) {
	if strings.Contains(url, "viacep") {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(viaCEPMockResponse)),
		}, nil
	} else if strings.Contains(url, "weatherapi") {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(weatherAPIMockResponse)),
		}, nil
	}
	return nil, fmt.Errorf("unexpected URL: %s", url)
}

func TestGetLocationByCEP(t *testing.T) {
	httpClient = &mockClient{}

	location, err := getLocationByCEP("01001000")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expectedLocation := "São Paulo"
	if location != expectedLocation {
		t.Errorf("expected %v, got %v", expectedLocation, location)
	}
}

func TestGetWeatherByLocation(t *testing.T) {
	httpClient = &mockClient{}
	os.Setenv("WEATHERAPI_KEY", "test_key")

	weather, err := getWeatherByLocation("São Paulo")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expectedTempC := 28.5
	if weather.Current.TempC != expectedTempC {
		t.Errorf("expected %v, got %v", expectedTempC, weather.Current.TempC)
	}
}

func TestCelsiusToKelvin(t *testing.T) {
	celsius := 28.5
	expected := 301.7
	result := celsiusToKelvin(celsius)

	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestHandleWeatherRequest(t *testing.T) {
	httpClient = &mockClient{}

	t.Run("valid CEP", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/clima?cep=01001000", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handleWeatherRequest)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := WeatherResponse{
			TempC: 28.5,
			TempF: 83.3,
			TempK: 301.7,
		}

		var actual WeatherResponse

		if err := json.NewDecoder(rr.Body).Decode(&actual); err != nil {
			t.Errorf("could not decode response: %v", err)
		}
		if actual != expected {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	})

	t.Run("invalid CEP", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/clima?cep=123", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handleWeatherRequest)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusUnprocessableEntity {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnprocessableEntity)
		}

		expected := "invalid zipcode\n"
		if rr.Body.String() != expected {
			t.Errorf("expected %v, got %v", expected, rr.Body.String())
		}
	})
}
