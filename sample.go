package main

import (
	"encoding/json"
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"fmt"
	"io/ioutil"
	"strconv"
)

type Person struct {
	ID      string   `json:"id,omitempty"`
	FIRST   string   `json:"firstname,omitempty"`
	LAST    string   `json:"lastname,omitempty"`
	Address *Address `json:"address,omitempty"`
}

/* Every response variable that is possible must be added
	to the object below to make sending responses easy
*/
type Response struct{
	Message 	string 	`json:"message,omitempty"`
	HueColor 	int 	`json:"hue,omitempty"`
}

func jSONResponse(resp Response) string {
	j, err := json.Marshal(&resp)
	if err != nil {
		fmt.Println("Error Json Marshaling")
	}
	return string(j)
}

type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

type UserSettings struct{
	BeaconColorMap map[int]int `json:"color,omitempty"`
}

var people []Person
var userSettings UserSettings

var CONTENT_TYPE = "Content-Type"
var APP_JSON = "application/json"

type ColorMatch struct{
	Major int `json:"major,omitempty"`
	Color int `json:"hue,omitempty"`
}


func SetColor(w http.ResponseWriter, r *http.Request){
	//params := mux.Vars(r)
	w.Header().Set(CONTENT_TYPE,APP_JSON)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
	}
	var u ColorMatch
	err = json.Unmarshal(body, &u)
	if(userSettings.BeaconColorMap == nil){
		userSettings = UserSettings{map[int]int{u.Major:u.Color}}
		bs,_ := json.Marshal(userSettings)
		fmt.Fprintf(w,string(bs))
		return
	}else{
		userSettings.BeaconColorMap[u.Major] = u.Color
		bs,_ := json.Marshal(userSettings)
		fmt.Fprintf(w,string(bs))
		return
	}

}

func GetAllColors(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set(CONTENT_TYPE,APP_JSON)
	if userSettings.BeaconColorMap == nil{
		fmt.Fprintf(w,jSONResponse(Response{Message:"No beacons registered with color"}))
		return
	}
	data,err := json.Marshal(userSettings)
	if err != nil{
		fmt.Fprintf(w,jSONResponse(Response{Message:"Error in get all colors"}))
		return
	}
	fmt.Fprintf(w,string(data))
	return
}

func GetColor(w http.ResponseWriter, r *http.Request)  {
	params := mux.Vars(r)
	w.Header().Set(CONTENT_TYPE,APP_JSON)
	number,_ := strconv.ParseInt(params["id"],10,0)
	val,prs := userSettings.BeaconColorMap[int(number)]
	if(prs){
		res := Response{
			HueColor : val,
		}
		fmt.Fprintf(w,jSONResponse(res))
		return
	}else {
		res := Response{
			Message : "No color associated to the beacon",
		}
		fmt.Fprintf(w,jSONResponse(res))
		return
	}

}

func GetPerson(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	w.Header().Set("Content-Type", "application/json")
	for _, item := range people {
		if item.ID == params["id"] {
			bs,_ := json.Marshal(item)
			fmt.Fprintf(w,string(bs))
			return
		}
	}
	json.NewEncoder(w).Encode(&Person{})
}

func GetPeople(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(people)
}

func CreatePerson(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var person Person
	_ = json.NewDecoder(req.Body).Decode(&person)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode("Person with ID already exists")
			return
		}
	}
	person.ID = params["id"]
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
}

func DeletePerson(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
		}
	}
	json.NewEncoder(w).Encode(people)
}

func main() {
	router := mux.NewRouter()
	people = append(people, Person{ID: "1", FIRST: "Koushik", LAST: "KASHOJJULA", Address: &Address{City: "Charlotte", State: "NC"}})
	people = append(people, Person{ID: "2", FIRST: "Kittu", LAST: "K"})
	router.HandleFunc("/people", GetPeople).Methods("GET")
	router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
	router.HandleFunc("/people/{id}", CreatePerson).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")
	router.HandleFunc("/color",SetColor).Methods("POST")
	router.HandleFunc("/color/{id}",GetColor).Methods("GET")
	router.HandleFunc("/colors",GetAllColors).Methods("GET")
	log.Fatal(http.ListenAndServe(":3000", router))
}
