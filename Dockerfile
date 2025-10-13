FROM golang:1.22-alpine AS build
WORKDIR /app

# copie todo o projeto e resolva dependências
COPY . .
RUN go mod tidy

# build estático
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/app ./cmd/server

FROM gcr.io/distroless/base-debian12
ENV PORT=8080
EXPOSE 8080
COPY --from=build /bin/app /app
ENTRYPOINT ["/app"]
