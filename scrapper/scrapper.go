package scrapper

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

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

type Pagination struct {
	pages   []Page
	current int
}
type Page struct {
	Url    string
	Number int
}

type Scrapper struct {
	Collector colly.Collector
	Host      string
}

type CarOptions func(*Car, *colly.HTMLElement)

func NewPagination() Pagination {
	return Pagination{
		pages:   []Page{},
		current: 1,
	}
}

func New() Scrapper {
	return Scrapper{
		Collector: *colly.NewCollector(),
		Host:      "https://www.standvirtual.com",
	}
}

func (s *Scrapper) Scrape(url string, cars *[]Car, carOptions ...CarOptions) {
	pagination := NewPagination()

	s.Collector.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting page: %d\n", pagination.current)
	})

	s.GetPaginationInfo(&pagination)
	s.StartCars(cars, getCarModel, getCarPrice, getCarPower, getCarDetails)
	err := s.Collector.Visit(url)
	if err != nil {
		panic(err)
	}

	s.GoToNextPage(&pagination)
}

// Function to start fetching car information
// Opts is used to select the fields to fetch
func (s *Scrapper) StartCars(cars *[]Car, opts ...CarOptions) {
	// Fetches each card and processes
	s.Collector.OnHTML(carCard, func(e *colly.HTMLElement) {
		var car Car
		for _, opt := range opts {
			opt(&car, e)
		}
		// Check if it's already in DB
		*cars = append(*cars, car)
		// Notify
	})
}

// Function to start fetching pagination information
// Gets the last page number and visits all pages
func (s *Scrapper) GetPaginationInfo(pages *Pagination) {
	// Fetch Pagination Info and process it
	s.Collector.OnHTML(paginationSection, func(e *colly.HTMLElement) {
		// If there is no pagination info, get it
		if len(pages.pages) == 0 {
			s.getPaginationInfo(pages, e)
		}
	})
}

// If current page is not next visit the next one
func (s *Scrapper) GoToNextPage(p *Pagination) {
	if isFinalPage(p) {
		return
	}
	p.current += 1
	url := p.pages[p.current-1].Url
	err := s.Collector.Visit(url)
	if err != nil {
		panic(fmt.Errorf("error visiting page! %w", err))
	}
	s.GoToNextPage(p)
}

// Check if current page is the last one
func isFinalPage(p *Pagination) bool {
	return p.current == p.pages[len(p.pages)-1].Number
}

func (s *Scrapper) getPaginationInfo(p *Pagination, e *colly.HTMLElement) {
	e.ForEach(paginationItem, func(_ int, e1 *colly.HTMLElement) {
		i, err := strconv.Atoi(e1.Text)
		if err != nil {
			panic(err)
		}
		url, err := url.JoinPath(s.Host, e1.Attr("href"))
		if err != nil {
			panic(err)
		}
		p.pages = append(p.pages, Page{
			Url:    url,
			Number: i,
		})
	})
}

// TitleSection contains a title with the brand and model and href to the link
// Title follows the following format: Brand Model
func getCarModel(c *Car, e *colly.HTMLElement) {
	e.ForEach(titleSection, func(_ int, el *colly.HTMLElement) {
		c.Model = strings.TrimSpace(el.Text)
		el.ForEach("a", func(_ int, i_el *colly.HTMLElement) {
			c.Link = i_el.Attr("href")
		})
	})
}

// Power Section has the displacement followed by the horsepower
// "X XXX cm3 • XXX cv"
func getCarPower(c *Car, e *colly.HTMLElement) {
	e.ForEach(carPowerSection, func(_ int, el *colly.HTMLElement) {
		c.Power = strings.TrimSpace(el.Text)
	})
}

// Details section contains the mileage, fuelType and year
// each of which it's identified by a detailParameter
func getCarDetails(c *Car, e *colly.HTMLElement) {
	// Details Section
	e.ForEach(detailsSection, func(_ int, el *colly.HTMLElement) {
		// Detail
		e.ForEach(detail, func(_ int, i_el *colly.HTMLElement) {
			dataParameter := i_el.Attr(detailParameter)
			switch dataParameter {
			case mileageParameter:
				c.Mileage = strings.TrimSpace(i_el.Text)
			case fuelParameter:
				c.Fuel = strings.TrimSpace(i_el.Text)
			case yearParameter:
				c.Year = strings.TrimSpace(i_el.Text)
			}
		})
	})
}

func getCarPrice(c *Car, e *colly.HTMLElement) {
	// Price
	e.ForEach(priceSection, func(_ int, el *colly.HTMLElement) {
		c.Price = strings.TrimSpace(el.Text)
	})
}
