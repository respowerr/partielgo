package db

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Room struct {
	ID       int
	Name     string
	Capacity int
}

type Reservation struct {
	ID        int
	RoomID    int
	Date      string
	StartTime string
	EndTime   string
}

type Database struct {
	Rooms        []Room
	Reservations []Reservation
}

func New() *Database {
	return &Database{
		Rooms:        []Room{},
		Reservations: []Reservation{},
	}
}

func (db *Database) AddRoom(name string, capacity int) error {
	newID := len(db.Rooms) + 1
	db.Rooms = append(db.Rooms, Room{ID: newID, Name: name, Capacity: capacity})
	return nil
}

func (db *Database) RemoveRoom(roomID int) error {
	for i, room := range db.Rooms {
		if room.ID == roomID {
			db.Rooms = append(db.Rooms[:i], db.Rooms[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Room with ID %d not found", roomID)
}

func (db *Database) AddReservation(reservation Reservation) error {
	if !db.IsRoomAvailable(reservation.RoomID, reservation.Date, reservation.StartTime, reservation.EndTime) {
		return fmt.Errorf("room is not available")
	}
	reservation.ID = len(db.Reservations) + 1
	db.Reservations = append(db.Reservations, reservation)
	return nil
}

func (db *Database) CancelReservation(reservationID int) error {
	for i, res := range db.Reservations {
		if res.ID == reservationID {
			db.Reservations = append(db.Reservations[:i], db.Reservations[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("reservation with ID %d not found", reservationID)
}

func (db *Database) ListRooms() []Room {
	return db.Rooms
}

func (db *Database) ListAllReservations() []Reservation {
	return db.Reservations
}

func (db *Database) IsRoomAvailable(roomID int, date, startTime, endTime string) bool {
	for _, res := range db.Reservations {
		if res.RoomID == roomID && res.Date == date && (startTime < res.EndTime && endTime > res.StartTime) {
			return false
		}
	}
	return true
}

func (db *Database) ExportJSON(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(db.Reservations); err != nil {
		return fmt.Errorf("failed to encode reservations: %w", err)
	}

	return nil
}

func (db *Database) ExportCSV(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"ID", "RoomID", "Date", "StartTime", "EndTime"}); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
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
			return fmt.Errorf("failed to write record: %w", err)
		}
	}
	return nil
}
