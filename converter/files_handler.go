package converter

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Refetches the data from the source and replace the old file.
func rebuildCache(fileWithPath string) {
	url := "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml"

	response, err := http.Get(url)

	if response == nil {
		log.Errorf("Response is nil. It will use the old file.")
		return
	}

	if err != nil || response.StatusCode >= 300 {
		if err == nil {
			err = errors.New("Response Code: " + string(response.StatusCode))
		}
		log.Errorf("Error reading the data from the source. It will use the old file. Err: %s", err.Error())
		return
	}
	response.Close = true
	defer response.Body.Close()

	log.Infoln("Fetching from URL")
	body, err := ioutil.ReadAll(response.Body)

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
		log.Fatalf("Error opening file. Err: %s", err.Error())
	}

	_, err = file.Read(data)

	if err != nil {
		log.Fatalf("Error reading from file. Err: %s", err.Error())
	}

	return data, err
}
