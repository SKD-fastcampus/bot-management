FROM golang:1.23 as builder

WORKDIR /app

COPY go.work go.work.sum ./
COPY pkg pkg
COPY services/bot-mgmt-server services/bot-mgmt-server
COPY services/geo services/geo 
# Copy geo just in case go.work needs it or references exist, but ideally only bot-mgmt-server is needed
# However, go.work requires all modules to be present.
COPY tools tools

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server services/bot-mgmt-server/cmd/server/main.go

FROM gcr.io/distroless/static-debian12

COPY --from=builder /bin/server /server

EXPOSE 8080

CMD ["/server"]
