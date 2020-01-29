package converter

import (
	"fmt"
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

	if invalidDate, ok := data.validDate("Some-invalid-data"); ok {
		t.Errorf("Expected empty invalid date, got '%s'", invalidDate)
	}

	_ = os.Chdir("converter")
}

func TestConvert(t *testing.T) {
	_ = os.Chdir("..")

	defaultCurrency := "EUR"

	amount := 10.0
	validSource := "BRL"
	validTarget := "CHF"
	validDate := "2020-01-24"
	expectedResult := 2.324451

	currentResult, newSource, newTarget := Convert(validSource, validTarget, validDate, amount)

	if v := compare(currentResult, expectedResult); !v {
		t.Errorf("Expected %f, got %f.", expectedResult, currentResult)
	}

	if newSource != validSource {
		t.Errorf("Expected %s, got %s.", validSource, newSource)
	}

	if newTarget != validTarget {
		t.Errorf("Expected %s, got %s.", validTarget, newTarget)
	}

	invalidSource := "BRL1"
	validTarget = "CHF"
	validDate = "2020-01-24"
	expectedResult = 10.712
	currentResult, newSource, newTarget = Convert(invalidSource, validTarget, validDate, amount)

	if v := compare(currentResult, expectedResult); !v {
		t.Errorf("Expected %f, got %f, %t", expectedResult, currentResult, v)
	}

	if newSource != defaultCurrency {
		t.Errorf("Expected %s, got %s.", defaultCurrency, newSource)
	}

	_ = os.Chdir("converter")
}

func ExampleConvert() {
	_ = os.Chdir("..")
	fmt.Println(Convert("BRL", "CHF", "2020-01-24", 10.0))
	// Output: 2.3244510025171428 BRL CHF
	_ = os.Chdir("converter")
}
