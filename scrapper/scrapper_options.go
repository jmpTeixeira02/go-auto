package scrapper

import (
	"strings"

	"github.com/gocolly/colly"
)

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
