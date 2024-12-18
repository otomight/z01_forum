# Set the builder image
FROM golang:1.23.1-alpine AS builder

ENV CERTOUT_FILE=server.crt
ENV KEYOUT_FILE=server.key

# Install compilation tools from Alpine (required by golang)
RUN apk add --no-cache gcc musl-dev sqlite-dev git openssl nodejs npm

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build with the rest of the files (use flags to remove debug symbols)
COPY . .
RUN npx tsc
RUN openssl req -x509 -config openssl.cnf \
		-out ${CERTOUT_FILE} -keyout ${KEYOUT_FILE}
RUN go build -ldflags="-s -w" -o main .


# Set the final image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main /app/forum.sql \
		/app/${CERTOUT_FILE} /app/${KEYOUT_FILE} ./
COPY --from=builder /app/web/static ./web/static
COPY --from=builder /app/web/templates ./web/templates

EXPOSE 80
EXPOSE 443

LABEL Name=Forum Version=0.1

CMD ["./main"]
