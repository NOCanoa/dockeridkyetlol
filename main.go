package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	print("Starting server on port 8080")

	http.HandleFunc("/up", up)
	http.HandleFunc("/", serveHTTP)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func up(w http.ResponseWriter, r *http.Request) {
	print("up")
	data := map[string]string{
		"bunny":  "net",
		"status": "up",
	}
	responseJSON(w, data)
}

func serveHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch {
	case path == "/":
		description := "Hey como chegate para aqui?\nIsto e basicamente uma proxy para a api da VamusAlgarve\nvamus-api -> esta proxi -> meu site\nassim poupando recursos para mim e aos servidores deles.\n\nSe tiveres alguma questao mandame um mail em\nhttps://thatcanoa.org/about.html"
		rotas := "/up uptime\n\n/routs rotas\n\n/v1/{route}/{bus} autocarro\n\n/v1/{route} linha autocarros ativos\n\n/bus/all todos autocarros\n\n/map/all map de todas as estacoes"
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusTeapot)
		fmt.Fprintf(w, "%s\n\n%s", description, rotas)

	case path == "/up":
		data := map[string]string{
			"bunny":  "net",
			"status": "up",
		}
		responseJSON(w, data)

	case path == "/bus/all":
		busAPI := "https://api.moverick.es/tracking/last_position_buses?template=passenger_route&administrative_area_id=608f9ee02fa54078921d723d&calculated=true&per=0"
		response, err := fetchJSON(busAPI)
		if err != nil {
			http.Error(w, "Error fetching bus data", http.StatusInternalServerError)
			return
		}
		responseJSON(w, response)

	case path == "/routs":
		routsAPI := "https://api.moverick.es/topology/sublines?municipality_id=60c879a7d9f552002f4a99da&template=passenger&per=0"
		response, err := fetchJSON(routsAPI)
		if err != nil {
			http.Error(w, "Error fetching routes data", http.StatusInternalServerError)
			return
		}
		responseJSON(w, response)

	default:
		http.NotFound(w, r)
	}
}

func fetchJSON(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func responseJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, X-Requested-With")
	json.NewEncoder(w).Encode(data)
}
