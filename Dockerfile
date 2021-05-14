FROM golang:latest

RUN mkdir /build
WORKDIR /build

RUN export GO111MODULE=on
RUN go get github.com/HeyParkerJ/10guys1cup-discord-bot/
RUN cd /build && git clone https://github.com/HeyParkerJ/10guys1cup-discord-bot.git

RUN cd /build/10guys1cup-discord-bot && go build

EXPOSE 8080

ENTRYPOINT [ "/build/10guys1cup-discord-bot/10guys1cup-discord-bot" ]
