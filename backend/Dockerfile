FROM golang

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o /main ./cmd/app/main.go

CMD [ "/main" ]