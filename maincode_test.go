package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCarPost(t *testing.T) {
	db, _ = sql.Open("mysql", "root:root@123@tcp(localhost:3306)/cardealership")

	const UUID = "87184816-3d9b-4a39-94b1-55b320e0ef41" //for duplicate id case

	carInput := car{"", "Ferrari11", 2019, "Ferrari", "Petrol", enginetype{3, 5, 100}}
	carInput2 := car{UUID, "model 1", 2019, "Tesla", "Diesel", enginetype{3, 4, 100}}
	invalidBrand := car{"", "model 1", 2019, "dvsjc", "Diesel", enginetype{3, 4, 100}}
	emptyCar := car{}
	invalidIdInput := car{"", "model 11", 1970, "Ferrari", "Diesel", enginetype{3, 4, 100}}
	//invalidIdInput2 := car{"sdjfds", "model 11", 1999, "Ferrari", "Diesel", enginetype{3, 4, 100}}

	testcases := []struct {
		desc   string
		input  car
		output car
		status int
	}{
		{"normal entry case", carInput, carInput, http.StatusCreated},
		{"empty entry case", emptyCar, emptyCar, http.StatusBadRequest},
		{"Invalid brand", invalidBrand, invalidBrand, http.StatusBadRequest},
		{"duplicate id case", carInput2, carInput2, http.StatusInternalServerError},
		{"Invalid year", invalidIdInput, invalidIdInput, http.StatusBadRequest},
		//{"Invalid UUID", invalidIdInput2, emptyCar, http.StatusBadRequest},
	}
	for i, v := range testcases {
		body, err := json.Marshal(v.input)
		if err != nil {
			t.Errorf("Marshall err,test:%v failed", i)

			continue
		}
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
		w := httptest.NewRecorder()
		carHandler(w, req)
		res := w.Result()
		var d car
		respBody, _ := ioutil.ReadAll(res.Body)
		err = json.Unmarshal(respBody, &d)
		//log.Println(d)
		//log.Println(v.input)
		assert.Equal(t, v.status, w.Result().StatusCode, "Test case failed %v", i)
		assert.Equal(t, v.input, d, "Test case failed %v", i)
	}
}
func TestCarDelete(t *testing.T) {
	var UUID string = "246cb94d-2ce0-4394-9957-817a71e67f0a"
	testcases := []struct {
		desc   string
		id     string
		status int
	}{
		{"passed id to delete", UUID, http.StatusNoContent},
		{"empty entry sent to delete", "", http.StatusBadRequest},
		{"Invalid UUID", "fduygs", http.StatusBadRequest},
		{"Valid UUID not present in database", uuid.NewString(), http.StatusBadRequest},
	}
	for i, tc := range testcases {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req = mux.SetURLVars(req, map[string]string{"id": tc.id})
		w := httptest.NewRecorder()
		runById(w, req)
		if tc.status != w.Result().StatusCode {
			t.Errorf("Test Failed %v.", i)
		}
	}
}

func TestCarPut(t *testing.T) {
	const UUID = "87184816-3d9b-4a39-94b1-55b320e0ef41"
	input1 := car{UUID, "model 1", 2019, "Tesla", "Diesel", enginetype{3, 4, 100}}
	emptyCar := car{"", "", 0, "", "", enginetype{0, 0, 0}}
	log.Println(emptyCar)
	input2 := car{uuid.NewString(), "model 1", 2019, "Tesla", "Diesel", enginetype{3, 4, 100}}
	invalidYearInput := car{UUID, "model 1", 1909, "Tesla", "Diesel", enginetype{3, 4, 100}}
	invalidBrand := car{UUID, "model 1", 2009, "Tla", "Diesel", enginetype{3, 4, 100}}
	testcases := []struct {
		desc   string
		input  car
		output car
		status int
	}{
		{"normal update case", input1, input1, http.StatusOK},
		{"empty data", emptyCar, emptyCar, http.StatusBadRequest},
		{"Valid UUID but not present in database", input2, input2, http.StatusBadRequest},
		{"Invalid brand", invalidBrand, invalidBrand, http.StatusBadRequest},
		{"Invalid year", invalidYearInput, invalidYearInput, http.StatusBadRequest},
	}
	for i, v := range testcases {
		body, err := json.Marshal(v.input)
		if err != nil {
			t.Errorf("Marshall err,test:%v failed", i)

			continue
		}
		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": v.input.Id})
		w := httptest.NewRecorder()
		runById(w, req)
		res := w.Result()
		respBody, _ := ioutil.ReadAll(res.Body)
		assert.Equal(t, v.status, w.Result().StatusCode, "Test case failed %v", i)
		assert.Equal(t, body, respBody, "Test case failed %v", i)
	}
}
func TestCarGetById(t *testing.T) {
	const UUID = "1da6c99c-b4f9-419a-a209-c8d77f193f11"
	output1 := car{UUID, "model 11", 2019, "Tesla", "Diesel", enginetype{3, 4, 100}}
	emptyCar := car{}
	testcases := []struct {
		desc   string
		id     string
		output car
		status int
	}{
		{"Valid UUID present in database", UUID, output1, http.StatusOK},
		{"Valid UUID not present in database", uuid.NewString(), emptyCar, http.StatusBadRequest},
		{"Invalid UUID", "bhxh", emptyCar, http.StatusBadRequest},
		{"Invalid UUID", "2324", emptyCar, http.StatusBadRequest},
	}
	for i, v := range testcases {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		w := httptest.NewRecorder()
		runById(w, req)
		var data car
		d, _ := ioutil.ReadAll(w.Body)
		err := json.Unmarshal(d, &data)
		if err != nil {
			log.Println(err)
		}
		assert.Equal(t, v.status, w.Result().StatusCode, "Test case failed %v", i)
		assert.Equal(t, v.output, data, "Test case failed %v", i)
	}
}

func TestCarGetByBrand(t *testing.T) {
	output1 := []car{{"e11bb37d-1a30-471f-9d23-ea8dd40e1c87", "Ferrari11", 2019, "Ferrari", "Petrol", enginetype{3, 5, 100}}}
	//output2 := []car{{"e11bb37d-1a30-471f-9d23-ea8dd40e1c87", "Ferrari11", 2019, "Ferrari", "Petrol", enginetype{0, 0, 0}}}

	var emptyCar []car

	testcases := []struct {
		desc   string
		url    string
		output []car
		status int
	}{
		{"Normal get by brand", "/?brand=Ferrari&engine=included", output1, http.StatusOK},
		{"passed empty brand", "/?brand=", emptyCar, http.StatusBadRequest},
		{"passed invalid brand", "/?brand=fdgfg", emptyCar, http.StatusBadRequest},
		{"Normal get by brand without engine", "/?brand=Ferrari", emptyCar, http.StatusBadRequest},
		{"passed empty engine", "/?brand=ferrari&engine=", emptyCar, http.StatusBadRequest},
	}
	for i, v := range testcases {
		req := httptest.NewRequest(http.MethodGet, v.url, nil)
		w := httptest.NewRecorder()
		carHandler(w, req)
		var data []car
		d, _ := ioutil.ReadAll(w.Body)
		err := json.Unmarshal(d, &data)
		if err != nil {
			log.Println(err)
		}
		assert.Equal(t, v.output, data, "Test case failed %v", i)
	}
}
