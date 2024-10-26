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
	expected := Car{
		Model: "Mercedes-Benz A 200 d AMG Line",
		Link:  "https://www.standvirtual.com/carros/anuncio/mercedes-benz-a-200-d-amg-line-ID8PKMqL.html",
	}

	//Act
	var cars []Car
	s.StartCars(&cars, GetCarModel)

	//Assert
	err := s.Visit(server.URL)
	if err != nil {
		panic(err)
	}

	assert.NotNil(t, cars)
	assert.EqualValues(t, expected, cars[0])
}

func TestGetCarPower(t *testing.T) {
	// Arrange
	server := localHtmlSetup()
	s := New()
	expected := Car{
		Power: "2 143 cm3 • 136 cv",
	}

	//Act
	var cars []Car
	s.StartCars(&cars, GetCarPower)

	//Assert
	err := s.Visit(server.URL)
	if err != nil {
		panic(err)
	}

	assert.NotNil(t, cars)
	assert.EqualValues(t, expected, cars[0])
}

func TestGetCarDetails(t *testing.T) {
	// Arrange
	server := localHtmlSetup()
	s := New()
	expected := Car{
		Mileage: "92 580 km",
		Fuel:    "Diesel",
		Year:    "2017",
	}

	//Act
	var cars []Car
	s.StartCars(&cars, GetCarDetails)

	//Assert
	err := s.Visit(server.URL)
	if err != nil {
		panic(err)
	}

	assert.NotNil(t, cars)
	assert.EqualValues(t, expected, cars[0])
}

func TestGetCarPrice(t *testing.T) {
	// Arrange
	server := localHtmlSetup()
	s := New()
	expected := Car{
		Price: "23 990",
	}

	//Act
	var cars []Car
	s.StartCars(&cars, GetCarPrice)

	//Assert
	err := s.Visit(server.URL)
	if err != nil {
		panic(err)
	}

	assert.NotNil(t, cars)
	assert.EqualValues(t, expected, cars[0])
}

func TestPagination(t *testing.T) {
	// Arrange
	server := localHtmlSetup()
	s := New()
	expected := 1325

	//Act
	var pagination Pagination
	s.StartPagination(&pagination)

	//Assert
	err := s.Visit(server.URL)
	if err != nil {
		panic(err)
	}
	actual := pagination.pages[len(pagination.pages)-1].Number

	assert.NotNil(t, pagination)
	assert.EqualValues(t, expected, actual)
}
