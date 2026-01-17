# Stage 1: Build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod ./
# COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

# Stage 2: Run
FROM alpine:latest  

# --- เพิ่มบรรทัดนี้: ติดตั้ง Certificates ---
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/main .

CMD ["./main"]