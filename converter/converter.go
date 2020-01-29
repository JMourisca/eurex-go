// Package converter makes conversions between different currencys based on "Euro foreign exchange reference rates"
// https://www.ecb.europa.eu/stats/policy_and_exchange_rates/euro_reference_exchange_rates/html/index.en.html
package converter

import (
	"encoding/xml"
	"errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type CurrencyRate struct {
	Currency string `xml:"currency,attr"`
	Rate string `xml:"rate,attr"`
}

type RateTime struct {
	Time string `xml:"time,attr"`
	CurrenciesRates []CurrencyRate `xml:"Cube"`
}

type Rates struct {
	Rates []RateTime `xml:"Cube>Cube"`
}

type RatesMap map[string]map[string]float64

// source: ISO Code of the source currency
// target: ISO Code of the target currency
// date: date of rate in ISO format YYYY-MM-DD
// fromFile: whether the data should be retrived from the webshite directly or from the cached file.
func Convert(source string, target string, date string, amount float64) float64 {
	ratesMap := rates()

	if v := ratesMap.validDate(date); v == "" {
		log.Errorf("Date %s not found. Please try another.", date)
		os.Exit(1)
	}
	target, targetVal := ratesMap.rate(date, target)
	source, sourceVal := ratesMap.rate(date, source)
	res := amount * (targetVal / sourceVal)

	return res
}

func (r RatesMap) validDate(date string) string {
	if _, ok := r[date]; ok {
		return date
	}

	return ""
}

// Returns the rate for given currency and date. When currency doesn't exists, defaults to EUR and rate 1.0
func (r RatesMap) rate(date string, currency string) (string, float64) {
	if _, ok := r[date][currency]; ok {
		return currency, r[date][currency]
	}
	log.Infof("The currency %s was not found. Setting EUR as default.", currency)
	return "EUR", 1.0
}

// Gets the data from the file in XML and translate into the structs.
// Returns a Map with the date as main key to make search easier and faster
func rates() RatesMap {
	data, _ := getFromFile()
	rates := Rates{}

	_ = xml.Unmarshal(data, &rates)
	return ratesMapHash(rates)
}

// Converts the data into a Map
func ratesMapHash(rates Rates) RatesMap {
	rm := RatesMap{}

	for _, ratesTimes := range rates.Rates {
		rm[ratesTimes.Time] = map[string]float64{}
		for _, rt := range ratesTimes.CurrenciesRates {
			rm[ratesTimes.Time][rt.Currency], _ = strconv.ParseFloat(rt.Rate, 64)
		}
	}
	return rm
}

// Refetches the data from the source and replace the old file.
func rebuildCache(fileWithPath string) {
	url := "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml"
	//data := make([]byte, 200 * 1024)
	request, err := http.Get(url)
	request.Close = true

	if err != nil || request.StatusCode >= 300 {
		if err == nil {
			err = errors.New(string(request.StatusCode))
		}
	}
	defer request.Body.Close()

	log.Infoln("Fetching from URL")
	body, err := ioutil.ReadAll(request.Body)

	if err != nil {
		log.Errorf("Error reading the data from the source. It will use the old file. Err: %s", err.Error())
		return
	}

	log.Info("Writing file.")
	err = ioutil.WriteFile(fileWithPath, body, 0666)
	if err != nil {
		log.Errorf("Error writing the file. It will use the old one. Err: %s", err.Error())
		return
	}
}

// Gets the rates from exisiting file. If the file is more than 24 hours old, it re-caches from the source.
func getFromFile() ([]byte, error) {
	data := make([]byte, 100 * 1024)
	dir, _ := os.Getwd()
	fileWithPath := filepath.Join(dir, "data", "eurofxref-hist-90d.xml")

	fileInfo, err := os.Stat(fileWithPath)

	now := time.Now()
	if os.IsNotExist(err) || now.Sub(fileInfo.ModTime()).Hours() > 24 {
		log.Infoln("Data is too old or doesn't exist. Will re-create.")
		rebuildCache(fileWithPath)
	}

	file, err := os.Open(fileWithPath)

	if err != nil {
		log.Errorf("Error opening file. Err: %s", err.Error())
		os.Exit(1)
	}
	_, err = file.Read(data)

	if err != nil {
		log.Errorf("Error reading from file. Err: %s", err.Error())
		return data, err
	}

	return data, err
}
