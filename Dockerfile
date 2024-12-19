# Set the builder image
FROM golang:1.23.1-alpine AS builder

ARG CERTOUT_FILE
ARG KEYOUT_FILE
ARG MAIN_SCSS_FILE
ARG MAIN_CSS_OUT_FILE

# Install compilation tools from Alpine (required by golang)
RUN apk add --no-cache gcc musl-dev sqlite-dev git openssl nodejs npm

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build with the rest of the files (use flags to remove debug symbols)
COPY . .
RUN npx tsc
RUN npx sass ${MAIN_SCSS_FILE}:${MAIN_CSS_OUT_FILE} --style=compressed
RUN openssl req -x509 -config openssl.cnf \
		-out ${CERTOUT_FILE} -keyout ${KEYOUT_FILE}
RUN go build -ldflags="-s -w" -o main .


# Set the final image
FROM alpine:latest

ARG CERTOUT_FILE
ARG KEYOUT_FILE

WORKDIR /app

COPY --from=builder /app/main /app/forum.sql \
		/app/${CERTOUT_FILE} /app/${KEYOUT_FILE} ./
COPY --from=builder /app/web/static ./web/static
COPY --from=builder /app/web/templates ./web/templates

EXPOSE 80
EXPOSE 443

LABEL Name=Forum Version=0.1

CMD ["./main"]
