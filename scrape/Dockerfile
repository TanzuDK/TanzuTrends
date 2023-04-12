# Use the official golang image as the base image
FROM golang:1.20

# Set the working directory inside the container
WORKDIR /app

# Copy the Go app source code to the working directory
COPY . .

# Download dependencies
RUN go mod download

# Build the Go app
RUN go build -o scrape .


# Start the app
CMD ["./scrape"]