# Set the builder image
FROM golang:1.23.1-alpine AS builder

# Install compilation tools from Alpine (required by golang)
RUN apk add --no-cache gcc musl-dev sqlite-dev git openssl nodejs npm

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build with the rest of the files (use flags to remove debug symbols)
COPY . .
RUN npx tsc
RUN openssl req -x509 -nodes -days 365 -newkey rsa:4096 \
		-keyout server.key -out server.crt \
		-subj "/C=FR/ST=Normandie/L=Rouen/O=Zone01/OU=P12024/CN=localhost"
RUN go build -ldflags="-s -w" -o main .


# Set the final image
FROM alpine:latest

# Certificate HTTPS
# RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/main /app/forum.sql /app/server.crt /app/server.key ./
COPY --from=builder /app/web/static ./web/static
COPY --from=builder /app/web/templates ./web/templates

EXPOSE 443

LABEL Name=Forum Version=0.1

CMD ["./main"]
