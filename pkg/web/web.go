package web

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"partielgo/pkg/db"
	"strconv"
)

func Serve(database *db.Database) {
	// Serve static files (like index.html, CSS, JS)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Handle requests to the root or home page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "static/index.html")
	})

	// Handle requests for room data
	http.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			rooms, err := database.ListRooms()
			if err != nil {
				http.Error(w, "Failed to fetch rooms", http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(rooms)
		}
	})

	// Handle adding a new room
	http.HandleFunc("/addRoom", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var room db.Room
			err := json.NewDecoder(r.Body).Decode(&room)
			if err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}
			if err = database.AddRoom(room.Name, room.Capacity); err != nil {
				http.Error(w, "Failed to add room", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
		}
	})

	// Handle removing a room
	http.HandleFunc("/deleteRoom", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			roomID, err := strconv.Atoi(r.URL.Query().Get("id"))
			if err != nil {
				http.Error(w, "Invalid room ID", http.StatusBadRequest)
				return
			}
			if err = database.RemoveRoom(roomID); err != nil {
				http.Error(w, "Failed to delete room", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		}
	})

	// Handle requests for reservation data
	http.HandleFunc("/reservations", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			reservations, err := database.ListReservations()
			if err != nil {
				http.Error(w, "Failed to fetch reservations", http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(reservations)
		}
	})

	// Handle adding a new reservation
	http.HandleFunc("/addReservation", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var reservation db.Reservation
			err := json.NewDecoder(r.Body).Decode(&reservation)
			if err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}
			if err = database.AddReservation(reservation.RoomID, reservation.Date, reservation.StartTime, reservation.EndTime); err != nil {
				http.Error(w, "Failed to add reservation", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
		}
	})

	// Handle removing a reservation
	http.HandleFunc("/deleteReservation", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			reservationID, err := strconv.Atoi(r.URL.Query().Get("id"))
			if err != nil {
				http.Error(w, "Invalid reservation ID", http.StatusBadRequest)
				return
			}
			if err = database.CancelReservation(reservationID); err != nil {
				http.Error(w, "Failed to delete reservation", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		}
	})

	// Handle exporting reservations to JSON or CSV
	http.HandleFunc("/exportReservations", func(w http.ResponseWriter, r *http.Request) {
		format := r.URL.Query().Get("format")
		switch format {
		case "json":
			reservations, err := database.ListReservations()
			if err != nil {
				http.Error(w, "Failed to fetch reservations", http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(reservations)
		case "csv":
			reservations, err := database.ListReservations()
			if err != nil {
				http.Error(w, "Failed to fetch reservations", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/csv")
			w.Header().Set("Content-Disposition", "attachment;filename=reservations.csv")
			writer := csv.NewWriter(w)
			defer writer.Flush()
			writer.Write([]string{"ID", "Room ID", "Date", "Start Time", "End Time"})
			for _, res := range reservations {
				writer.Write([]string{strconv.Itoa(res.ID), strconv.Itoa(res.RoomID), res.Date, res.StartTime, res.EndTime})
			}
		default:
			http.Error(w, "Invalid format", http.StatusBadRequest)
		}
	})

	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
