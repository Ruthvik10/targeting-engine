FROM golang:1.23-alpine AS BuildStage
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /bin/targeting-engine ./main.go
EXPOSE 8080

FROM alpine:latest
WORKDIR /
COPY --from=BuildStage /bin/targeting-engine /bin/targeting-engine
COPY --from=BuildStage /app/app.env /
CMD ["/bin/targeting-engine"]