FROM golang:1.22.2

RUN apt update && apt install -y nodejs npm

WORKDIR /app

RUN npm i tailwindcss -g 
RUN go install github.com/a-h/templ/cmd/templ@latest

COPY go.mod go.sum .
RUN go mod download

COPY package.json package-lock.json .
RUN npm i

COPY . .

ENV CGO_ENABLED=1
RUN tailwindcss --minify -i ./frontend/input.css -o ./assets/styles.css && templ generate
RUN go build -o bin .

ENTRYPOINT "/app/bin"

EXPOSE 3000