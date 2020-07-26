package main

type BurgerItem struct {
	ItemName    string  `json:"itemName" bson:"itemName"`
	ItemId      string  `json:"itemId" bson:"itemId"`
	Price       float32 `json:"price" bson:"price"`
	ItemType    string  `json:"itemType" bson:"itemType"`
	Description string  `json:"description" bson:"description"`
}

type BurgerOrder struct {
	OrderId     string       `json:"orderId" bson:"orderId"`
	UserId      string       `json:"userId" bson:"userId"`
	OrderStatus string       `json:"orderStatus" bson:"orderStatus"`
	Order_Cart  []BurgerItem `json:"items" bson:"items"`
	TotalAmount float32      `json:"totalAmount" bson:"totalAmount"`
	// IpAddress   string  `json:"ipaddress" bson:"ipaddress"`
}

type RequiredPayload struct {
	OrderId     string  `json:"orderId" bson:"orderId"`
	UserId      string  `json:"userId" bson:"userId"`
	ItemName    string  `json:"itemName" bson:"itemId"`
	ItemId      string  `json:"itemId" bson:"itemId"`
	ItemType    string  `json:"itemType" bson:"itemType"`
	Price       float32 `json:"price" bson:"price"`
	Description string  `json:"description" bson:"description"`
}

var orders map[string]BurgerOrder
