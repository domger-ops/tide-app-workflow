#!/bin/bash

# Pushgateway URL
pushgateway_url="http://pushgateway.monitoring.svc.cluster.local:9091/metrics/job/tide_app"

# Function to fetch and display tide and weather details for a given location using the Stormglass API
fetch_and_display_tide_and_weather() {
  local park_name="$1"
  local city="$2"
  local state="$3"
  local lat="$4"
  local lon="$5"
  local station="$6"
  local fallback_lat="${7:-}" # Fallback latitude, optional
  local fallback_lon="${8:-}" # Fallback longitude, optional

  # Current date in YYYYMMDD format
  local date=$(date -u +"%Y-%m-%d") # Format: YYYY-MM-DD
  local today=$(date -u +"%Y%m%d") # Format: YYYYMMDD, used for API calls

  # Display location details
  echo -e "Fetching lowest tide details for $date:\n$park_name"

  # Modify NOAA API URL to use the station variable
  local noaa_tide_api_url="https://api.tidesandcurrents.noaa.gov/api/prod/datagetter?begin_date=${today}&range=24&station=${station}&product=predictions&datum=MLLW&time_zone=lst_ldt&interval=hilo&units=english&application=TidePool&format=json"

  # Fetch tide data using the NOAA API
  local tide_data=$(curl -s "${noaa_tide_api_url}")

  # Extract tide information from the JSON response
  # Value: Lowest Tide Time (HH:MM) - eliminates the date, allowing for easier use of value when needed for weather api parsing
  local lowest_tide_time=$(echo "$tide_data" | jq -r '.predictions | sort_by(.v | tonumber) | .[0].t' | cut -d' ' -f2)

  # Value: Lowest Tide Number
  local lowest_tide_value=$(echo "$tide_data" | jq -r '.predictions | sort_by(.v | tonumber) | .[0].v')

  # Assuming lowest_tide_time is in the format "HH:MM"
  lowest_tide_hour=$(echo "$lowest_tide_time" | cut -d':' -f1) # Extracts the hour part
  lowest_tide_time_number=$((10#$lowest_tide_hour)) # Converts the hour to a number to use for weather API

  # Display tide details
  echo "Lowest tide time: $lowest_tide_time, Tide number: $lowest_tide_value"

  # Fetch and display weather details using the National Weather Service API
  fetch_weather_data() {
    local lat="$1"
    local lon="$2"
    local urls=$(curl -s "https://api.weather.gov/points/${lat},${lon}" | jq -r '.properties.forecastHourly')
    local weather_data=$(curl -s "$urls")
    local temp=$(echo "$weather_data" | jq -r --arg index "$lowest_tide_time_number" '.properties.periods[$index|tonumber].temperature')
    echo "$temp"
  }

  local temp=$(fetch_weather_data "$lat" "$lon")
  if [ "$temp" = "null" ] || [ -z "$temp" ] && [ -n "$fallback_lat" ] && [ -n "$fallback_lon" ]; then
    echo "Primary location failed, trying nearby location..."
    temp=$(fetch_weather_data "$fallback_lat" "$fallback_lon")
  fi

  if [ "$temp" = "null" ] || [ -z "$temp" ]; then
    echo "Temperature data not available."
  else
    echo "Temperature: ${temp}°F"
  fi
  echo

  # Push metrics to Pushgateway
  echo "Pushing metrics to Pushgateway..."
  curl -X POST --data-binary @- "${pushgateway_url}" <<EOF
tide_lowest_value{park_name="$park_name", city="$city", state="$state"} $lowest_tide_value
tide_lowest_time{park_name="$park_name", city="$city", state="$state"} $lowest_tide_time
weather_temperature{park_name="$park_name", city="$city", state="$state"} $temp
EOF

  echo "Metrics pushed successfully."
}

# Array of locations
locations=(
  "Point Loma Tide Pools,San Diego,CA,32.6731,-117.2425,9410170,32.7157,-117.1611"
  "Crystal Cove State Park,Laguna Beach,CA,33.5665,-117.8090,9410580,33.5427,-117.7854"
  "Leo Carrillo State Park,Malibu,CA,34.0453,-118.9358,9410230,34.0259,-118.7798"
  "Santa Rosa Island Tide Pools,Channel Islands National Park,CA,33.9950,-120.0805,9410840,34.0147,-119.6982"
  "Point Lobos State Natural Reserve,Carmel,CA,36.5159,-121.9480,9413450,36.5552,-121.9233"
  "Cape Perpetua Tide Pools,Yachats,OR,44.2811,-124.1089,9432780,44.3118,-124.1037"
  "Kalaloch Beach Tide Pools,Forks,WA,47.6136,-124.3740,9437540,47.7109,-124.4154"
  "Shi Shi Beach Tide Pools,Neah Bay,WA,48.2705,-124.6884,9443090,48.3687,-124.6252"
)

# Loop through each location and fetch/display tide and weather details
for location in "${locations[@]}"; do
  IFS=',' read -r park_name city state lat lon station fallback_lat fallback_lon <<< "$location"
  fetch_and_display_tide_and_weather "$park_name" "$city" "$state" "$lat" "$lon" "$station" "$fallback_lat" "$fallback_lon"
done