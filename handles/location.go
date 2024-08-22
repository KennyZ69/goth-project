package handles

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func fetchCountries() ([]map[string]string, error) {
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

	var countries []map[string]string
	for _, item := range data {
		if name, ok := item["name"].(map[string]interface{}); ok {
			if commonName, ok := name["common"].(string); ok {
				if cca2, ok := item["cca2"].(string); ok {
					countries = append(countries, map[string]string{"name": commonName, "code": cca2})
				}
			}
		}
	}
	return countries, nil
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
		if query == "" || strings.Contains(strings.ToLower(country["name"]), strings.ToLower(query)) {
			fmt.Fprintf(w, `<option value="%s">%s</option>`, country["name"], country["name"])
		}
	}
	return nil
}
