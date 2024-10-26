package scrapper

import (
	"strings"

	"github.com/gocolly/colly"
)

const (
	carCard          = ".ooa-yca59n.epwfahw0"
	titleSection     = ".epwfahw9.ooa-1ed90th.er34gjf0"
	carPowerSection  = ".epwfahw10.ooa-1tku07r.er34gjf0"
	detailsSection   = ".ooa-1uwk9ii.epwfahw11"
	detail           = ".ooa-1omlbtp.epwfahw13"
	detailParameter  = "data-parameter"
	mileageParameter = "mileage"
	fuelParameter    = "fuel_type"
	yearParameter    = "first_registration_year"
	priceSection     = ".epwfahw16.ooa-1n2paoq.er34gjf0"
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
