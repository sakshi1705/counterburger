# Golang REST API for Order microservice

## Running Golang API locally

#### 0. Open a terminal

#### 1. Set your GOPATH to the project directoy

 
export GOPATH="Your Project directory"


- Note you might need to setup your environment before running the API

#### 2. Build your app

go build order


#### 3. Run the app from terminal

./order


#### 4. See where the app runs

[negroni] listening on :8000

#### 5. Test APIs using postman

Try ping the API

curl -X GET \
  http://localhost:8000/order/ping \
  -H 'Postman-Token: ceca2e35-2963-4bd8-9a41-6504f9186f78' \
  -H 'cache-control: no-cache'

{
    "Test": "Burger Order API Server Working on machine"
}


## Running the GO API in EC2 using docker

#### 1. Install Docker 

#### 2. Start Docker

sudo systemctl start docker<br>
sudo systemctl is-active docker


#### 3. Login to your docker hub account

sudo docker login


#### 4. Create Docker file 

sudo vi Dockerfile

```
FROM golang:latest 
EXPOSE 3000
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
ENV GOPATH /app
RUN cd /app ; go install order
CMD ["/app/bin/order"]
```


#### 5. Build the docker image locally

sudo docker build -t order . <br>
sudo docker images


#### 6. Push docker image to dockerhub

docker push aditi1203/order:v1


#### 7. Create Public EC2 Instance
```
Configuration:
1. AMI:             CentOS 7 (x86_64) - with Updates HVM
2. Instance Type:   t2.micro
3. VPC:             cmpe281
4. Network:         Private subnet (us-west-1c)
5. Auto Public IP:  No(default)
6. SG Open Ports:   22, 80, 8080, 3000, 8000
7. Key Pair:        cmpe281-lab1
```

#### 8. ssh to your ec2 instance from jumpbox, user name is centos

#### 9. Create docker-compose yml file (with the environment variables set up)

#### 10. Deploy go API for order sevice after pulling the image and running it

sudo docker pull aditi1203/order:v1

Run the command to create the container:

```
 sudo docker run --name order -e Server=ip-10-0-1-70.us-west-2.compute.internal:27017 
 -e Database=admin -e Collection=order -e User=cmpe281 -e Pass=cmpe281 -p 8000:8000 -d aditi1203/order:v1
```
![upload](https://user-images.githubusercontent.com/28626925/57173972-af5e1e80-6e54-11e9-8f20-e39d305cec8f.png)


#### 11. Clean Up docker environment if needed

```
docker stop order
docker rm order
docker rmi {imageid}
```
 
## Configuring MongoDB Replica Set

![mongodbcluster](https://user-images.githubusercontent.com/28626925/57173621-4c6a8880-6e50-11e9-954b-64d5e7a0deb6.png)



## Configuring Network Load Balancer

#### 1. Confugure load balancer over two private docker instances

1. Name:            networklbfinalproject
2. VPC:             CMPE281
3. Ports:           ELB Port: 80
                    Instance Port: 8080
4. Subnet:          Public Subnet
5. Security Group:  users-elb
                    Open Ports: 80
6. Health Check:    /users/test/ping
7. EC2 Instances:   Two private docker instances

#### 1. Images of Load Balancer

Network Load Balancer:

![nlb1](https://user-images.githubusercontent.com/28626925/57173607-16c59f80-6e50-11e9-84f3-0be212c5290d.png)

Target Groups of the Network Load Balancer

![nlb2](https://user-images.githubusercontent.com/28626925/57173614-378df500-6e50-11e9-95d9-e0fdf88f1454.png)
