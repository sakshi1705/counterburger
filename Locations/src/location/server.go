package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
	"io/ioutil"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/unrolled/render"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"errors"
)

var mongodb_server = "52.88.212.255:27017"
var mongodb_database = "location"
var mongodb_collection = "location"
var mongo_user = "cmpe281"
var mongo_pass = "cmpe281"
var adminDatabase = "admin"


var debug = true
var server1 = "http://riakLoadBalancer-265740929.us-west-2.elb.amazonaws.com:8000" // set in environment
var server2 = "http://riakLoadBalancer-265740929.us-west-2.elb.amazonaws.com:8000" // set in environment
var server3 = "http://riakLoadBalancer-265740929.us-west-2.elb.amazonaws.com:8000" // set in environment

type Client struct {
	Endpoint string
	*http.Client
}

type ErrorMessage struct {
	message string
}


var tr = &http.Transport{
	MaxIdleConns:       10,
	IdleConnTimeout:    30 * time.Second,
	DisableCompression: true,
}

func NewClient(server string) *Client {
	return &Client{
		Endpoint:  	server,
		Client: 	&http.Client{Transport: tr},
	}
}

func (c *Client) Ping() (string, error) {
	resp, err := c.Get(c.Endpoint + "/ping" )
	if err != nil {
		fmt.Println("[RIAK DEBUG] " + err.Error())
		return "Ping Error!", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if debug { fmt.Println("[RIAK DEBUG] GET: " + c.Endpoint + "/ping => " + string(body)) }
	return string(body), nil
}



func getSystemIp() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr).String()
	address := strings.Split(localAddr, ":")
	fmt.Println("address: ", address[0])
	return address[0]
}

func NewServerConfiguration() *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	n := negroni.Classic()
	mx := mux.NewRouter()
	initRoutes(mx, formatter)
	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})

	n.UseHandler(handlers.CORS(allowedHeaders, allowedMethods, allowedOrigins)(mx))
	return n
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func init() {

	//riak
	fmt.Println("Riak Server1:", server1 )	
	fmt.Println("Riak Server2:", server2 )	
	fmt.Println("Riak Server3:", server3 )	

		

	// Riak KV Setup	
	c1 := NewClient(server1)
	msg, err := c1.Ping( )
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Riak Ping Server1: ", msg)		
	}

	c2 := NewClient(server2)
	msg, err = c2.Ping( )
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Riak Ping Server2: ", msg)		
	}

	c3 := NewClient(server3)
	msg, err = c3.Ping( )
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Riak Ping Server3: ", msg)		
	}
}


func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/location/ping", pingHandler(formatter)).Methods("GET")
	mx.HandleFunc("/location", addLocationHandler(formatter)).Methods("POST")
	mx.HandleFunc("/locations", getAllLocationsHandler(formatter)).Methods("GET")
	mx.HandleFunc("/location/{locationId}", getLocationByIDHandler(formatter)).Methods("GET")
	mx.HandleFunc("/location/{locationId}", deleteLocationHandler(formatter)).Methods("DELETE")
	mx.HandleFunc("/location/zipcode/{zipcode}", getLocationByZipHandler(formatter)).Methods("GET")
	//riak routes
	mx.HandleFunc("/ping", pingHandlerRiak(formatter)).Methods("GET")
	mx.HandleFunc("/location/getLocation/{key}", getLocationHandler(formatter)).Methods("GET")
	
	// mx.handleFunc("/restaurant/{locationId}", updateLocationHandler(formatter)).Methods("PUT")
}


/* RIAK APIS*/

func pingHandlerRiak(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"Riak APIs version 1.0 alive!"})
	}
}


func getLocationHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		setupResponse(&w, req)
		if (*req).Method == "OPTIONS" {
			fmt.Println("PREFLIGHT Request")
			return
		}
		params := mux.Vars(req)
		fmt.Println("params")
		fmt.Println(params)
		var key string = params["key"]
		fmt.Println("Key : ", key)
		if key == ""  {
			formatter.JSON(w, http.StatusBadRequest, "Invalid Request. Cart Key Missing.")
		} else {
			c1 := NewClient(server1)
			var ord []Location
			ord, err := c1.getLocation(key)
			fmt.Println("ord")
			fmt.Println(ord)
			fmt.Println("Key ---> : ", err != nil)
			if err != nil {
				fmt.Println("err : ", err)
				formatter.JSON(w, http.StatusBadRequest, struct{ Test string }{"Cart not found!"})
			} else {
				formatter.JSON(w, http.StatusOK, ord)
			}
		}

	}
}

func (c *Client) getLocation(key string) ([]Location, error) {

	var ord_nil []Location
	resp, err := c.Get(c.Endpoint + "/buckets/bucket/keys/"+key )
	fmt.Println("resp")
	fmt.Println(resp.Body)
	fmt.Println(resp.StatusCode)
	if err != nil {
		fmt.Println("[RIAK DEBUG] ===> " + err.Error())
		return ord_nil, err
	}
	if resp.StatusCode != 200 {
		return ord_nil, errors.New("Key not found..")
		
		//return ord_nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if debug { fmt.Println("[RIAK DEBUG] GET: " + c.Endpoint + "/buckets/bucket/keys/"+key +" => " + string(body)) }
	var ord []Location
	fmt.Println("ord in getLocation")
	fmt.Println(ord)
	if err := json.Unmarshal(body, &ord); err != nil {
		fmt.Println("RIAK DEBUG] JSON unmarshaling failed: %s", err)
		return ord_nil, err
	}
	return ord, nil
}



/* MONGO APIS*/

func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		message := "Locations API Server IP: " + getSystemIp()
		formatter.JSON(w, http.StatusOK, struct{ Test string }{message})
	}
}




func getAllLocationsHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		//w.Header().Set("Content-Type", "application/json")

		session, err := mgo.Dial(mongodb_server)
		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		defer session.Close()

		if err := session.DB(adminDatabase).Login(mongo_user, mongo_pass); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		session.SetMode(mgo.Monotonic, true)
		collection := session.DB(mongodb_database).C(mongodb_collection)

		// params := mux.Vars(req)
		// var zipcode string = params["zipcode"]
		// fmt.Println(zipcode);

		var locationArray []Location
		err = collection.Find(bson.M{}).All(&locationArray)

		if locationArray == nil || len(locationArray) <= 0 {
			formatter.JSON(w, http.StatusNotFound, "No restaurants found")
		} else {
			fmt.Println("Result: ", locationArray)
			formatter.JSON(w, http.StatusOK, locationArray)
		}
	}
}

func addLocationHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// w.Header().Set("Content-Type", "application/json")
		uuidForRestaurant, _ := uuid.NewV4()
		session, err := mgo.Dial(mongodb_server)
		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		defer session.Close()

		if err := session.DB(adminDatabase).Login(mongo_user, mongo_pass); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		session.SetMode(mgo.Monotonic, true)
		collection := session.DB(mongodb_database).C(mongodb_collection)

		var location Location
		_ = json.NewDecoder(req.Body).Decode(&location)

		location.LocationId = uuidForRestaurant.String()
		fmt.Println("Locations: ", location)
		err = collection.Insert(location)
		if err != nil {
			formatter.JSON(w, http.StatusNotFound, "Error occurred. Cannot add restaurant")
			return
		}
		formatter.JSON(w, http.StatusOK, location)
	}
}

func getLocationByIDHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		//w.Header().Set("Content-Type", "application/json")

		session, err := mgo.Dial(mongodb_server)
		if err != nil {
			log.Fatal(err)
			formatter.JSON(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		defer session.Close()

		if err := session.DB(adminDatabase).Login(mongo_user, mongo_pass); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		session.SetMode(mgo.Monotonic, true)
		collection := session.DB(mongodb_database).C(mongodb_collection)

		params := mux.Vars(req)
		var locationId string = params["locationId"]
		// fmt.Println("All paramaters:", params)
		fmt.Println("location id is : ", locationId)

		var location Location
		err = collection.Find(bson.M{"locationid": locationId}).One(&location)

		if err != nil {
			formatter.JSON(w, http.StatusNotFound, "Cannot find restaurant")
			return
		} else {
			fmt.Println("Result: ", location)
			// res := json.NewEncoder(w).Encode(res)
			formatter.JSON(w, http.StatusOK, location)
		}
	}
}

func getLocationByZipHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		//w.Header().Set("Content-Type", "application/json")

		session, err := mgo.Dial(mongodb_server)
		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		defer session.Close()

		if err := session.DB(adminDatabase).Login(mongo_user, mongo_pass); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		session.SetMode(mgo.Monotonic, true)
		collection := session.DB(mongodb_database).C(mongodb_collection)

		params := mux.Vars(req)
		var zipcode string = params["zipcode"]
		fmt.Println(zipcode)

		var res []Location
		err = collection.Find(bson.M{"zipcode": zipcode}).All(&res)

		if res == nil || len(res) <= 0 {
			formatter.JSON(w, http.StatusNotFound, "Cannot find any restaurants for that zipcode")
		} else {
			fmt.Println("Result: ", res)
			formatter.JSON(w, http.StatusOK, res)
		}
	}
}

func deleteLocationHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		session, err := mgo.Dial(mongodb_server)
		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		defer session.Close()

		if err := session.DB(adminDatabase).Login(mongo_user, mongo_pass); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		session.SetMode(mgo.Monotonic, true)
		c := session.DB(mongodb_database).C(mongodb_collection)

		params := mux.Vars(req)

		var result Location
		err = c.Find(bson.M{"locationid": params["locationId"]}).One(&result)
		if err == nil {
			c.Remove(bson.M{"locationid": params["locationId"]})
			formatter.JSON(w, http.StatusOK, result)
		} else {
			formatter.JSON(w, http.StatusNotFound, "Restaurant not found for delete")
		}
	}
}
