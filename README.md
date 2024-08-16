
# Tide App Shell Script

## How it Works
The Tide App Shell Script fetches tide and weather data for specified locations using external APIs. It processes this data to determine the times of the lowest and highest tides, along with the corresponding temperatures, and then outputs this information to the user and/or pushes the values to a monitoring service.

## Code Structure
The script is structured into functions that handle specific tasks:

- Fetching Tide Data: Retrieves tide information from a tide data API.
- Fetching Weather Data: Obtains weather data, including temperatures, from a weather API.
- Data Processing: Extracts and processes the relevant tide and weather information from the API responses.
- Error Handling: Checks for errors in data retrieval or processing and attempts to use fallback data sources if necessary.

## Usage
To use the Tide App, execute the script from the command line with required parameters (if any, such as location). Ensure you have the necessary permissions and environment variables set for accessing the APIs.

```./tide-app.sh```

## Output
The script outputs the following information:

Temperature at the time of the lowest tide.
Temperature at the time of the highest tide.
Times for the lowest and highest tides.

Example:

>Fetching tide details for:
>Point Lobos State Natural Reserve - 2024-07-16
>Lowest tide time: 02:14, Tide number: 0.735
>Highest tide time: 19:08, Tide number: 5.462
>Primary location failed, trying nearby location...
>Temperature at low tide: 64°F
>Temperature at high tide: 55°F

## Error Handling
The script includes basic error handling for:

- Missing or null data from APIs.
- Failure to connect to the APIs, with attempts to use a nearby location as a fallback.
- Invalid responses, with user-friendly error messages displayed.

## Future Improvements
- **Dynamic Location Input**: Allow users to input their own city name, latitude, and longitude instead of using predefined locations.
- **Extended Forecast**: Fetch and display a more detailed forecast, including chances of precipitation, humidity, etc.
- **Error Handling Enhancements**: Implement more robust error handling, including retries for network requests and user-friendly error messages.
- **Performance Optimization**: Optimize the script for faster execution, possibly by parallelizing the data fetching for different locations.
