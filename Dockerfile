# Attain based image
FROM golang:latest AS build-env

# Setup govendor
WORKDIR /go/src/bendy-bot

# Copy over important go mod stuff
COPY go.mod go.sum ./

# Retreive dependencies
RUN go mod download
# copy files
COPY internal/ internal/
COPY main.go main.go

# Compile the application
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bendy-bot .


# Second stage of building to reduce image size
FROM alpine:latest

WORKDIR  /go/src/bendy-bot
COPY --from=build-env /go/src/bendy-bot/bendy-bot . 
COPY --from=build-env /go/src/bendy-bot/internal/storage internal/storage/

# Set gopath
ENV GOPATH /go


# Expose all 3 of these ports as discordgo just picks one of them randomly for some reason?
EXPOSE 443
EXPOSE 80
EXPOSE 8080

CMD ["./bendy-bot"]