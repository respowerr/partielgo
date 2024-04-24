package db

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// Room représente une salle avec son identifiant, nom et capacité.
type Room struct {
	ID       int
	Name     string
	Capacity int
}

// Reservation représente une réservation avec son identifiant, l'identifiant de salle, la date et les heures de début et de fin.
type Reservation struct {
	ID        int
	RoomID    int
	Date      string
	StartTime string
	EndTime   string
}

// Database contient les listes des salles et des réservations.
type Database struct {
	Rooms        []Room
	Reservations []Reservation
}

// New initialise une nouvelle base de données avec des listes vides.
func New() *Database {
	return &Database{
		Rooms:        []Room{},
		Reservations: []Reservation{},
	}
}

// AddRoom ajoute une nouvelle salle à la base de données.
func (db *Database) AddRoom(name string, capacity int) error {
	newID := len(db.Rooms) + 1
	db.Rooms = append(db.Rooms, Room{ID: newID, Name: name, Capacity: capacity})
	return nil
}

// RemoveRoom supprime une salle par son identifiant.
func (db *Database) RemoveRoom(roomID int) error {
	for i, room := range db.Rooms {
		if room.ID == roomID {
			db.Rooms = append(db.Rooms[:i], db.Rooms[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Salle avec l'ID %d non trouvée", roomID)
}

// AddReservation ajoute une nouvelle réservation si la salle est disponible.
func (db *Database) AddReservation(reservation Reservation) error {
	if !db.IsRoomAvailable(reservation.RoomID, reservation.Date, reservation.StartTime, reservation.EndTime) {
		return fmt.Errorf("La salle n'est pas disponible")
	}
	reservation.ID = len(db.Reservations) + 1
	db.Reservations = append(db.Reservations, reservation)
	return nil
}

// CancelReservation annule une réservation par son identifiant.
func (db *Database) CancelReservation(reservationID int) error {
	for i, res := range db.Reservations {
		if res.ID == reservationID {
			db.Reservations = append(db.Reservations[:i], db.Reservations[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Réservation avec l'ID %d non trouvée", reservationID)
}

// ListRooms retourne la liste de toutes les salles.
func (db *Database) ListRooms() []Room {
	return db.Rooms
}

// ListAllReservations retourne la liste de toutes les réservations.
func (db *Database) ListAllReservations() []Reservation {
	return db.Reservations
}

// IsRoomAvailable vérifie si une salle est disponible pour une réservation donnée.
func (db *Database) IsRoomAvailable(roomID int, date, startTime, endTime string) bool {
	for _, res := range db.Reservations {
		if res.RoomID == roomID && res.Date == date && (startTime < res.EndTime && endTime > res.StartTime) {
			return false
		}
	}
	return true
}

// ExportReservations exporte les réservations au format JSON ou CSV.
func (db *Database) ExportReservations(filename string, format string) error {
	if format == "json" {
		return db.ExportJSON(filename)
	} else if format == "csv" {
		return db.ExportCSV(filename)
	}
	return fmt.Errorf("Format non supporté")
}

// ExportJSON écrit les réservations dans un fichier JSON.
func (db *Database) ExportJSON(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Échec de la création du fichier : %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(db.Reservations); err != nil {
		return fmt.Errorf("Échec de l'encodage des réservations : %v", err)
	}

	return nil
}

// ExportCSV écrit les réservations dans un fichier CSV.
func (db *Database) ExportCSV(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Échec de la création du fichier : %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"ID", "RoomID", "Date", "StartTime", "EndTime"}); err != nil {
		return fmt.Errorf("Échec de l'écriture de l'en-tête : %v", err)
	}

	for _, reservation := range db.Reservations {
		record := []string{
			strconv.Itoa(reservation.ID),
			strconv.Itoa(reservation.RoomID),
			reservation.Date,
			reservation.StartTime,
			reservation.EndTime,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("Échec de l'écriture de l'enregistrement : %v", err)
		}
	}
	return nil
}
