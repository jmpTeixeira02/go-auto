package scrapper

import (
	"github.com/gocolly/colly"
)

const (
	carCard      = ".ooa-yca59n.epwfahw0"
	titleSection = ".epwfahw9.ooa-1ed90th.er34gjf0"
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

// Returns a tuple of strings with the model and link for the car
// Car follows the following format: Brand Model
func GetCarModel(c *Car, e *colly.HTMLElement) {
	e.ForEach(titleSection, func(_ int, el *colly.HTMLElement) {
		c.Model = el.Text
		el.ForEach("a", func(_ int, i_el *colly.HTMLElement) {
			c.Link = i_el.Attr("href")
		})
	})
}

// Returns a string containing the displacement and hp
// "X XXX cm3 • XXX cv"
func GetCarPower(c *Car, e *colly.HTMLElement) {
	e.ForEach(".epwfahw10.ooa-1tku07r.er34gjf0", func(_ int, el *colly.HTMLElement) {
		c.Power = el.Text
	})
}
