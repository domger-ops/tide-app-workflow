#!/bin/bash

pushgateway_url="http://pushgateway.monitoring.svc.cluster.local:9091/metrics/job/tide_app"

fetch_and_display_tide_and_weather() {
  local park_name="$1"
  local city="$2"
  local state="$3"
  local lat="$4"
  local lon="$5"
  local station="$6"
  local fallback_lat="${7:-}"
  local fallback_lon="${8:-}"

  local date=$(date -u +"%Y-%m-%d")
  local today=$(date -u +"%Y%m%d")

  echo -e "Fetching tide details for:\n$park_name - $date"

  local noaa_tide_api_url="https://api.tidesandcurrents.noaa.gov/api/prod/datagetter?begin_date=${today}&range=24&station=${station}&product=predictions&datum=MLLW&time_zone=lst_ldt&interval=hilo&units=english&application=TidePool&format=json"
  local tide_data=$(curl -s "${noaa_tide_api_url}")

  local lowest_tide_time=$(echo "$tide_data" | jq -r '.predictions | sort_by(.v | tonumber) | .[0].t' | cut -d' ' -f2)
  local highest_tide_time=$(echo "$tide_data" | jq -r '.predictions | sort_by(.v | tonumber) | reverse | .[0].t' | cut -d' ' -f2)
  local lowest_tide_value=$(echo "$tide_data" | jq -r '.predictions | sort_by(.v | tonumber) | .[0].v')
  local highest_tide_value=$(echo "$tide_data" | jq -r '.predictions | sort_by(.v | tonumber) | reverse | .[0].v')
  lowest_tide_hour=$(echo "$lowest_tide_time" | cut -d':' -f1)
  lowest_tide_time_number=$((10#$lowest_tide_hour))
  highest_tide_hour=$(echo "$highest_tide_time" | cut -d':' -f1)
  highest_tide_time_number=$((10#$highest_tide_hour))
  
  echo "Lowest tide time: $lowest_tide_time, Tide number: $lowest_tide_value"
  echo "Highest tide time: $highest_tide_time, Tide number: $highest_tide_value"

  fetch_weather_data "$lat" "$lon" "$lowest_tide_time_number" "$highest_tide_time_number" "$fallback_lat" "$fallback_lon"
}

fetch_weather_data() {
    local lat="$1"
    local lon="$2"
    local lowest_tide_time_number="$3"  
    local highest_tide_time_number="$4"
    local fallback_lat="$5"
    local fallback_lon="$6"
    local urls=$(curl -s "https://api.weather.gov/points/${lat},${lon}" | jq -r '.properties.forecastHourly')
    local weather_data=$(curl -s "$urls")

    local temp_low=$(echo "$weather_data" | jq -r --arg index "$lowest_tide_time_number" '.properties.periods[$index|tonumber].temperature')
    local temp_high=$(echo "$weather_data" | jq -r --arg index "$highest_tide_time_number" '.properties.periods[$index|tonumber].temperature')

    if [[ "$temp_low" == "null" || -z "$temp_low" || "$temp_high" == "null" || -z "$temp_high" ]]; then
      if [[ -n "$fallback_lat" && -n "$fallback_lon" ]]; then
        echo "Primary location failed, trying nearby location..."
        urls=$(curl -s "https://api.weather.gov/points/${fallback_lat},${fallback_lon}" | jq -r '.properties.forecastHourly')
        weather_data=$(curl -s "$urls")
        temp_low=$(echo "$weather_data" | jq -r --arg index "$lowest_tide_time_number" '.properties.periods[$index|tonumber].temperature')
        temp_high=$(echo "$weather_data" | jq -r --arg index "$highest_tide_time_number" '.properties.periods[$index|tonumber].temperature')
      fi
    fi

    if [[ "$temp_low" == "null" || -z "$temp_low" || "$temp_high" == "null" || -z "$temp_high" ]]; then
      echo "Temperature data not available."
    else
      echo "Temperature at low tide: ${temp_low}°F" 
      echo "Temperature at high tide: ${temp_high}°F"
    fi
    echo
}

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

for location in "${locations[@]}"; do
  IFS=',' read -r park_name city state lat lon station fallback_lat fallback_lon <<< "$location"
  fetch_and_display_tide_and_weather "$park_name" "$city" "$state" "$lat" "$lon" "$station" "$fallback_lat" "$fallback_lon"
done