package scrapper

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/gocolly/colly"
)

const (
	carCard           = ".ooa-yca59n.epwfahw0"
	titleSection      = ".epwfahw9.ooa-1ed90th.er34gjf0"
	carPowerSection   = ".epwfahw10.ooa-1tku07r.er34gjf0"
	detailsSection    = ".ooa-1uwk9ii.epwfahw11"
	detail            = ".ooa-1omlbtp.epwfahw13"
	detailParameter   = "data-parameter"
	mileageParameter  = "mileage"
	fuelParameter     = "fuel_type"
	yearParameter     = "first_registration_year"
	priceSection      = ".epwfahw16.ooa-1n2paoq.er34gjf0"
	paginationSection = ".pagination-list.ooa-1vdlgt7"
	paginationItem    = ".ooa-g4wbjr.e1y5xfcl0"
)

type Car struct {
	Model   string
	Price   string
	Mileage string
	Fuel    string
	Year    string
	Power   string
	Link    string
}

type pagination struct {
	pages   []page
	current int
}
type page struct {
	Url    string
	Number int
}

type Scrapper struct {
	Collector colly.Collector
	Host      string
}

type carOptions func(*Car, *colly.HTMLElement)

func newPagination() pagination {
	return pagination{
		pages:   []page{},
		current: 1,
	}
}

func New() Scrapper {
	return Scrapper{
		Collector: *colly.NewCollector(),
		Host:      "https://www.standvirtual.com",
	}
}

// Function used to scrape the website
// URL should already contain all the filters and be on page 1
// Pass the desired details to be fetched on carOptions
func (s *Scrapper) Scrape(url string, cars *[]Car, carOptions ...carOptions) error {
	pagination := newPagination()

	s.Collector.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting page: %d\n", pagination.current)
	})

	s.getPaginationInfo(&pagination)
	s.startCars(cars, carOptions...)
	err := s.Collector.Visit(url)
	if err != nil {
		return fmt.Errorf("error visiting url %w", err)
	}

	return s.goToNextPage(&pagination)
}

// Function to start fetching car information
// Opts is used to select the fields to fetch
func (s *Scrapper) startCars(cars *[]Car, opts ...carOptions) {
	// Fetches each card and processes
	s.Collector.OnHTML(carCard, func(e *colly.HTMLElement) {
		var car Car
		for _, opt := range opts {
			opt(&car, e)
		}
		*cars = append(*cars, car)
	})
}

// Function to start fetching pagination information
// Gets the last page number and visits all pages
func (s *Scrapper) getPaginationInfo(pages *pagination) {
	// Fetch Pagination Info and process it
	s.Collector.OnHTML(paginationSection, func(e *colly.HTMLElement) {
		// If there is no pagination info, get it
		if len(pages.pages) == 0 {
			s.updatePaginationInfo(pages, e)
		}
	})
}

// If current page is not next visit the next one
func (s *Scrapper) goToNextPage(p *pagination) error {
	if isFinalPage(p) {
		return nil
	}
	url := p.pages[p.current].Url
	p.current += 1
	err := s.Collector.Visit(url)
	if err != nil {
		return fmt.Errorf("error visiting page! %w", err)
	}
	return s.goToNextPage(p)
}

// Check if current page is the last one
func isFinalPage(p *pagination) bool {
	return p.current == p.pages[len(p.pages)-1].Number
}

func (s *Scrapper) updatePaginationInfo(p *pagination, e *colly.HTMLElement) {
	e.ForEach(paginationItem, func(_ int, e1 *colly.HTMLElement) {
		i, err := strconv.Atoi(e1.Text)
		if err != nil {
			panic(err)
		}
		url, err := url.JoinPath(s.Host, e1.Attr("href"))
		if err != nil {
			panic(err)
		}
		p.pages = append(p.pages, page{
			Url:    url,
			Number: i,
		})
	})
}
