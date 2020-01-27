package converter

import (
	"math"
	"os"
	"testing"
)

func compare(x, y float64) bool {
	tolerance := 0.00001
	diff := math.Abs(x - y)
	mean := math.Abs(x + y)/2.0
	return (diff/mean) < tolerance
}

func TestGetFromFile(t *testing.T) {
	_ = os.Chdir("..")
	_, err := getFromFile()
	_ = os.Chdir("converter")
	if err != nil {
		t.Error("Can't open file. Err: ", err)
	}
}

func TestGetRates(t *testing.T) {
	_ = os.Chdir("..")
	if data := rates(); len(data) == 0 {
		t.Error("No data found.")
	}
	_ = os.Chdir("converter")
}

func TestRate(t *testing.T) {
	_ = os.Chdir("..")
	data := rates()

	if invalidDate := data.validDate("Some-invalid-data"); invalidDate != "" {
		t.Error("Expected empty invalid date, got ", invalidDate)
	}

	_ = os.Chdir("converter")
}

func TestConvert(t *testing.T) {
	_ = os.Chdir("..")
	validSource := "BRL"
	validTarget := "CHF"
	validDate := "2020-01-24"
	expectedResult := 0.232445
	currentResult := Convert(validSource, validTarget, validDate)

	if v := compare(currentResult, expectedResult); !v {
		t.Errorf("Expected %f, got %f, %t", expectedResult, currentResult, v)
	}

	invalidSource := "BRL1"
	validTarget = "CHF"
	validDate = "2020-01-24"
	expectedResult = 1.071200
	currentResult = Convert(invalidSource, validTarget, validDate)

	if v := compare(currentResult, expectedResult); !v {
		t.Errorf("Expected %f, got %f, %t", expectedResult, currentResult, v)
	}

	_ = os.Chdir("converter")
}
