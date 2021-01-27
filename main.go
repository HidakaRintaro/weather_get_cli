package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
)

type Weather struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

type Graph struct {
	Name string
	Temp float64
	Max  float64
	Min  float64
}

func (data *Weather) ins() {
	db, err := sql.Open("mysql", "root:@(127.0.0.1)/weather_go")
	if err != nil {
		log.Fatalln("db error")
	}
	defer db.Close()

	ins, err := db.Prepare("INSERT INTO weather (name, temp, max, min) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatalln(err)
	}
	ins.Exec(data.Name, data.Main.Temp, data.Main.TempMax, data.Main.TempMin)
}

func sec() {
	db, err := sql.Open("mysql", "root:@(127.0.0.1)/weather_go")
	if err != nil {
		log.Fatalln("db error")
	}
	defer db.Close()

	rows, err := db.Query("SELECT name, temp, max, min FROM weather ORDER BY id")
	if err != nil {
		log.Fatalln(err)
	}
	var glist []Graph
	for rows.Next() {
		g := Graph{}
		if err := rows.Scan(&g.Name, &g.Temp, &g.Max, &g.Min); err != nil {
			log.Fatalln(err)
		}
		glist = append(glist, g)
	}
	for _, v := range glist {
		fmt.Println("name:", v.Name, ", temp:", v.Temp, ", max:", v.Max, ", min:", v.Min)
	}
}

func main() {
	city := "Osaka,JP"
	key := "****************"
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%v&appid=%v&lang=ja&units=metric", city, key)
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	byteArray, _ := ioutil.ReadAll(resp.Body)

	jsonBytes := ([]byte)(byteArray)
	var data Weather
	err := json.Unmarshal(jsonBytes, &data)
	if err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return
	}
	fmt.Println(data)

	fmt.Printf("都市名 %v, 平均気温 %v, 最高気温 %v, 最低気温 %v \n", data.Name, data.Main.Temp, data.Main.TempMax, data.Main.TempMin)
	data.ins()

	sec()
}
