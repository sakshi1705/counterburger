#### [Reference link used for App in private subnet](https://medium.com/@manoj.bhagwat60/deploying-the-application-google-cloud-platform-public-private-subnets-9b77986ebd51)

##### Firewall rule for allowing SSH
```
gcloud compute firewall-rules create cmpe281-vpc-allow-ssh --allow tcp:22 --network cmpe281
```
##### Allow other ports for internal machines
```
gcloud compute firewall-rules create cmpe281-vpc-allow-internal-network \
--allow tcp:1-65535,udp:1-65535,icmp \
--source-ranges 10.0.0.0/16 \
--network cmpe281
```
##### Create a NAT Gateway 
```
gcloud compute instances create nat-gateway --network cmpe281 --can-ip-forward \
  --zone us-west1-a \     
  --image-family debian-9\     
  --image-project debian-cloud \
  --subnet public \
  --tags nat-instance\
```
##### Create an instance for Mongo DB
```
gcloud compute instances create mongo1 --network cmpe281 --no-address \
    --zone us-west1-b \
    --image-family ubuntu-1804-lts \
    --subnet private \
    --image-project gce-uefi-images \
    --tags private-instance
```

##### Change routes of all private instances through NAT gateway instance

```
gcloud compute routes create cmpe281-vpc-no-ip-internet-route --network cmpe281 \
    --destination-range 0.0.0.0/0 \
    --next-hop-instance nat-gateway \
    --next-hop-instance-zone us-west1-a \
    --tags private-instance --priority 800
```

#### On your NAT instance, configure iptables:
```
sudo sysctl -w net.ipv4.ip_forward=1
sudo iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
```


_______________________________________________________________________________________________

[GCP Documentation for Privtae Kubernetes Cluster](https://cloud.google.com/kubernetes-engine/docs/how-to/private-clusters)
[GCP Documentation for Kubernetes cluster command line](https://cloud.google.com/kubernetes-engine/docs/tutorials/hello-app)
[GCP Documentation for pulling docker images in private Kubernetes Cluster](https://cloud.google.com/kubernetes-engine/docs/how-to/private-clusters#pulling_images_from_docker_hub)

##### Get Ip address of Cloud shell
```
dig +short myip.opendns.com @resolver1.opendns.com
```

##### Update the cluster to access through this cloud shell

```
gcloud container clusters update payment \
    --zone us-west1-b \
    --enable-master-authorized-networks \
    --master-authorized-networks 100.0.0.0/28,104.196.231.234/32
```

##### Get credentails to access the cluster
```
gcloud container clusters get-credentials payment \
    --zone us-west1-b \
    --project payment-238702
```
##### Change the network tag of all Kubernetes cluster instances
```
gcloud compute instances add-tags gke-payment-default-pool-86ea36a4-x1dp \
      --zone us-west1-b \
      --tags private-instance

gcloud compute instances add-tags gke-payment-default-pool-86ea36a4-jfvj \
      --zone us-west1-b \
      --tags private-instance

gcloud compute instances add-tags gke-payment-default-pool-86ea36a4-bd1x \
      --zone us-west1-b \
      --tags private-instance
 ```
 ##### Run deployment inside the kubernetes cluster
 ```
kubectl run payment-web --image=adityadoshatti/payment:v7 --port 3000 \
  --env "AWS_MONGODB=34.216.42.35:27017" --env "MONGODB_DBNAME=payments" \
  --env "MONGODB_COLLECTION=payments_details" --env "MONGODB_USERNAME=cmpe281" \
  --env "MONGODB_PASSWORD=****"
```

kubectl scale deployment payment-web --replicas=3
