package converter

import (
	"encoding/xml"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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

func Convert(source string, target string, date string) float64 {
	ratesMap := rates()

	if v := ratesMap.validDate(date); v == "" {
		log.Errorf("Date %s not found. Please try something else.", date)
		os.Exit(1)
	}
	target, targetVal := ratesMap.rate(date, target)
	source, sourceVal := ratesMap.rate(date, source)
	res := targetVal / sourceVal

	fmt.Printf("%s = %f %s in %s\n", source, res, target, date)
	return res
}

func (r RatesMap) validDate(date string) string {
	if _, ok := r[date]; ok {
		return date
	}

	return ""
}

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

func (r RatesMap) rate(date string, currency string) (string, float64) {
	if _, ok := r[date][currency]; ok {
		return currency, r[date][currency]
	}
	log.Infof("The currency %s was not found. Setting EUR as default.", currency)
	return "EUR", 1.0
}

func rates() RatesMap {
	//data, err := getFromUrl()
	//
	//if err != nil {
	//	log.Warning("Error fetching from URL. Will try to fetch from file. Err: ", err)
		data, _ := getFromFile()
	//}

	rates := Rates{}

	_ = xml.Unmarshal(data, &rates)
	return ratesMapHash(rates)
}

func getFromUrl() ([]byte, error) {
	url := "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml"
	data := make([]byte, 100 * 1024)
	urlData, err := http.Get(url)
	if err != nil || urlData.StatusCode >= 300 {
		if err == nil {
			err = errors.New("error")
		}
		return data, err
	}
	defer urlData.Body.Close()
	fmt.Println("Fetching from URL")
	_, err = urlData.Body.Read(data)
	return data, err
}

func getFromFile() ([]byte, error) {
	data := make([]byte, 100 * 1024)
	dir, _ := os.Getwd()
	fileWithPath := filepath.Join(dir, "data", "eurofxref-hist-90d.xml")
	file, err := os.Open(fileWithPath)

	if err != nil {
		log.Error("Error opening file. ", err)
		os.Exit(1)
	}
	_, err = file.Read(data)

	if err != nil {
		log.Error("Error reading from file", err)
		return data, err
	}

	return data, err
}
