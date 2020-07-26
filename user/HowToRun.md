# GOLANG REST API for user microservice

## Running GOLANG API locally

#### 0. Open a terminal

#### 1. Set your GOPATH to the project directoy

``` 
export GOPATH="Your Project directory"
```

- Note you might need to setup your environment before running the API

#### 2. Build your app
```
go build users
```

#### 3. Run the app from terminal
```
./users
```

#### 4. See where the app runs
```
[negroni] listening on :8080
```
#### 5. Test APIs using postman
```
Try ping the API

curl -X GET \
  http://localhost:8000/users/test/ping \
  -H 'Postman-Token: ceca2e35-2963-4bd8-9a41-6504f9186f78' \
  -H 'cache-control: no-cache'

{
    "Test": "Burger Users API Server Working on machine: 10.0.1.29"
}

```
## Running the GO API in EC2 using docker

#### 1. Install Docker 

#### 2. Start Docker
```
sudo systemctl start docker
sudo systemctl is-active docker
```

#### 3. Login to your docker hub account
```
sudo docker login
```

#### 4. Create Docker file 
```
sudo vi Dockerfile

FROM golang:latest 
EXPOSE 8080
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
ENV GOPATH /app
RUN cd /app ; go install users
CMD ["/app/bin/users"]
```

#### 5. Build the docker image locally
```
sudo docker build -t users .
sudo docker images
```

#### 6. Push docker image to dockerhub
```
docker push ramyaiyer06/users:v2
```

#### 7. Create Public EC2 Instance

Configuration:
1. AMI:             CentOS 7 (x86_64) - with Updates HVM
2. Instance Type:   t2.micro
3. VPC:             cmpe281
4. Network:         Private subnet (us-west-1c)
5. Auto Public IP:  No(default)
6. SG Open Ports:   22, 80, 8080, 3000, 8000
7. Key Pair:        cmpe281-lab1

#### 8. ssh to your ec2 instance from jumpbox, user name is centos

#### 9. Create docker-compose yml file (with the environment variables set up)

#### 10. Deploy go API for order sevice after pulling the image and running it
```
sudo docker pull ramyaiyer06/users:v2
docker run --name users -td -p 8080:8080 ramyaiyer06/users:v2
```

#### 11. Clean Up docker environment when finished
```
docker stop users
docker rm users
docker rmi {imageid}
``` 
## Configuring RDS for MySQL
#### 1. select RDS service from AWS Console and create a database with configuration as follows:

```
Version: 5.5.61
Instance Type: "db.t2.micro"
Instance Name: user
DB User: user
DB Pass: cmpe281pass
VPC: cmpe281
DB Subnet Group: mysql
Public accessibility: No
AZ Info: No Preference
Security Group: rds-mysql (with port 3306 open)
DB Name: user
Port: 3306
IAM DB authentication: Disable
Disable auto minor version upgrade
Remaining Options: Use Defaults
```

#### 2. Create a replica for RDS with name "user-replica"

#### 3. RDS Dashboard images
![Screen Shot 2019-05-01 at 5 08 15 PM](https://user-images.githubusercontent.com/43103509/57050978-fde0a280-6c33-11e9-81fe-f1087da537dd.png)
![Screen Shot 2019-05-01 at 5 08 54 PM](https://user-images.githubusercontent.com/43103509/57050997-151f9000-6c34-11e9-92b5-7a2603d82cf8.png)

## Configuring Load Balancer

#### 1. Confugure load balancer over two private docker instances
```
1. Name:            users-classic-public
3. VPC:             CMPE281
4. Ports:           ELB Port: 80
                    Instance Port: 8080
6. Subnet:          Public Subnet
7. Security Group:  users-elb
                    Open Ports: 80
8. Health Check:    /users/test/ping
9. EC2 Instances:   Two private docker instances
```
#### 1. Images of Load Balancer

![Screen Shot 2019-05-01 at 5 21 17 PM](https://user-images.githubusercontent.com/43103509/57051258-8f044900-6c35-11e9-8b9e-d23c7d02e339.png)
