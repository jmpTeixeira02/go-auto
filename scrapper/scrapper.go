package scrapper

import (
	"fmt"
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
	*colly.Collector
}

func New() Scrapper {
	return Scrapper{
		Collector: colly.NewCollector(),
	}
}

func (s *Scrapper) StartCars(cars *[]Car, opts ...func(*Car, *colly.HTMLElement)) {
	// Fetches each card and processes
	s.OnHTML(carCard, func(e *colly.HTMLElement) {
		var car Car
		for _, opt := range opts {
			opt(&car, e)
		}
		// Check if it's already in DB
		*cars = append(*cars, car)
		// Notify
	})
}

func (s *Scrapper) StartPagination(pages *Pagination) {
	// Fetch Pagination Info and process it
	s.OnHTML(paginationSection, func(e *colly.HTMLElement) {
		// If there is no pagination info, get it
		if len(pages.pages) == 0 {
			GetPaginationInfo(pages, e)
		}
		// Wait until pagination info is filled
		s.OnScraped(func(_ *colly.Response) {
			GoToNextPage(pages, e)
		})
	})
}

func GoToNextPage(p *Pagination, e *colly.HTMLElement) {
	if isFinalPage(p) {
		return
	}
	p.current += 1
	err := e.Request.Visit(p.pages[p.current-1].Url)
	if err != nil {
		panic(fmt.Errorf("error visiting page! %w", err))
	}
}

func isFinalPage(p *Pagination) bool {
	return p.current != p.pages[len(p.pages)-1].Number
}

func GetPaginationInfo(p *Pagination, e *colly.HTMLElement) {
	e.ForEach(paginationItem, func(_ int, e1 *colly.HTMLElement) {
		i, err := strconv.Atoi(e1.Text)
		if err != nil {
			panic(err)
		}
		p.pages = append(p.pages, Page{Url: e1.Attr("href"), Number: i})
	})
}

// TitleSection contains a title with the brand and model and href to the link
// Title follows the following format: Brand Model
func GetCarModel(c *Car, e *colly.HTMLElement) {
	e.ForEach(titleSection, func(_ int, el *colly.HTMLElement) {
		c.Model = strings.TrimSpace(el.Text)
		el.ForEach("a", func(_ int, i_el *colly.HTMLElement) {
			c.Link = i_el.Attr("href")
		})
	})
}

// Power Section has the displacement followed by the horsepower
// "X XXX cm3 • XXX cv"
func GetCarPower(c *Car, e *colly.HTMLElement) {
	e.ForEach(carPowerSection, func(_ int, el *colly.HTMLElement) {
		c.Power = strings.TrimSpace(el.Text)
	})
}

// Details section contains the mileage, fuelType and year
// each of which it's identified by a detailParameter
func GetCarDetails(c *Car, e *colly.HTMLElement) {
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

func GetCarPrice(c *Car, e *colly.HTMLElement) {
	// Price
	e.ForEach(priceSection, func(_ int, el *colly.HTMLElement) {
		c.Price = strings.TrimSpace(el.Text)
	})
}
