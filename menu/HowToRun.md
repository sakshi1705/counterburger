# Golang REST API for Menu microservice

## Running Golang API locally

#### 1. Set your GOPATH to the project directoy

``` 
export GOPATH="Your Project directory"
```

- Note you might need to setup your environment before running the API

#### 2. Build your app
```
go build menu
```

#### 3. Run the app from terminal
```
./menu
```

#### 4. See where the app runs
```
[negroni] listening on :8001
```
#### 5. Test APIs using postman
```
Try ping the API

curl -X GET \
  http://localhost:8001/menu/ping \
  -H 'Postman-Token: ceca2e35-2963-4bd8-9a41-6504f9186f78' \
  -H 'cache-control: no-cache'

{
    "Test": "Burger Users API Server Working on machine: 10.0.2.20"
}

```
#### 6. Create Docker file 
```
sudo vi Dockerfile

FROM golang:latest 
EXPOSE 8001
RUN mkdir /menu 
ADD . /menu/ 
WORKDIR /menu 
ENV GOPATH /menu
RUN cd /menu ; go install menu
CMD ["/menu/bin/menu"]
```

#### 7. Build the docker image locally
```
sudo docker build -t menu .
sudo docker images
```

#### 8. Push docker image to dockerhub
```
docker push sanjna712/menu
```

## Setup MongoDB Replica Set for Sharding

#### MongoDB Sharding:

![MongoDB Shardon](https://user-images.githubusercontent.com/43122063/55284645-bbbffa80-532f-11e9-8df9-fbc63f99547f.jpg)
```
Set up two config servers, two shard server replication sets and one mongos AMI Linux Instances.
The Mongos instance is in public subnet, rest all are in private subnet.
```
#### AWS Screenshot showing all the servers.

<img width="1440" alt="Screen Shot 2019-05-03 at 6 49 53 PM" src="https://user-images.githubusercontent.com/43122063/57172409-55b60f80-6dd4-11e9-98d3-9830b2369d62.png">

#### 7. Create a Menu cluster which has 3 Menu Nodes on GCP using Kubernetes Service (set environment variables before creating)

<img width="1440" alt="Screen Shot 2019-05-03 at 6 54 16 PM" src="https://user-images.githubusercontent.com/43122063/57172448-002e3280-6dd5-11e9-8bdd-7dd56bc98cfb.png">

<img width="1440" alt="Screen Shot 2019-05-03 at 6 57 18 PM" src="https://user-images.githubusercontent.com/43122063/57172472-5602da80-6dd5-11e9-9041-381f175d735c.png">

#### 8. Create a Menu workload on GCP using Kubernetes Service

<img width="1440" alt="Screen Shot 2019-05-03 at 6 54 26 PM" src="https://user-images.githubusercontent.com/43122063/57172457-2c49b380-6dd5-11e9-97e5-5061bdddf1bf.png">

#### 9. Use Kubernetes Service to expose as a Load Balancer for the Menu cluster

<img width="1440" alt="Screen Shot 2019-05-03 at 6 59 48 PM" src="https://user-images.githubusercontent.com/43122063/57172510-c7db2400-6dd5-11e9-83fa-ef9f3f1ebd87.png">

#### 10. Test the endpoint on browser by using the ping command

```
http://34.83.82.198:3000/menu/ping
```

