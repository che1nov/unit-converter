package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// ConversionData holds data for displaying the conversion result
type ConversionData struct {
	InputValue     float64 // The value entered by the user
	FromUnit       string  // Unit to convert from
	ToUnit         string  // Unit to convert to
	ConvertedValue float64 // The converted value
}

// Constants for conversion factors
var (
	lengthConversions = map[string]float64{
		"millimeter": 0.001,
		"centimeter": 0.01,
		"meter":      1,
		"kilometer":  1000,
		"inch":       0.0254,
		"foot":       0.3048,
		"yard":       0.9144,
		"mile":       1609.34,
	}

	weightConversions = map[string]float64{
		"milligram": 0.000001,
		"gram":      0.001,
		"kilogram":  1,
		"ounce":     0.0283495,
		"pound":     0.453592,
	}

	tempConversions = map[string]func(float64) float64{
		"celsius":    func(c float64) float64 { return c },
		"fahrenheit": func(f float64) float64 { return (f - 32) * 5 / 9 },
		"kelvin":     func(k float64) float64 { return k - 273.15 },
	}
)

// Handler for the home page
func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

// Handler for performing conversions
func convertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Retrieve form data
	inputValueStr := r.FormValue("value")
	fromUnit := r.FormValue("from_unit")
	toUnit := r.FormValue("to_unit")
	conversionType := r.FormValue("conversion_type")

	// Convert input value to a number
	inputValue, err := strconv.ParseFloat(inputValueStr, 64)
	if err != nil || inputValue < 0 {
		http.Error(w, "Invalid input value", http.StatusBadRequest)
		return
	}

	var convertedValue float64

	// Perform conversion based on the selected type
	switch conversionType {
	case "length":
		convertedValue = convertLength(inputValue, fromUnit, toUnit)
	case "weight":
		convertedValue = convertWeight(inputValue, fromUnit, toUnit)
	case "temperature":
		convertedValue = convertTemperature(inputValue, fromUnit, toUnit)
	default:
		http.Error(w, "Invalid conversion type", http.StatusBadRequest)
		return
	}

	// Prepare data for the result page
	data := ConversionData{
		InputValue:     inputValue,
		FromUnit:       fromUnit,
		ToUnit:         toUnit,
		ConvertedValue: convertedValue,
	}

	// Render the result page
	tmpl := template.Must(template.ParseFiles("templates/result.html"))
	tmpl.Execute(w, data)
}

// Converts length units
func convertLength(value float64, fromUnit, toUnit string) float64 {
	baseValue := value * lengthConversions[fromUnit]
	return baseValue / lengthConversions[toUnit]
}

// Converts weight units
func convertWeight(value float64, fromUnit, toUnit string) float64 {
	baseValue := value * weightConversions[fromUnit]
	return baseValue / weightConversions[toUnit]
}

// Converts temperature units
func convertTemperature(value float64, fromUnit, toUnit string) float64 {
	// Convert to base unit (Celsius)
	baseValue := tempConversions[fromUnit](value)

	// Convert from base unit to target unit
	switch toUnit {
	case "celsius":
		return baseValue
	case "fahrenheit":
		return baseValue*9/5 + 32
	case "kelvin":
		return baseValue + 273.15
	default:
		return 0
	}
}

func main() {
	// Define routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/convert", convertHandler)

	// Start the server
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
