package scraper

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/king-jam/ft-alert-bot/models"
)

var timeStamp string

func doEvery(d time.Duration, f func() (*models.Forecast, error)) {
	for range time.Tick(d) {
		f()
	}
}

// ScrapeAndParse goes and gets the data every set duration
func ScrapeAndParse(d time.Duration) (*models.SnowPlaces, error) {
	doEvery(d, Scraper)

}

// Scraper scrapes
func Scraper() (*models.Forecast, error) {
	// Request the HTML page.
	resp, err := http.Get("https://www.weather.gov/source/box/winter/snow_prob.xml")
	if err != nil {
		// handle error
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var forecast models.Forecast
	xml.Unmarshal(body, &forecast)
	return &forecast, nil
}

// Parser parses
func Parser(forecast models.Forecast) (*models.SnowPlaces, error) {

	// snowPlaces := make(models.SnowPlaces)
	// snowCity := make(models.SnowCity)
	snowPlaces := make(models.SnowPlaces, 0)
	line := strings.Split(forecast.Text, "\n")
	for _, v := range line {
		// input validation for bad strings from xml parse
		if v == " " || v == "" {
			continue
		}
		lineItems := strings.Split(v, ",")
		lowEndSnowfall, err := strconv.ParseFloat(lineItems[4], 64)
		if err != nil {
			return nil, err
		}
		expectedSnowfall, err := strconv.ParseFloat(lineItems[5], 64)
		if err != nil {
			return nil, err
		}
		highEndSnowfall, err := strconv.ParseFloat(lineItems[6], 64)
		if err != nil {
			return nil, err
		}
		chanceMoreThanZero, err := strconv.ParseFloat(lineItems[7], 64)
		if err != nil {
			return nil, err
		}
		chanceMoreThanOne, err := strconv.ParseFloat(lineItems[8], 64)
		if err != nil {
			return nil, err
		}
		chanceMoreThanTwo, err := strconv.ParseFloat(lineItems[9], 64)
		if err != nil {
			return nil, err
		}
		chanceMoreThanFour, err := strconv.ParseFloat(lineItems[10], 64)
		if err != nil {
			return nil, err
		}
		chanceMoreThanSix, err := strconv.ParseFloat(lineItems[11], 64)
		if err != nil {
			return nil, err
		}
		chanceMoreThanEight, err := strconv.ParseFloat(lineItems[12], 64)
		if err != nil {
			return nil, err
		}
		chanceMoreThanTwelve, err := strconv.ParseFloat(lineItems[13], 64)
		if err != nil {
			return nil, err
		}
		chanceMoreThanEighteen, err := strconv.ParseFloat(lineItems[14], 64)
		if err != nil {
			return nil, err
		}
		snow := models.Snow{
			Place:                  lineItems[0],
			State:                  lineItems[1],
			County:                 lineItems[2],
			LowEndSnowfall:         lowEndSnowfall,
			ExpectedSnowfall:       expectedSnowfall,
			HighEndSnowfall:        highEndSnowfall,
			ChanceMoreThanZero:     chanceMoreThanZero,
			ChanceMoreThanOne:      chanceMoreThanOne,
			ChanceMoreThanTwo:      chanceMoreThanTwo,
			ChanceMoreThanFour:     chanceMoreThanFour,
			ChanceMoreThanSix:      chanceMoreThanSix,
			ChanceMoreThanEight:    chanceMoreThanEight,
			ChanceMoreThanTwelve:   chanceMoreThanTwelve,
			ChanceMoreThanEighteen: chanceMoreThanEighteen,
		}
		snowPlaces = append(snowPlaces, snow)
		// snowCity[lineItems[1]] = snow
		// snowPlaces[lineItems[0]] = snowCity
	}

	for _, s := range snowPlaces {
		if s.Place == "Hopkinton" && s.State == "MA" {
			fmt.Printf("%+v\n", s)
		}
	}

	return &snowPlaces, nil
}
