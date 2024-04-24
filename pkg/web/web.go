package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"partielgo/pkg/db"
)

func Serve(database *db.Database) {
	// Serveur de fichiers statiques pour les assets comme le CSS et le JavaScript
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Gestionnaire de la racine qui sert le fichier HTML principal
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("static/index.html"))
		tmpl.Execute(w, nil)
	})

	// Gestionnaire pour lister toutes les salles
	http.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		rooms := database.ListRooms()
		json.NewEncoder(w).Encode(rooms)
	})

	// Gestionnaire pour lister toutes les réservations
	http.HandleFunc("/reservations", func(w http.ResponseWriter, r *http.Request) {
		reservations := database.ListAllReservations()
		json.NewEncoder(w).Encode(reservations)
	})

	// Gestionnaire pour créer une réservation
	http.HandleFunc("/createReservation", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Méthode non supportée", http.StatusMethodNotAllowed)
			return
		}
		var res db.Reservation
		if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := database.AddReservation(res); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(res)
		}
	})

	// Gestionnaire pour annuler une réservation
	http.HandleFunc("/cancelReservation", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Méthode non supportée", http.StatusMethodNotAllowed)
			return
		}
		var data struct{ ID int }
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := database.CancelReservation(data.ID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Réservation annulée avec succès")
		}
	})

	// Gestionnaire pour ajouter une salle
	http.HandleFunc("/addRoom", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Méthode non supportée", http.StatusMethodNotAllowed)
			return
		}
		var room db.Room
		if err := json.NewDecoder(r.Body).Decode(&room); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := database.AddRoom(room.Name, room.Capacity); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(room)
		}
	})

	// Gestionnaire pour supprimer une salle
	http.HandleFunc("/removeRoom", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Méthode non supportée", http.StatusMethodNotAllowed)
			return
		}
		var data struct{ ID int }
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := database.RemoveRoom(data.ID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Salle supprimée avec succès")
		}
	})

	log.Println("Serveur démarré sur le port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
