package test

import (
	"encoding/json"
	"fmt"
)

type User struct {
	ID          int               `json:"id"`
	Age         int               `json:"age"`
	Name        string            `json:"name"`
	Password    *string           `json:"password"`
	Coordinate  Coordinate        `json:"coordinate"`
	Coordinate2 *Coordinate       `json:"coordinate2"`
	Map         map[string]string `json:"map"`
	Address     []string          `json:"address"`
}

type Coordinate struct {
	Lat float64 `json:"latitude"`
	Lng float64 `json:"longitude"`
}

func JsonTest() {
	user := User{
		ID:         1,
		Password:   nil,
		Coordinate: Coordinate{Lat: 37.7749, Lng: -122.4194},
		Map:        map[string]string{},
		Address:    []string{},
	}

	user2 := User{
		ID:         1,
		Password:   nil,
		Coordinate: Coordinate{Lat: 37.7749, Lng: -122.4194},
	}
	jsonData, _ := json.Marshal(user)
	jsonData2, _ := json.Marshal(user2)
	fmt.Println("JSON Data:", string(jsonData))
	fmt.Println("JSON Data2:", string(jsonData2))
}

func JsonTest3() {
	const jsonData = `{"map":null,"address":null,id":1,"name":"","password":<nil>,"coordinate":null,"coordinate2":null}`
	var user User
	json.Unmarshal([]byte(jsonData), &user)
	fmt.Println(user)
	fmt.Println("address is nil? ", user.Address == nil)
	fmt.Println("map is nil? ", user.Map == nil)

	fmt.Println("-----------------------------")
	const jsonData2 = `{"map":{},"address":[],id":1,"name":"","password":<nil>,"coordinate":null,"coordinate2":null}`
	var user2 User
	json.Unmarshal([]byte(jsonData2), &user2)
	fmt.Println(user2)
	fmt.Println("address is nil? ", user2.Address == nil)
	fmt.Println("map is nil? ", user2.Map == nil)
	fmt.Println("-----------------------------")
	const jsonData3 = `{id":1,"name":"","password":<nil>,"coordinate":null,"coordinate2":null}`
	var user3 User
	json.Unmarshal([]byte(jsonData3), &user3)
	fmt.Println(user3)
	fmt.Println("address is nil? ", user3.Address == nil)
	fmt.Println("map is nil? ", user3.Map == nil)

}

func JsonTest4() {
	var user User

	fmt.Println("user.Name  ", user.Name)
	fmt.Println("user.ID  ", user.ID)
	fmt.Println("user.Map  ", user.Map)
	fmt.Println("user.Map is nil? ", user.Map == nil)
	fmt.Println("user.Map  ", user.Map["key1"])
	user.Map["key1"] = "value1"

}
