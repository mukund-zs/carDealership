package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func carHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		y, _ := ioutil.ReadAll(req.Body)
		var data car
		var err error
		err = json.Unmarshal(y, &data)
		if err != nil {
			log.Println("Unmarshal error", err)
		}

		//Empty entry case
		c := car{"", "", 0, "", "", enginetype{0, 0, 0}}
		if data == c {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		//Invalid year case
		if data.Year < 1980 || data.Year > time.Now().Year() {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(y))
			return
		}
		//Invalid brand case
		brand := []string{"Tesla", "BMW", "Ferrari", "Porsche", "Mercedes"}
		f := 0
		for _, b := range brand {
			if data.Brand == b {
				f = 1
				break
			}
		}
		if f == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(y))
			return
		}
		var id string
		if data.Id == "" {
			id = uuid.NewString()
		} else {
			id = data.Id
		}
		db := dbConn()
		//Normal entry case
		_, err3 := db.Exec("Insert into car values (?,?,?,?,?)", id, data.Name, data.Year, data.Brand, data.FuelType)

		//duplicate id case
		if err3 != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(y))
			return
		}
		//Sending data to engine
		_, err4 := db.Exec("Insert into engine values (?,?,?,?)", id, data.Engine.Displacement, data.Engine.Noc, data.Engine.Rng)
		w.WriteHeader(http.StatusCreated)
		if err4 != nil {
			log.Print("Error adding entries to engine", err4)
			return
		}
		_, err6 := w.Write([]byte(y))
		if err6 != nil {
			log.Print("Error while writing to response", err4)
		}
	} else if req.Method == http.MethodGet {
		brand := req.URL.Query().Get("brand")
		engine := req.URL.Query().Get("engine")
		if engine == "" || brand == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		db := dbConn()
		var data []car
		rows, err := db.Query("Select * from car where brand=?", brand)
		if err != nil {
			log.Println("Error while fetching data")
		}
		defer rows.Close()
		var i = 0
		for rows.Next() {
			var id string
			var name string
			var year int
			var brand string
			var fueltype string
			//var noc int
			//var displacement int
			//var rnge int
			err1 := rows.Scan(&id, &name, &year, &brand, &fueltype)
			var d car
			d = car{id, name, year, brand, fueltype, enginetype{0, 0, 0}}
			data = append(data, d)
			if err1 != nil {
				log.Println("Error while entering data to struct")
			}
			i++
		}
		if engine == "included" {
			i := 0
			for i < len(data) {
				var id string
				row := db.QueryRow("Select * from engine where id=?", data[i].Id)
				row.Scan(&id, &data[i].Engine.Displacement, &data[i].Engine.Noc, &data[i].Engine.Rng)
				i++
			}
		}
		y, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Error while marshaling")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(y))

	}
}

func runById(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodDelete {
		path := mux.Vars(req)
		db := dbConn()
		res, err := db.Exec("delete from car where id= ?", path["id"])
		if err != nil {
			log.Println("Error while deleting!", err)
		}
		r, _ := res.RowsAffected()
		if r == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err2 := db.Exec("delete from engine where id= ?", path["id"])
		if err2 != nil {
			log.Println("Error while deleting!", err2)
		}
		w.WriteHeader(http.StatusNoContent)

		////////////////////////////////////////////
	} else if req.Method == http.MethodPut {
		y, _ := ioutil.ReadAll(req.Body)
		var data car
		var err error
		err = json.Unmarshal(y, &data)
		if err != nil {
			log.Println("Unmarshal error", err)
		}
		path := mux.Vars(req)
		//Empty entry case
		c := car{"", "", 0, "", "", enginetype{0, 0, 0}}
		if data == c {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(y))
			return
		}

		//Invalid year case
		if data.Year < 1980 || data.Year > time.Now().Year() {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(y))
			return
		}
		//Invalid brand case
		brand := []string{"Tesla", "BMW", "Ferrari", "Porsche", "Mercedes"}
		f := 0
		for _, b := range brand {
			if data.Brand == b {
				f = 1
				break
			}
		}
		if f == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(y))
			return
		}
		db := dbConn()
		//Normal entry case
		res, err3 := db.Exec("update car set name=?,year=?,brand=?,fueltype=? where id=?", data.Name, data.Year, data.Brand, data.FuelType, path["id"])
		if err3 != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(y))
			return
		}
		r, _ := res.RowsAffected()
		if r == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(y))

			return
		}
		//Sending data to engine
		_, err4 := db.Exec("update engine set displacement=?,noc=?,rng=? where id=?", data.Engine.Displacement, data.Engine.Noc, data.Engine.Rng, path["id"])
		w.WriteHeader(http.StatusOK)
		if err4 != nil {
			log.Print("Error adding entries to engine", err4)
			return
		}
		_, err6 := w.Write([]byte(y))
		if err6 != nil {
			log.Print("Error while writing to response", err4)
		}
	} else if req.Method == http.MethodGet {
		path := mux.Vars(req)
		var data car
		db := dbConn()
		res := db.QueryRow("select * from car where id= ?", path["id"])
		res.Scan(&data.Id, &data.Name, &data.Year, &data.Brand, &data.FuelType)
		if data.Id == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		res2 := db.QueryRow("select * from engine where id= ?", path["id"])
		res2.Scan(&data.Id, &data.Engine.Displacement, &data.Engine.Noc, &data.Engine.Rng)
		y, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Error while marshaling")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(y))
	}
}
