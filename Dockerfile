FROM golang:1.17.1 as build-env

RUN mkdir /vaccinator

WORKDIR /vaccinator

# <- COPY go.mod and go.sum files to the workspace
COPY go.mod . 
COPY go.sum .
#COPY run.sh .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o go/bin/vaccinator

######## Start a new stage from scratch #######
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=build-env vaccinator .

EXPOSE 8080

CMD ["./go/bin/vaccinator"]