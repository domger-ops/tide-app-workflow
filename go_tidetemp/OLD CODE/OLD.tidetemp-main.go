package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	"tidetemp.go/tidepools"
)

type Period struct {
	StartTime   time.Time `json:"startTime"`
	Temperature int       `json:"temperature"`
}

type Properties struct {
	Periods []Period `json:"periods"`
}

type ForecastHourly struct {
	Properties Properties `json:"properties"`
}

type Prediction struct {
	T    time.Time `json:"t"`
	V    string    `json:"v"`
	Type string    `json:"type"`
}

type PredictionsData struct {
	Predictions []Prediction `json:"predictions"`
}

func fetchForecastHourly(forecastURL string) (*ForecastHourly, error) {
	resp, err := http.Get(forecastURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var forecastHourly ForecastHourly
	if err := json.Unmarshal(body, &forecastHourly); err != nil {
		return nil, err
	}

	return &forecastHourly, nil
}

func fetchForecastHourlyURL(lat, lon float64) (string, error) {
	url := fmt.Sprintf("https://api.weather.gov/points/%f,%f", lat, lon)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	var result map[string]interface{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	propertiesMap, ok := result["properties"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("failed to extract 'properties' as map[string]interface{}")
	}

	forecastHourlyURL, ok := propertiesMap["forecastHourly"].(string)
	if !ok {
		return "", fmt.Errorf("'forecastHourly' URL not found in properties")
	}
	return forecastHourlyURL, nil
}

func findLowestTideAndFetchTemperature(data []string, location tidepools.TidePool) {
	var lowestTide Prediction
	var lowestValue float64 = -1
	for _, jsonData := range data {
		var predictionsData PredictionsData
		if err := json.Unmarshal([]byte(jsonData), &predictionsData); err != nil {
			fmt.Println("Error parsing JSON:", err)
			continue
		}

		for _, prediction := range predictionsData.Predictions {
			if prediction.Type == "L" {
				value, err := strconv.ParseFloat(prediction.V, 64)
				if err != nil {
					fmt.Println("Error parsing value:", err)
					continue
				}
				if lowestValue == -1 || value < lowestValue {
					lowestValue = value
					lowestTide = prediction
				}
			}
		}
	}

	if lowestValue == -1 {
		fmt.Println("No low tide found.")
		return
	}

	forecastURL, err := fetchForecastHourlyURL(location.Lat, location.Long)
	if err != nil {
		fmt.Println("Error fetching forecast URL:", err)
		return
	}

	forecastHourly, err := fetchForecastHourly(forecastURL)
	if err != nil {
		fmt.Println("Error fetching forecast data:", err)
		return
	}

	lowestTideTime := lowestTide.T.Format(time.RFC3339)
	for _, period := range forecastHourly.Properties.Periods {
		if period.StartTime.Format(time.RFC3339) == lowestTideTime {
			fmt.Printf("Lowest tide at %s, Temperature: %d°F\n", lowestTide.T.Format("2006-01-02T15:04:05"), period.Temperature)
			return
		}
	}
	fmt.Println("Temperature for the lowest tide time not found.")
}

func fetchTideData(station string, today string, wg *sync.WaitGroup, dataChan chan<- string) {
	defer wg.Done()

	noaaTideApiUrl := fmt.Sprintf("https://api.tidesandcurrents.noaa.gov/api/prod/datagetter?begin_date=%s&range=48&station=%s&product=predictions&datum=MLLW&time_zone=lst_ldt&interval=hilo&units=english&application=TidePool&format=json", today, station)
	resp, err := http.Get(noaaTideApiUrl)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var predictionsData PredictionsData
	if err := json.Unmarshal(body, &predictionsData); err != nil {
		fmt.Println("Error marshalling predictions data:", err)
		return
	}

	data, err := json.Marshal(predictionsData)
	if err != nil {
		fmt.Println("Error marshalling predictions data:", err)
		return
	}

	dataChan <- string(data)
}

func main() {
	today := time.Now().Format("20060102")
	stations := tidepools.GetAllStationNumbers()
	var wg sync.WaitGroup
	dataChan := make(chan string, len(stations))

	for _, station := range stations {
		wg.Add(1)
		go fetchTideData(station, today, &wg, dataChan)
	}

	go func() {
		wg.Wait()
		close(dataChan)
	}()

	var data []string
	for d := range dataChan {
		data = append(data, d)
	}

	for _, location := range tidepools.Locations {
		findLowestTideAndFetchTemperature(data, location)
	}
}
