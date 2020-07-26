# APIs developed in GO and containerized the services and orchestration using Kubernetes

## Development in APIs in Go

### Steps to build Go binary


#### 1. Go get all the packages required for development API
```
go get github.com/golang/example/hello
```
#### 2. Export Go Path to current directory
```
export GOPATH="$CURRENT_PWD_PATH"
```
#### 3. Execute go build or go install
```
go build
go install payment
```
#### 4. Execute the generated binary to have API server runnning
```
./payment
```
#### 5. Publish the Microservice as docker image using Dockerfile
```
FROM golang:latest 
EXPOSE 3000
RUN mkdir /payment
ADD . /payment/ 
WORKDIR /payment 
ENV GOPATH /payment
RUN cd /payment ; go install payment
CMD ["/payment/bin/payment"]

docker build -t adityadoshatti/payment:v9
```
#### 6. Once the docker image is ready and published need to follow the steps to create Kubernetes cluster, Follow the following steps
[Google cloud platform steps](GCP_Steps.md)