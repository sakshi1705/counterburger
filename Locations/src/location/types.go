package main

type Location struct {
	LocationName  string `json:"locationName"`
	LocationId    string `json:"locationId"`
	Zipcode       string `json:"zipcode"`
	AddressLine1  string `json:"addressLine1"`
	AddressLine2  string `json:"addressLine2"`
	City          string `json:"city"`
	State         string `json:"state"`
	Country       string `json:"country"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	Hours         string `json:"hours"`
	AcceptedCards string `json:"acceptedCards"`
	Distance      string `json:"distance"`
}
