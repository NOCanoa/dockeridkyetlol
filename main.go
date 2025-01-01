package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

func main() {
	print("Starting server on port 8080")

	http.HandleFunc("/routs", routs)
	http.HandleFunc("/bus", bus)
	http.HandleFunc("/bus2", bus2)
	http.HandleFunc("/busliveall", busliveall)
	http.HandleFunc("/up", up)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})
	http.HandleFunc("/static/rep2.jpg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/rep2.jpg")
	})
	http.HandleFunc("/static/rings.svg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/rings.svg")
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func routs(w http.ResponseWriter, r *http.Request) {
	print("routs")
	routsAPI := "https://api.moverick.es/topology/sublines?municipality_id=60c879a7d9f552002f4a99da&template=passenger&per=0"
	response, err := fetchJSON(routsAPI)
	if err != nil {
		http.Error(w, "Error fetching routes data", http.StatusInternalServerError)
		return
	}
	responseJSON(w, response)
}

func bus(w http.ResponseWriter, r *http.Request) {
	print("bus")
	print(r.URL.Query().Get("bus-line"))
	sublineId := r.URL.Query().Get("bus-line")
	print(sublineId + "\n")
	busAPI := "https://api.moverick.es/tracking/last_position_buses?template=passenger_route&subline_id=" + sublineId + "&calculated=true"
	print(busAPI + "\n")

	response, err := fetchJSON(busAPI)
	if err != nil {
		http.Error(w, "Error fetching bus data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get data array from response
	data, ok := response.(map[string]interface{})["data"].([]interface{})
	if !ok {
		http.Error(w, "Invalid data format", http.StatusInternalServerError)
		return
	}

	// Print raw data for debugging
	/* fmt.Printf("Raw data: %+v\n", data) */

	funcMap := template.FuncMap{
		"multiply": func(a, b float64) float64 {
			return a * b
		},
	}
	/* css idea style="background: linear-gradient(90deg, #00bf009e {{multiply .distance_over_route 100}}% , white {{multiply .distance_over_route 100}}% ); */
	tmpl := template.New("bus").Funcs(funcMap)
	tmpl, err = tmpl.Parse(`
		<div class="bus-container">
			{{range .}}
				<div class="bus">
					<h2>{{.next_point_name}}</h2>
					<div class="buss">
						<p>{{.bus_plate}}</p>
						<progress class="progress" value="{{multiply .distance_over_route 100}}" max="100">{{multiply .distance_over_route 100}}%</progress>
					</div>
					<p>{{.route_origin}} → {{.route_destination}}</p>
				</div>
			{{end}}
		</div>
	`)
	if err != nil {
		http.Error(w, "Template parse error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func bus2(w http.ResponseWriter, r *http.Request) {
	print("bus")
	print(r.URL.Query().Get("bus-line"))
	sublineId := r.URL.Query().Get("bus-line")
	print(sublineId + "\n")
	busAPI := "https://api.moverick.es/tracking/last_position_buses?template=passenger_route&subline_id=" + sublineId + "&calculated=true"
	print(busAPI + "\n")

	response, err := fetchJSON(busAPI)
	if err != nil {
		http.Error(w, "Error fetching bus data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get data array from response
	data, ok := response.(map[string]interface{})["data"].([]interface{})
	if !ok {
		http.Error(w, "Invalid data format", http.StatusInternalServerError)
		return
	}

	// Print raw data for debugging
	/* fmt.Printf("Raw data: %+v\n", data) */

	funcMap := template.FuncMap{
		"multiply": func(a, b float64) float64 {
			return a * b
		},
	}
	/* css idea style="background: linear-gradient(90deg, #00bf009e {{multiply .distance_over_route 100}}% , white {{multiply .distance_over_route 100}}% ); */
	tmpl := template.New("bus").Funcs(funcMap)
	// Modify the existing template to include map:
	tmpl, err = tmpl.Parse(`
<div class="bus-container">
	{{range .}}
		<div class="bus">
			<h2>{{.next_point_name}}</h2>
			<div class="buss">
				<p>{{.bus_plate}}</p>
				<progress class="progress" value="{{multiply .distance_over_route 100}}" max="100">{{multiply .distance_over_route 100}}%</progress>
			</div>
			<p>{{.route_origin}} → {{.route_destination}}</p>
			<div id="map-{{.bus_plate}}" style="height: 200px; width: 100%; margin-top: 10px;border-radius: 5px;"></div>
			<script>
		 		var map = L.map('map-{{.bus_plate}}').setView([{{.lat}}, {{.lng}}], 13);
				L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
					attribution: '© OpenStreetMap contributors'
				}).addTo(map);
				L.marker([{{.lat}}, {{.lng}}]).addTo(map);
			</script>
		</div>
	{{end}}
</div>
`)
	if err != nil {
		http.Error(w, "Template parse error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func busliveall(w http.ResponseWriter, r *http.Request) {
	print("busliveall")
	busAPI := "https://api.moverick.es/tracking/last_position_buses?template=passenger_route&administrative_area_id=608f9ee02fa54078921d723d&calculated=true&per=0"
	response, err := fetchJSON(busAPI)
	if err != nil {
		http.Error(w, "Error fetching bus data", http.StatusInternalServerError)
		return
	}
	responseJSON(w, response)
}

func up(w http.ResponseWriter, r *http.Request) {
	print("up")
	data := map[string]string{
		"bunny":  "net",
		"status": "up",
	}
	responseJSON(w, data)
}

/* func root(w http.ResponseWriter, r *http.Request) {
	print("root")
	description := "Hey como chegate para aqui?\nIsto e basicamente uma proxy para a api da VamusAlgarve\nvamus-api -> esta proxi -> meu site\nassim poupando recursos para mim e aos servidores deles.\n\nSe tiveres alguma questao mandame um mail em\nhttps://thatcanoa.org/about.html"
	rotas := "/up uptime\n\n/routs rotas\n\n/v1/{route}/{bus} autocarro\n\n/v1/{route} linha autocarros ativos\n\n/bus/all todos autocarros\n\n/map/all map de todas as estacoes"
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusTeapot)
	fmt.Fprintf(w, "%s\n\n%s", description, rotas)
} */

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

// Helper function to safely parse float values
func parseFloat(v interface{}) float64 {
	switch v := v.(type) {
	case float64:
		return v
	case string:
		f, _ := strconv.ParseFloat(v, 64)
		return f
	default:
		return 0
	}
}
