FROM golang:1.25-alpine AS build
WORKDIR /src
RUN apk add --no-cache git ca-certificates tzdata
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/moleben .

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata wget
COPY --from=build /out/moleben /app/moleben

# ВАЖНО: добавляем файлы, на которые ссылается конфиг по умолчанию
COPY ./prompt /app/prompt
COPY ./migrations /app/migrations

ENV GIN_MODE=release
EXPOSE 8080
HEALTHCHECK --interval=10s --timeout=3s --retries=3 CMD wget -qO- http://localhost:8080/health || exit 1
USER 65532:65532
ENTRYPOINT ["/app/moleben"]
