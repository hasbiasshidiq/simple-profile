# Dockerfile definition for Backend application service.

# From which image we want to build. This is basically our environment.
FROM golang:1.19-alpine as Build

# Set the working directory to /app
WORKDIR /app

# This will copy all the files in our repo to the inside the container at root location.
COPY . .


# Build our binary at /app
RUN go build -o /app/main cmd/main.go

####################################################################
# This is the actual image that we will be using in production.
FROM alpine:latest

# Set the working directory to /app
WORKDIR /app

# Copy the binary from the build image to the production image.
COPY --from=Build /app/main .

# Copy the /cert directory from the build image to the production image.
COPY --from=Build /app/cert ./cert

# This is the port that our application will be listening on.
EXPOSE 1323

# This is the command that will be executed when the container is started.
ENTRYPOINT ["./main"]