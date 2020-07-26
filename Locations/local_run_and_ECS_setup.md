### Instructions to run the restaurant location application locally ### <br>
<br>
1. Go to the root directory of the project and execute a pwd command <br>
2. Set go path for proper execution <br>

```
export GOPATH=output/of/pwd/command 
export PATH=$PATH:$GOPATH/bin
```
<br>
3. Run the go files

```
go run *.go
````
<br>
4. Running the above command should give something as follows

```
Riak Server1: http://riakDNSaddress.us-west-2.elb.amazonaws.com:8000
Riak Server2: http://riakDNSaddress.us-west-2.elb.amazonaws.com:8000
Riak Server3: http://riakDNSaddress.us-west-2.elb.amazonaws.com:8000
[RIAK DEBUG] GET: http://riakLoadBalancer-265740929.us-west-2.elb.amazonaws.com:8000/ping => OK
2019/05/03 02:24:02 Riak Ping Server1:  OK
[RIAK DEBUG] GET: http://riakLoadBalancer-265740929.us-west-2.elb.amazonaws.com:8000/ping => OK
2019/05/03 02:24:02 Riak Ping Server2:  OK
[RIAK DEBUG] GET: http://riakLoadBalancer-265740929.us-west-2.elb.amazonaws.com:8000/ping => OK
2019/05/03 02:24:02 Riak Ping Server3:  OK
[negroni] listening on :8000
```
<br>
5. Ping the application to check if it is working by opening browser and going to http://localhost:8000/ping. it should show the following

```
{
  "Test": "Location API ping Handler working!!"
}
```

## Setting up ECS Cluster for location service ##

Below are the steps used by us to setup and run our ECS cluster for our Locations microservice

1. Of the various options provided, we choose the appropriate container image template 

![containerDefinition](https://user-images.githubusercontent.com/47696913/57172610-3c629280-6dd7-11e9-866b-8e756f0bb377.png)

2. After choosing the template, we define the container in container definition pane

![containerDefinition1](https://user-images.githubusercontent.com/47696913/57172615-42f10a00-6dd7-11e9-9382-c635f16bfec6.png)

3. We also need to define port mappings and health checks 

![containerDefinition2](https://user-images.githubusercontent.com/47696913/57172618-484e5480-6dd7-11e9-80dc-eef14d34ba87.png)

4. After defining the container, we define the task. We enter details for task definitions

![taskDefinition](https://user-images.githubusercontent.com/47696913/57172622-4e443580-6dd7-11e9-92dc-922e328871a7.png)

5. Next we define the service definition where we mention the number of the name of the service and number of tasks.

![servicedefinition](https://user-images.githubusercontent.com/47696913/57172624-51d7bc80-6dd7-11e9-8fac-9de369b7b5fd.png)

6. Then we specify that we want to use load balancer for our tasks

![serviceDefinition](https://user-images.githubusercontent.com/47696913/57172626-543a1680-6dd7-11e9-8316-77b5d3d38ffb.png)

7. Finally, we define the cluster at the end

![clusterDefinition](https://user-images.githubusercontent.com/47696913/57172630-5d2ae800-6dd7-11e9-9f2a-9e2be5761d84.png)
