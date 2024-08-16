# Import necessary libraries
import requests  # Used to make HTTP requests
from datetime import datetime  # Used to work with dates and times

# Define a function to fetch and display tide and weather information
def fetch_and_display_tide_and_weather(park_name, city, state, lat, lon, station, fallback_lat=None, fallback_lon=None):
    # Get today's date in two different formats
    date = datetime.now().strftime("%Y-%m-%d")  # Format: YYYY-MM-DD
    today = datetime.now().strftime("%Y%m%d")  # Format: YYYYMMDD

    # Print the park name and today's date
    print(f"Fetching tide details for:\n{park_name} - {date}")

    # Prepare the URL to fetch tide data from NOAA's API
    noaa_tide_api_url = f"https://api.tidesandcurrents.noaa.gov/api/prod/datagetter?begin_date={today}&range=24&station={station}&product=predictions&datum=MLLW&time_zone=lst_ldt&interval=hilo&units=english&application=TidePool&format=json"
    # Make the HTTP request and get the response in JSON format
    tide_data = requests.get(noaa_tide_api_url).json()

    # Sort the tide predictions by their value to find the lowest and highest tides
    predictions = sorted(tide_data['predictions'], key=lambda x: float(x['v']))
    lowest_tide = predictions[0]  # The first item after sorting is the lowest tide
    highest_tide = predictions[-1]  # The last item is the highest tide

    # Extract and print the time and value of the lowest and highest tides
    print(f"Lowest tide time: {lowest_tide['t'].split(' ')[1]}, Tide number: {lowest_tide['v']}")
    print(f"Highest tide time: {highest_tide['t'].split(' ')[1]}, Tide number: {highest_tide['v']}")

    # Extract the hour from the lowest and highest tide times
    lowest_tide_hour = int(lowest_tide['t'].split(' ')[1].split(':')[0])
    highest_tide_hour = int(highest_tide['t'].split(' ')[1].split(':')[0])

    # Fetch and display weather data for the given latitude and longitude
    fetch_weather_data(lat, lon, lowest_tide_hour, highest_tide_hour, fallback_lat, fallback_lon)
    
# Define a function to fetch weather data
def fetch_weather_data(lat, lon, lowest_tide_hour, highest_tide_hour, fallback_lat, fallback_lon):
    # Make an HTTP request to get weather data URLs
    urls_response = requests.get(f"https://api.weather.gov/points/{lat},{lon}").json()
    # Extract the URL for hourly forecasts
    urls = urls_response['properties']['forecastHourly']
    # Fetch the hourly forecast data
    weather_data = requests.get(urls).json()

    try:
        # Try to extract temperature data for the times of lowest and highest tides
        temp_low = weather_data['properties']['periods'][lowest_tide_hour]['temperature']
        temp_high = weather_data['properties']['periods'][highest_tide_hour]['temperature']
    except (IndexError, KeyError):
        # If there's an error (e.g., data not available), try the fallback location
        if fallback_lat and fallback_lon:
            print("Primary location failed, trying nearby location...")
            urls_response = requests.get(f"https://api.weather.gov/points/{fallback_lat},{fallback_lon}").json()
            urls = urls_response['properties']['forecastHourly']
            weather_data = requests.get(urls).json()
            temp_low = weather_data['properties']['periods'][lowest_tide_hour]['temperature']
            temp_high = weather_data['properties']['periods'][highest_tide_hour]['temperature']

    # If temperature data was found, print it
    if temp_low and temp_high:
        print(f"Temperature at low tide: {temp_low}°F")
        print(f"Temperature at high tide: {temp_high}°F")
    else:
        # If no temperature data was found, print a message
        print("Temperature data not available.")

    # Print a message indicating that metrics can be pushed to a monitoring service
    print("Metrics available to be pushed to monitoring service.")

# List of locations to fetch tide and weather information for
locations = [
    "Point Loma Tide Pools,San Diego,CA,32.6731,-117.2425,9410170,32.7157,-117.1611",
    "Crystal Cove State Park,Laguna Beach,CA,33.5665,-117.8090,9410580,33.5427,-117.7854",
    "Leo Carrillo State Park,Malibu,CA,34.0453,-118.9358,9410230,34.0259,-118.7798",
    "Santa Rosa Island Tide Pools,Channel Islands National Park,CA,33.9950,-120.0805,9410840,34.0147,-119.6982",
    "Point Lobos State Natural Reserve,Carmel,CA,36.5159,-121.9480,9413450,36.5552,-121.9233",
    "Cape Perpetua Tide Pools,Yachats,OR,44.2811,-124.1089,9432780,44.3118,-124.1037",
    "Kalaloch Beach Tide Pools,Forks,WA,47.6136,-124.3740,9437540,47.7109,-124.4154",
    "Shi Shi Beach Tide Pools,Neah Bay,WA,48.2705,-124.6884,9443090,48.3687,-124.6252"
    # Add more locations as needed
]

# Loop through each location and fetch its tide and weather information
for location in locations:
    # Split the location string into its components
    park_name, city, state, lat, lon, station, fallback_lat, fallback_lon = location.split(',')
    # Call the function to fetch and display information for this location
    fetch_and_display_tide_and_weather(park_name, city, state, lat, lon, station, fallback_lat, fallback_lon)