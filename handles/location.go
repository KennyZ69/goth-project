package handles

import (
	"encoding/json"
	"fmt"
	"github.com/biter777/countries"
	"io"
	"net/http"
	"strings"
)

func fetchCountries() ([]string, error) {
	resp, err := http.Get("https://restcountries.com/v3.1/all")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data []map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	var countries []string
	for _, item := range data {
		if name, ok := item["name"].(map[string]interface{}); ok {
			if commonName, ok := name["common"].(string); ok {
				countries = append(countries, commonName)
			}
		}
	}
	return countries, nil
}

func fetchCities(country string) ([]string, error) {
	apiKey := "kennyz69"
	countryAlphaCode := countries.ByName(country)
	fmt.Printf("there is the alpha code: %v", countryAlphaCode)
	if countryAlphaCode == countries.Unknown {
		return nil, fmt.Errorf("country code in unknown")
	}
	countryCode := countryAlphaCode.Alpha2()
	url := fmt.Sprintf("http://api.geonames.org/searchJSON?country=%s&featureClass=P&maxRows=10&username=%s", countryCode, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data struct {
		Geonames []struct {
			Name string `json:"name"`
		} `json:"geonames"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	var cities []string
	for _, item := range data.Geonames {
		cities = append(cities, item.Name)
	}
	return cities, nil
}

func CityHandler(w http.ResponseWriter, r *http.Request) error {
	countryName := r.URL.Query().Get("country-select")
	fmt.Printf("There is the country code: %v", countryName)
	if countryName == "" {
		return fmt.Errorf("Country not provided")
	}
	fmt.Printf("this is the countryName: %v\n", countryName)
	// cities, err := fetchCities(countryName)
	cities, err := fetchCities(countryName)
	if err != nil {
		http.Error(w, "Unable to fetch cities", http.StatusInternalServerError)
		return err
	}
	// w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "text/html")
	for _, city := range cities {
		// Sending each city as an option
		fmt.Fprintf(w, `<option value="%s">%s</option>`, city, city)
	}
	// json.NewEncoder(w).Encode(cities)
	return nil
}

func CountryHandler(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query().Get("country")
	countries, err := fetchCountries()
	if err != nil {
		http.Error(w, "Unable to fetch countries", http.StatusInternalServerError)
		return nil
	}

	w.Header().Set("Content-Type", "text/html")
	for _, country := range countries {
		if query == "" || strings.Contains(strings.ToLower(country), strings.ToLower(query)) {
			// Filter based on input query
			fmt.Fprintf(w, `<option value="%s">%s</option>`, country, country)
		}
	}
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(countries)
	return nil
}
