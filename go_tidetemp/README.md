# TIDE TEMP GO APP

### what it does
Delivers **time and temperature and height** for the **lowest tide** over the **next 24 hours** to the user for various popular tide pool locations along the pacific coast.

### why?
Considering that some times (8am rather than 11pm) and temperatures (78F rather than 28F) are more ideal than others when it comes to a tide pool experience, and coupling that with the fact that tide heights (-.98 rather than 3.24) also play a large factor in the overall experience, this application gives the user information that they can utilize to see when, and decide how ideal the conditions are for tide pool enjoyment at the lowest tide point over the next 24 hours. 

### ya, but why?
General advancment of knowlege, advanced api calls, packaged functions, and concurency:
- The code contains multiple api calls that are diffently dependent on eachother. one involving an api endpoint that is parsed from an initial api call. another taking a value from an inital api call, and creating another api call based on information from that data set.
- The code seperates various functions and structures into packages to make organization, reusablity and maintainability.
- The code utilizes wait groups to manage multiple goroutines concurrently rather than sequentially, increasing speed, processing efficiency, and resource utilization.

### packages
-tidetemp-main.go (main app)
-tidepools
    -tidepools.go (list of locations)

### how it works
API solution:
1. Fetch and Parse Tide API Response: The `fetchTideAPI` function fetches the tide data, and `parseTideAPIResponse` parses this data to find the time and value of the lowest tide. The time of the lowest tide is formatted into a human-readable string (`lowestTideTimeFormatted`).

2. Fetch and Parse Weather API Response: The `fetchWeatherAPI` function fetches the general weather data, and `parseWeatherAPIResponse` parses this data to extract the URL for the hourly forecast.

3. Fetch Hourly Weather Data: The `fetchHourlyWeatherAPI` function uses the hourly forecast URL to fetch detailed hourly weather data.

4. Parse Hourly Weather Data Using Tide Time: The `parseHourlyWeatherAPIResponse` function takes the hourly weather data and the formatted time of the lowest tide (`lowestTideTimeFormatted`). It converts the lowest tide time back to a format that can be compared with the start times in the hourly weather data. The function then iterates through the hourly weather periods to find the period that matches the hour of the lowest tide time. Once a matching period is found, it extracts and returns the weather data (start time, temperature, and temperature unit) for that specific time.