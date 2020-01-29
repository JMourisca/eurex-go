// Package converter makes conversions between different currencys based on "Euro foreign exchange reference rates"
// https://www.ecb.europa.eu/stats/policy_and_exchange_rates/euro_reference_exchange_rates/html/index.en.html
package converter

import (
	"encoding/xml"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

// Corresponds to the node <Cube currency="USD" rate="1.1025"/>
type CurrencyRate struct {
	Currency string `xml:"currency,attr"`
	Rate string `xml:"rate,attr"`
}

// Correspods to the node <Cube time="2020-01-27"><Cube ... ></Cube></Cube>
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
// It returns source and target again in case the originals don't exist, so the response shows the proper
// default one (EUR)
func Convert(source string, target string, date string, amount float64) (float64, string, string) {
	ratesMap := rates()

	vd, ok := ratesMap.validDate(date)
	if !ok {
		log.Errorf("Date %s not found. Please try another.", date)
		os.Exit(1)
	}
	validTarget, targetVal := ratesMap.rate(vd, target)
	validSource, sourceVal := ratesMap.rate(vd, source)
	res := amount * (targetVal / sourceVal)

	return res, validSource, validTarget
}

// Gets the data from the file in XML and translate into the structs.
// Returns a Map with the date as main key to make search easier and faster
func rates() RatesMap {
	data, _ := getFromFile()
	rates := Rates{}

	_ = xml.Unmarshal(data, &rates)
	return rates.getMap()
}

// Converts the data into a Map so it make search easier and faster
func (r Rates) getMap() RatesMap {
	rm := RatesMap{}

	for _, ratesTimes := range r.Rates {
		rm[ratesTimes.Time] = map[string]float64{}
		for _, rt := range ratesTimes.CurrenciesRates {
			rm[ratesTimes.Time][rt.Currency], _ = strconv.ParseFloat(rt.Rate, 64)
		}
	}
	return rm
}

// Checks if the date exists in the map.
func (rm RatesMap) validDate(date string) (string, bool) {
	if _, ok := rm[date]; ok {
		return date, ok
	}

	// TODO: If the date doesn't exist, fetch the latest one.
	return "", false
}

// Returns the rate for given currency and date. When currency doesn't exists, defaults to EUR and rate 1.0
func (rm RatesMap) rate(date string, currency string) (string, float64) {
	if val, ok := rm[date][currency]; ok {
		return currency, val
	}
	log.Infof("The currency %s was not found. Setting EUR as default.", currency)
	return "EUR", 1.0
}