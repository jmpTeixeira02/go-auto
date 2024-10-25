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
}
