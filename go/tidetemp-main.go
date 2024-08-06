package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"tidetemp-main.go/tidepools"
)

const (
	tideUrlTemplate    = "https://api.tidesandcurrents.noaa.gov/api/prod/datagetter?begin_date=%s&range=24&station=%s&product=predictions&datum=MLLW&time_zone=lst_ldt&interval=hilo&units=english&application=TidePool&format=json"
	weatherUrlTemplate = "https://api.weather.gov/points/%f,%f"
)

type TideResponse struct {
	Predictions []struct {
		T    string `json:"t"`
		V    string `json:"v"`
		Type string `json:"type"`
	} `json:"predictions"`
}

type HourlyWeatherResponse struct {
	Properties struct {
		Periods []struct {
			StartTime       string `json:"startTime"`
			EndTime         string `json:"endTime"`
			Temperature     int    `json:"temperature"`
			TemperatureUnit string `json:"temperatureUnit"`
		} `json:"periods"`
	} `json:"properties"`
}

func fetchTideAPI(stationID string) (string, error) {
	today := time.Now().Format("20060102")
	formattedUrl := fmt.Sprintf(tideUrlTemplate, today, stationID)

	resp, err := http.Get(formattedUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading body: %w", err)
	}

	return string(body), nil
}

func parseTideAPIResponse(tideAPIResponse string) (string, string, error) {
	var response TideResponse
	err := json.Unmarshal([]byte(tideAPIResponse), &response)
	if err != nil {
		return "", "", fmt.Errorf("error parsing tide API response: %w", err)
	}

	if len(response.Predictions) == 0 {
		return "", "", errors.New("no valid tide predictions found")
	}

	lowestTideValue := float64(999)
	var lowestTideTime string

	for _, prediction := range response.Predictions {
		v, err := strconv.ParseFloat(prediction.V, 64)
		if err != nil {
			continue
		}
		if v < lowestTideValue {
			lowestTideValue = v
			lowestTideTime = prediction.T
		}
	}

	if lowestTideTime == "" {
		return "", "", errors.New("no valid tide predictions found")
	}

	parsedTime, err := time.Parse("2006-01-02 15:04", lowestTideTime)
	if err != nil {
		return "", "", fmt.Errorf("error parsing tide time: %w", err)
	}

	lowestTideTimeFormatted := parsedTime.Format("3:04 PM")

	return lowestTideTimeFormatted, fmt.Sprintf("%.3f", lowestTideValue), nil
}

func fetchWeatherAPI(lat, lon float64) (string, error) {
	formattedUrl := fmt.Sprintf(weatherUrlTemplate, lat, lon)
	resp, err := http.Get(formattedUrl)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading body: %w", err)
	}

	return string(body), nil
}

func parseWeatherAPIResponse(weatherAPIResponse string) (string, error) {
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(weatherAPIResponse), &result); err != nil {
		return "", fmt.Errorf("error parsing JSON: %w", err)
	}

	properties, ok := result["properties"].(map[string]interface{})
	if !ok {
		return "", errors.New("error accessing properties in response")
	}

	forecastHourly, ok := properties["forecastHourly"].(string)
	if !ok {
		return "", errors.New("error accessing forecastHourly in properties")
	}

	return forecastHourly, nil
}

func fetchHourlyWeatherAPI(hourlyForecastURL string) (string, error) {
	resp, err := http.Get(hourlyForecastURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading body: %w", err)
	}

	return string(body), nil
}

func parseHourlyWeatherAPIResponse(hourlyWeatherAPIResponse string, lowestTideTimeFormatted string) (string, int, string, error) {
	var response HourlyWeatherResponse
	err := json.Unmarshal([]byte(hourlyWeatherAPIResponse), &response)
	if err != nil {
		return "", 0, "", fmt.Errorf("error parsing hourly weather API response: %w", err)
	}

	lowestTideTime, err := time.Parse("3:04 PM", lowestTideTimeFormatted)
	if err != nil {
		return "", 0, "", fmt.Errorf("error converting tide time for comparison: %w", err)
	}

	for _, period := range response.Properties.Periods {
		startTime, err := time.Parse(time.RFC3339, period.StartTime)
		if err != nil {
			continue
		}

		if startTime.Hour() == lowestTideTime.Hour() {
			return period.StartTime, period.Temperature, period.TemperatureUnit, nil
		}
	}

	return "", 0, "", errors.New("no matching period found for the lowest tide time")
}

func processLocation(location tidepools.TidePool, wg *sync.WaitGroup) {
	defer wg.Done()

	tideResponse, err := fetchTideAPI(location.Station)
	if err != nil {
		fmt.Printf("%s: Error fetching tide API: %v\n", location.Name, err)
		return
	}

	lowestTideTimeFormatted, _, err := parseTideAPIResponse(tideResponse)
	if err != nil {
		fmt.Printf("%s: Error parsing tide data: %v\n", location.Name, err)
		return
	}

	initialWeatherResponse, err := fetchWeatherAPI(location.Lat, location.Long)
	initialHourlyForecastURL, initialURLError := "", ""
	if err != nil {
		initialURLError = fmt.Sprintf("Primary location failed, attempting backup location for weather API: %v", err)
	} else {
		initialHourlyForecastURL, err = parseWeatherAPIResponse(initialWeatherResponse)
		if err != nil {
			initialURLError = fmt.Sprintf("Error parsing weather API response: %v", err)
		}
	}

	backupWeatherResponse, backupErr := fetchWeatherAPI(location.BackupLat, location.BackupLong)
	backupHourlyForecastURL, backupURLError := "", ""
	if backupErr != nil {
		backupURLError = fmt.Sprintf("Error fetching weather API with backup location: %v", backupErr)
	} else {
		backupHourlyForecastURL, err = parseWeatherAPIResponse(backupWeatherResponse)
		if err != nil {
			backupURLError = fmt.Sprintf("Error parsing weather API response with backup location: %v", err)
		}
	}

	hourlyWeatherResponse, err := fetchHourlyWeatherAPI(initialHourlyForecastURL)
	if err != nil && backupHourlyForecastURL != "" {
		hourlyWeatherResponse, err = fetchHourlyWeatherAPI(backupHourlyForecastURL)
	}

	if err != nil {
		fmt.Printf("%s: Error fetching hourly weather API: %v\n", location.Name, err)
		return
	}

	_, _, _, err = parseHourlyWeatherAPIResponse(hourlyWeatherResponse, lowestTideTimeFormatted)
	if err != nil {
		errorMessage := fmt.Sprintf("%s: Error parsing hourly weather data: %v.", location.Name, err)
		if initialURLError != "" {
			errorMessage += fmt.Sprintf(" Initial URL error: %s", initialURLError)
		}
		if backupURLError != "" {
			errorMessage += fmt.Sprintf(" Backup URL error: %s", backupURLError)
		}
		if initialHourlyForecastURL != "" {
			errorMessage += fmt.Sprintf(" Initial URL: %s", initialHourlyForecastURL)
		}
		if backupHourlyForecastURL != "" {
			errorMessage += fmt.Sprintf(" Backup URL: %s", backupHourlyForecastURL)
		}
		fmt.Println(errorMessage)
		return
	}

	fmt.Printf("%s: OK\n", location.Name)
}

func main() {
	var wg sync.WaitGroup

	for _, location := range tidepools.Locations {
		wg.Add(1)
		go processLocation(location, &wg)
	}

	wg.Wait()
}
