package scrapper

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTestData() []byte {
	path, err := filepath.Abs("../testdata/scrapper.html")
	if err != nil {
		panic(fmt.Errorf("could not open file: %w", err))
	}
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("could read  file: %w", err))
	}
	return bytes
}

func localHtmlSetup() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write(getTestData())
	})
	return httptest.NewServer(handler)
}

func TestGetCarModel(t *testing.T) {
	// Arrange
	server := localHtmlSetup()
	s := New()

	//Act
	var cars []Car
	s.StartCars(&cars, GetCarModel)

	//Assert
	err := s.Visit(server.URL)
	if err != nil {
		panic(err)
	}

	assert.NotNil(t, cars)
	assert.EqualValues(t, "Mercedes-Benz A 200 d AMG Line", cars[0].Model)
}

func TestGetCarPower(t *testing.T) {
	// Arrange
	server := localHtmlSetup()
	s := New()

	//Act
	var cars []Car
	s.StartCars(&cars, GetCarPower)

	//Assert
	err := s.Visit(server.URL)
	if err != nil {
		panic(err)
	}

	assert.NotNil(t, cars)
	assert.EqualValues(t, "2 143 cm3 • 136 cv", cars[0].Power)
}

func TestGetCarDetails(t *testing.T) {
	// Arrange
	server := localHtmlSetup()
	s := New()

	//Act
	var cars []Car
	s.StartCars(&cars, GetCarDetails)

	//Assert
	err := s.Visit(server.URL)
	if err != nil {
		panic(err)
	}

	assert.NotNil(t, cars)
	assert.EqualValues(t, "92 580 km", cars[0].Mileage)
	assert.EqualValues(t, "Diesel", cars[0].Fuel)
	assert.EqualValues(t, "2017", cars[0].Year)
}

func TestGetCarPrice(t *testing.T) {
	// Arrange
	server := localHtmlSetup()
	s := New()

	//Act
	var cars []Car
	s.StartCars(&cars, GetCarPrice)

	//Assert
	err := s.Visit(server.URL)
	if err != nil {
		panic(err)
	}

	assert.NotNil(t, cars)
	assert.EqualValues(t, "23 990", cars[0].Price)
}
