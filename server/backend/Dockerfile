FROM golang:1.22-alpine3.19 as build
RUN apk add --no-cache gcc libc-dev ca-certificates && update-ca-certificates
WORKDIR /app

ENV CGO_ENABLED=0
ENV GO111MODULE=on

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o /app/main .

FROM scratch AS final
WORKDIR /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/main /app
EXPOSE 8088 8089
CMD [ "./main" ]