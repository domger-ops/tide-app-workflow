FROM alpine:latest

WORKDIR /app

# Copy the tide-app.sh script into the container
COPY tide-app.sh /app/tide-app.sh

# Install bash, jq, and curl
RUN apk update && apk add bash jq curl

# Make the script executable
RUN chmod +x /app/tide-app.sh

# Command to run the script
CMD ["/app/tide-app.sh"]