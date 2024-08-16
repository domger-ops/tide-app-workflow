# Tide and Weather Information Fetcher

## Overview

The Tide and Weather Information Fetcher is a Python application designed to provide users with detailed tide and weather information for various parks and locations. This tool fetches tide data from the NOAA's API and weather data from the Weather.gov API, presenting the lowest and highest tides of the day along with the corresponding weather conditions.

## Features

- **Tide Information**: Fetches and displays the times and values of the lowest and highest tides for the specified date and location.
- **Weather Information**: Provides temperature data for the times of the lowest and highest tides, enhancing planning and preparation for visitors.
- **Fallback Location Support**: In case the primary location fails to fetch weather data, the application attempts to retrieve data for a nearby fallback location.
- **User-Friendly Output**: Presents the fetched data in a clear, readable format, making it easy for users to understand and use the information.

## Prerequisites

Before you begin, ensure you have met the following requirements:

- **Python 3.6+**: The application is written in Python and requires Python 3.6 or newer.
- **Requests Library**: This project uses the [`requests`](command:_github.copilot.openSymbolFromReferences?%5B%7B%22%24mid%22%3A1%2C%22path%22%3A%22%2FLibrary%2FFrameworks%2FPython.framework%2FVersions%2F3.11%2Flib%2Fpython3.11%2Fsite-packages%2Frequests%2F__init__.py%22%2C%22scheme%22%3A%22file%22%7D%2C%7B%22line%22%3A0%2C%22character%22%3A0%7D%5D "../../../../Library/Frameworks/Python.framework/Versions/3.11/lib/python3.11/site-packages/requests/__init__.py") library to make HTTP requests. You can install it using pip:

  ```bash
  pip install requests
  ```

## Installation

To install the Tide and Weather Information Fetcher, follow these steps:

1. Clone the repository to your local machine:

   ```bash
   git clone https://github.com/your-username/tide-weather-fetcher.git
   ```

2. Navigate to the cloned repository:

   ```bash
   cd tide-weather-fetcher
   ```

3. Ensure you have the required dependencies installed:

   ```bash
   pip install -r requirements.txt
   ```

## Usage

To use the Tide and Weather Information Fetcher, you need to execute the script with Python. The script does not require any command-line arguments as it is currently set to fetch data for predefined locations.

```bash
python tidetemp_python.py
```

Upon execution, the script will fetch and display tide and weather information for the locations defined within the script.

## Configuration

The script includes a list of locations for which it fetches tide and weather information. To add or modify locations, edit the `locations` list in the script. Each location should be specified as a string with the following format:

```python
"Park Name,City,State,Latitude,Longitude,NOAA Station ID,Fallback Latitude,Fallback Longitude"
```

## Contributing

Contributions to the Tide and Weather Information Fetcher are welcome. To contribute:

1. Fork the repository.
2. Create a new branch for your feature (`git checkout -b feature/AmazingFeature`).
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4. Push to the branch (`git push origin feature/AmazingFeature`).
5. Open a pull request.

## Contact

For any queries or further assistance, please contact the repository owner.

---

This README provides a comprehensive guide to installing, configuring, and using the Tide and Weather Information Fetcher. Adjust the repository URL, contact information, and any other specific details as necessary to match your project's setup.