# FROM golang:1.23.4

# # WORKDIR /app

# # COPY go.mod ./
# # COPY go.sum ./
# # RUN go mod download

# # COPY . .

# # RUN go mod tidy
# # RUN go build -o main ./cmd/main.go


# # CMD ["./main"]


# FROM golang:1.23.4


# WORKDIR /app

# COPY go.mod ./
# COPY go.sum ./
# RUN go mod download

# COPY . .

# RUN go mod tidy
# RUN go build -o main ./cmd/main.go

# CMD ["./main"]
