# use official Golang ima
FROM golang:1.22.5

# set working directory
WORKDIR /app

# copy th source code
COPY . .

# download and install dependencies
RUN go get -d -v ./...

# build the app
RUN go build -o bin/kerigma

# expose the live port
EXPOSE 8080

# run executable
CMD ["./bin/kerigma"]