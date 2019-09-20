# Attain based image
FROM golang:latest

# Setup govendor
RUN go get github.com/kardianos/govendor

WORKDIR /go/src/bendy-bot
COPY . . 

# Install dependencies
RUN govendor init
RUN govendor fetch +missing

# Execute App
RUN go run main.go

# Expose all 3 of these ports as discordgo just picks one of them randomly for some reason?
EXPOSE 443
EXPOSE 80
EXPOSE 8080
