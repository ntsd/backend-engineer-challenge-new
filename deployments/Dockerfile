# ================= builder =================
FROM golang:1.18-alpine as builder

WORKDIR /app

# download dependency
COPY go.* ./
RUN go mod download

# copy the source file
COPY . ./

# build the app
RUN go build -o ./bin/main.bin ./cmd/run/main.go

# ================= runner =================
FROM golang:1.18-alpine as runner

WORKDIR /app

# copy the binary file and env
COPY --from=builder app/bin/main.bin ./main.bin

EXPOSE 8080

CMD sh -c ./main.bin
