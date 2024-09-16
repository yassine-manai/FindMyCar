# Step 1: Use official Golang image as the base image (adjusted to the actual Go version)
FROM golang:1.23-alpine

# Step 2: Set the working directory inside the container
WORKDIR /app

# Step 3: Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Step 4: Download the dependencies
RUN go mod download

# Step 5: Copy the rest of the application code to the working directory
COPY . .

# Step 6: Disable CGO for compatibility with Alpine and build the Go application
RUN CGO_ENABLED=0 go build -o main .



# Step 8: Command to run the Go application
CMD ["./main"]
