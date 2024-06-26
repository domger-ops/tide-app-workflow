#!/bin/bash

# Function to fetch and display weather details for a given location
fetch_and_display_weather() {
  local city="$1"
  local lat="$2"
  local lon="$3"

  echo "Fetching weather for $city..."

  local urls=$(curl -s "https://api.weather.gov/points/${lat},${lon}" | jq -r '.properties.forecast')
  if [ -z "$urls" ]; then
    echo "Failed to retrieve forecast URL for $city. Exiting."
    return 1
  fi

  local weather_data=$(curl -s "$urls")
  if [ -z "$weather_data" ]; then
    echo "Failed to retrieve forecast data for $city. Exiting."
    return 1
  fi

  local description=$(echo "$weather_data" | jq -r '.properties.periods[0].shortForecast')
  local temp=$(echo "$weather_data" | jq -r '.properties.periods[0].temperature')
  local winddir=$(echo "$weather_data" | jq -r '.properties.periods[0].windDirection')

  if [ -z "$description" ] || [ -z "$temp" ] || [ -z "$winddir" ]; then
    echo "Failed to retrieve one or more weather parameters for $city. Exiting."
    return 1
  fi

  echo "General description for the $city area is $description"
  echo "General temperature for the $city area is $temp F"
  echo "General wind direction for the $city area is $winddir"
  echo # Newline for readability
}

# Define locations with city name, latitude, and longitude
locations=(
  "Los Angeles,34.0522,-118.2437"
  "Atlanta,33.7488,-84.3877"
  "New York City,40.7128,-74.006"
)

# Loop through each location and fetch/display weather details
for location in "${locations[@]}"; do
  IFS=',' read -r city lat lon <<< "$location"
  fetch_and_display_weather "$city" "$lat" "$lon"
done