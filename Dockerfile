FROM golang:alpine3.19 as dependencies

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

FROM dependencies AS build
COPY . ./
RUN CGO_ENABLED=0 go build -o /main -ldflags=" -w -s" .

FROM golang:alpine3.19
COPY --from=build /main /main
ENV GIN_MODE=release

CMD [ "/main" ]