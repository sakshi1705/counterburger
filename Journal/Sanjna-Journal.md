# Progress Report with Ideas
Either going for Counter Burger or Clipper Future. 

Can implement Frontend in React and use MongoDB for sharding and replication.

WOW factor - Multi tenancy with having a market place like Amazon - which has electronics and clothes and home furniture.

Can implement different types of databases such as Riak and MongoDB to try out the AKF Cube like MongoDB sharding and Riak replication.

Selected counter burger as the main project.

Started coding the frontend in ReactJS.

Decided to deploy go backend on Docker containers.

Using MongoDB as the database.

Planning to shard the Database on the basis of the type of food.

Finished the menu backend

Started with Payments frontend in React.

Created MongoDB instance on AWS which will be later replicated and sharded.

Launch Linux AMI Server 2.2 Instance
  AMI: Linux Server 2.2
  Instance Type: t2.micro
  VPC: cmpe281
  Network: public subnet
  Auto Public IP: yes
  Security Group: mongodb
  SG open ports: 22, 27017, 27018, 27019
  Key Pair: cmpe281-us-west-1.pem
  Name: config-server1

MongoDB is sharded and frontend is deployed on heroku.
