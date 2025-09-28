package main

import (
	"log"
	"net/http"
)

type WebServer struct {
	api    *WildlifeNLAPI
	server *http.Server
}

func NewWebServer(addr string, api *WildlifeNLAPI) *WebServer {
	mux := http.NewServeMux()
	webServer := &WebServer{api: api, server: &http.Server{Addr: addr, Handler: mux}}
	mux.HandleFunc("/", webServer.rootHandler)
	mux.HandleFunc("/animals", webServer.animalsHandler)
	return webServer
}

func (s *WebServer) rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`
		<!DOCTYPE html>
		<html>
		<head>
		<meta charset="UTF-8">
		<title>WildlifeNL Animals</title>
		<link
			rel="stylesheet"
			href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css"
			integrity="sha256-p4NxAoJBhIIN+hmNHrzRCf9tD/miZyoHS5obTRR9BMY="
			crossorigin="" />
		<style>
			#map { height: 100vh; width: 100%; }
		</style>
		</head>
		<body>
		<div id="map"></div>
		<script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js" integrity="sha256-20nQCchB9co0qIjJZRGuk2/Z9VM+kNiyxNV1lvTlZBo=" crossorigin=""></script>
		<script>
		const map = L.map('map').setView([` + centralPoint.Location() + `], 18);
		L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors', maxZoom: 19}).addTo(map);
		fetch('/animals')
		.then(response => {
			return response.json();
		})
		.then(data => {
			data.forEach(item => {
				const lat = item.location?.latitude;
				const lng = item.location?.longitude;
				const marker = L.marker([lat, lng]).addTo(map);
				marker.bindPopup("<b>" + item.name + "</b><br/>" + item.species?.commonName + " (" + item.species.name + ")<br/>" + item.species.category + "<br/>");
			});
		});
		</script>
		</body>
		</html>
	`))
}

func (s *WebServer) animalsHandler(w http.ResponseWriter, r *http.Request) {
	animals, err := s.api.GetAnimals()
	if err != nil {
		log.Println(err)
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(animals)
}

func (s *WebServer) ListenAndServe() error {
	return s.server.ListenAndServe()
}
