package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"os"
	"strings"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/unrolled/render"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// const (
// 	MongoDBHosts = "ds143326.mlab.com:43326"
// 	AuthDatabase = "cmpe281"
// 	AuthUserName = "aditi1203"
// 	AuthPassword = "Aditi1203!"

// 	// TestDatabase = "goinggo"
// )

var mongodb_server = os.Getenv("Server")
var mongodb_database = os.Getenv("Database")
var mongodb_collection = os.Getenv("Collection")
var mongo_user = os.Getenv("User")
var mongo_pass = os.Getenv("Pass")

// var mongodb_server = "ds143326.mlab.com:43326"
// var mongodb_database = "cmpe281"
// var mongodb_collection = "orders"
// var mongo_user = "aditi1203"
// var mongo_pass = "Aditi1203!"
// var AWS_DB = "cmpe281"

func NewServer() *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	n := negroni.Classic()
	router := mux.NewRouter()
	initRoutes(router, formatter)
	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})

	n.UseHandler(handlers.CORS(allowedHeaders, allowedMethods, allowedOrigins)(router))
	return n
}

func getIp() string {
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

func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		message := "Burger order API Server Working on machine: "
		//  + getSystemIp()
		formatter.JSON(w, http.StatusOK, struct{ Test string }{message})
	}
}

func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/order/ping", pingHandler(formatter)).Methods("GET")
	mx.HandleFunc("/orders", getAllBurgers(formatter)).Methods("GET")
	mx.HandleFunc("/order/{orderId}", getBurgerByOrderId(formatter)).Methods("GET")
	mx.HandleFunc("/orders/{userId}", getBurgerByUserId(formatter)).Methods("GET")
	mx.HandleFunc("/order", orderBurger(formatter)).Methods("POST")
	mx.HandleFunc("/order/{orderId}", orderStatusUpdate(formatter)).Methods("PUT")
	mx.HandleFunc("/order/item/{orderId}", deleteOrderByItem(formatter)).Methods("POST")
	mx.HandleFunc("/order/{orderId}", orderDeleteByOrderId(formatter)).Methods("DELETE")
}

func getAllBurgers(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		//setup
		session, _ := mgo.Dial(mongodb_server)
		defer session.Close()
		session.SetMode(mgo.Monotonic, true)
		error := session.DB(AWS_DB).Login(mongo_user, mongo_pass)
		if error != nil {
			formatter.JSON(w, http.StatusInternalServerError, "Error: Authorization Error")
			return
		}
		c := session.DB(mongodb_database).C(mongodb_collection)

		//query
		var ordersarray []BurgerOrder
		err := c.Find(bson.M{}).All(&ordersarray)
		fmt.Println("Burger Order:", ordersarray)
		fmt.Println("Error in Burger Order:", err)
		formatter.JSON(w, http.StatusOK, ordersarray)
	}
}

func getBurgerByOrderId(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		//Setup
		session, _ := mgo.Dial(mongodb_server)
		defer session.Close()
		session.SetMode(mgo.Monotonic, true)
		error := session.DB(AWS_DB).Login(mongo_user, mongo_pass)
		if error != nil {
			formatter.JSON(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		c := session.DB(mongodb_database).C(mongodb_collection)

		//Code for finding record based on orderId
		parameter := mux.Vars(req)
		var orderid string = parameter["orderId"]
		fmt.Println("Inside get order by order id")
		fmt.Println("Input orderID: ", orderid)
		var burgerorder BurgerOrder
		ordererr := c.Find(bson.M{"orderId": orderid}).One(&burgerorder)
		if ordererr != nil {
			fmt.Println("Sorry! Error Occured")
			formatter.JSON(w, http.StatusNotFound, "Given OrderId Not Found")
			return
		}
		_ = json.NewDecoder(req.Body).Decode(&burgerorder)
		fmt.Println("Burger Order Details: ", burgerorder)
		formatter.JSON(w, http.StatusOK, burgerorder)
	}
}

func getBurgerByUserId(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		//setup
		session, _ := mgo.Dial(mongodb_server)
		defer session.Close()
		session.SetMode(mgo.Monotonic, true)
		error := session.DB(AWS_DB).Login(mongo_user, mongo_pass)
		if error != nil {
			formatter.JSON(w, http.StatusInternalServerError, "Error in Databse Connection")
			return
		}
		c := session.DB(mongodb_database).C(mongodb_collection)

		//Code for finding record
		parameter := mux.Vars(req)
		var userid string = parameter["userId"]
		fmt.Println("Inside get order by User id")
		fmt.Println("Input userId is: ", userid)
		var burgerorder BurgerOrder
		ordererror := c.Find(bson.M{"userId": userid,"orderStatus": "Active"}).One(&burgerorder)
		if ordererror != nil {
			fmt.Println("Sorry! Some error has occured")
			formatter.JSON(w, http.StatusNotFound, "Given UserId Not Found")
			return
		}
		_ = json.NewDecoder(req.Body).Decode(&burgerorder)
		fmt.Println("Result of Burger order:", burgerorder)
		formatter.JSON(w, http.StatusOK, burgerorder)
	}
}

//Post a burger order
func orderBurger(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		var user_order RequiredPayload
		_ = json.NewDecoder(req.Body).Decode(&user_order)

		//Setup
		session, error := mgo.Dial(mongodb_server)
		if error = session.DB(AWS_DB).Login(mongo_user, mongo_pass); error != nil {
			formatter.JSON(w, http.StatusInternalServerError, "Internal server Error has occured")
			return
		}
		defer session.Close()
		session.SetMode(mgo.Monotonic, true)
		c := session.DB(mongodb_database).C(mongodb_collection)

		//Code to insert record
		var ordertry BurgerOrder
		var item_input BurgerItem
		fmt.Println("argument1", user_order.Description)
		item_input.ItemId = user_order.ItemId
		item_input.ItemType = user_order.ItemType
		item_input.ItemName = user_order.ItemName
		item_input.Price = user_order.Price
		item_input.Description = user_order.Description

		record_error := c.Find(bson.M{"userId": user_order.UserId, "orderStatus": "Active"}).One(&ordertry)
		if record_error == nil {
			fmt.Println("Order Exists: Adding to existing order")
			ordertry.Order_Cart = append(ordertry.Order_Cart, item_input)
			ordertry.TotalAmount = (ordertry.TotalAmount + item_input.Price)
			c.Update(bson.M{"userId": user_order.UserId}, bson.M{"$set": bson.M{"items": ordertry.Order_Cart, "totalAmount": ordertry.TotalAmount}})
		} else {
			fmt.Println("New Order has been placed!")
			ordertry = BurgerOrder{
				UserId:  user_order.UserId,
				OrderId: uuid.NewV4().String(),
				Order_Cart: []BurgerItem{
					item_input},
				OrderStatus: "Active",
				TotalAmount: user_order.Price}
			ordererror := c.Insert(ordertry)
			if ordererror != nil {
				formatter.JSON(w, http.StatusInternalServerError, "Internal Server Error: Error in Placing order")
				return
			}
		}
		formatter.JSON(w, http.StatusOK, ordertry)
	}
}

//Delete entire order using orderId
func orderDeleteByOrderId(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		//Setup
		session, error := mgo.Dial(mongodb_server)
		if error = session.DB(AWS_DB).Login(mongo_user, mongo_pass); error != nil {
			formatter.JSON(w, http.StatusInternalServerError, "Error in Database Connection")
			return
		}
		defer session.Close()
		session.SetMode(mgo.Monotonic, true)
		c := session.DB(mongodb_database).C(mongodb_collection)

		//code to delete entire order by orderId
		parameter := mux.Vars(req)
		var orderid string = parameter["orderId"]
		fmt.Println("Inside delete order by orderId")
		fmt.Println("orderID: ", orderid)
		record_error := c.Remove(bson.M{"orderId": orderid})
		if record_error != nil {
			fmt.Println("Given orderid not present!")
			formatter.JSON(w, http.StatusNotFound, "Sorry!Order Not Found")
			return
		}
		formatter.JSON(w, http.StatusOK, "Order has been deleted: "+orderid)

	}
}

//Delete part of order using ItemId
func deleteOrderByItem(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		var user_order RequiredPayload
		_ = json.NewDecoder(req.Body).Decode(&user_order)
		parameter := mux.Vars(req)

		//Setup
		session, error := mgo.Dial(mongodb_server)
		if error = session.DB(AWS_DB).Login(mongo_user, mongo_pass); error != nil {
			formatter.JSON(w, http.StatusInternalServerError, "Error in Database Connection")
			return
		}
		defer session.Close()
		session.SetMode(mgo.Monotonic, true)
		c := session.DB(mongodb_database).C(mongodb_collection)

		//code to delete entire order by orderId

		var burgerorder BurgerOrder
		var orderId string = parameter["orderId"]
		fmt.Println("Input orderid is: ", orderId)
		fmt.Println("Input itemId is: ", user_order.ItemId)
		record_error := c.Find(bson.M{"orderId": orderId}).One(&burgerorder)
		if record_error != nil {
			fmt.Println("Given order not found")
			formatter.JSON(w, http.StatusNotFound, "Sorry!Given order not found")
			return
		}
		for i := 0; i < len(burgerorder.Order_Cart); i++ {
			if burgerorder.Order_Cart[i].ItemId == user_order.ItemId {
				burgerorder.TotalAmount = burgerorder.TotalAmount - burgerorder.Order_Cart[i].Price
				burgerorder.Order_Cart = append(burgerorder.Order_Cart[0:i], burgerorder.Order_Cart[i+1:]...)
				break
			}
		}
		c.Update(bson.M{"orderId": orderId}, bson.M{"$set": bson.M{"items": burgerorder.Order_Cart, "totalAmount": burgerorder.TotalAmount}})
		fmt.Println("Given orderId: ", user_order.OrderId)
		fmt.Println("Item deleted: ", user_order.ItemId)
		formatter.JSON(w, http.StatusOK, burgerorder)
	}
}

//Update Burger-order-status on payment
func orderStatusUpdate(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		parameter := mux.Vars(req)

		//Setup
		session, error := mgo.Dial(mongodb_server)
		if error = session.DB(AWS_DB).Login(mongo_user, mongo_pass); error != nil {
			formatter.JSON(w, http.StatusInternalServerError, "Database Connection Error")
			return
		}
		defer session.Close()
		session.SetMode(mgo.Monotonic, true)
		c := session.DB(mongodb_database).C(mongodb_collection)

		//Code to update order status
		var ordertry BurgerOrder
		var orderId string = parameter["orderId"]
		fmt.Println("Input orderid is: ", orderId)

		record_error := c.Find(bson.M{"orderId": orderId}).One(&ordertry)
		if record_error != nil {
			formatter.JSON(w, http.StatusNotFound, "Status not updated")
		}

		fmt.Println("orderStatus",ordertry)
		ordertry.OrderStatus = "Placed"
		c.Update(bson.M{"orderId": ordertry.OrderId}, bson.M{"$set": bson.M{"orderStatus": ordertry.OrderStatus}})
		formatter.JSON(w, http.StatusOK, ordertry)

	}
}
