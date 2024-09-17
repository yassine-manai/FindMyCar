# Step 1: Use official Golang image as the base image (adjusted to the actual Go version)
FROM golang:1.23-alpine

# Step 2: Set the working directory inside the container
WORKDIR /app

COPY go.mod . 
COPY go.sum .

# Step 4: Download the dependencies
RUN go mod download

# Step 5: Copy the rest of the application code to the working directory
COPY . .

# Step 6: Build the Go application
RUN go build -o bin/app .

# Step 7: Set the entry point for the container to run the application
ENTRYPOINT [ "/app/bin/app" ]
