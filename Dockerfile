FROM alpine:latest

WORKDIR /app

# Copy the weather-app.sh script into the container
COPY weather-app.sh /app/weather-app.sh

# Install bash, jq, and curl
RUN apk update && apk add bash jq curl

# Make the script executable
RUN chmod +x /app/weather-app.sh

# Command to run the script
CMD ["/app/weather-app.sh"]