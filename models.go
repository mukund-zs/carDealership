package main

type enginetype struct {
	Displacement int `json:"displacement"`
	Noc          int `json:"noc"`
	Rng          int `json:"rng"`
}
type car struct {
	Id       string     `json:"id"`
	Name     string     `json:"name"`
	Year     int        `json:"year"`
	Brand    string     `json:"brand"`
	FuelType string     `json:"fuelType"`
	Engine   enginetype `json:"engine"`
}

//uuid.UUID
