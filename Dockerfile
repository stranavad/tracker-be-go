FROM golang:alpine

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY . .
# COPY types/*.go ./types/
# COPY db/*.go ./db/
# COPY main.go ./main.go
# COPY .env ./.env

# Build
RUN CGO_ENABLED=0 GOOS=linux go build main.go


# Run
CMD ["./main"]
