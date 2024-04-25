package db

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type Room struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
}

type Reservation struct {
	ID        int    `json:"id"`
	RoomID    int    `json:"roomID"`
	Date      string `json:"date"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

type Database struct {
	Connection *sql.DB
}

func New() *Database {
	db, err := sql.Open("sqlite3", "./reservation.db")
	if err != nil {
		log.Fatal(err)
	}
	if err := createTables(db); err != nil {
		log.Fatal(err)
	}
	return &Database{Connection: db}
}

func createTables(db *sql.DB) error {
	_, err := db.Exec(`
    CREATE TABLE IF NOT EXISTS rooms (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT,
        capacity INTEGER
    );
    CREATE TABLE IF NOT EXISTS reservations (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        room_id INTEGER,
        date TEXT,
        start_time TEXT,
        end_time TEXT,
        FOREIGN KEY(room_id) REFERENCES rooms(id)
    );`)
	return err
}

func (db *Database) AddRoom(name string, capacity int) error {
	stmt, err := db.Connection.Prepare("INSERT INTO rooms (name, capacity) VALUES (?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(name, capacity)
	return err
}

func (db *Database) ListRooms() ([]Room, error) {
	rows, err := db.Connection.Query("SELECT id, name, capacity FROM rooms")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []Room
	for rows.Next() {
		var r Room
		if err := rows.Scan(&r.ID, &r.Name, &r.Capacity); err != nil {
			return nil, err
		}
		rooms = append(rooms, r)
	}
	return rooms, nil
}

func (db *Database) ListAvailableRooms(date, startTime, endTime string) ([]Room, error) {
	var availableRooms []Room
	query := `SELECT * FROM rooms WHERE id NOT IN (
		SELECT room_id FROM reservations WHERE date = ? AND NOT (
			(end_time <= ? OR start_time >= ?)
		)
	);`
	rows, err := db.Connection.Query(query, date, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var room Room
		if err = rows.Scan(&room.ID, &room.Name, &room.Capacity); err != nil {
			return nil, err
		}
		availableRooms = append(availableRooms, room)
	}
	return availableRooms, nil
}

func (db *Database) RemoveRoom(id int) error {
	_, err := db.Connection.Exec("DELETE FROM rooms WHERE id = ?", id)
	return err
}

func (db *Database) AddReservation(roomID int, date, startTime, endTime string) error {
	availableRooms, err := db.ListAvailableRooms(date, startTime, endTime)
	if err != nil {
		return err
	}
	isAvailable := false
	for _, room := range availableRooms {
		if room.ID == roomID {
			isAvailable = true
			break
		}
	}
	if !isAvailable {
		return fmt.Errorf("the room is not available at the selected time")
	}

	stmt, err := db.Connection.Prepare("INSERT INTO reservations (room_id, date, start_time, end_time) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(roomID, date, startTime, endTime)
	return err
}

func (db *Database) ListReservations() ([]Reservation, error) {
	rows, err := db.Connection.Query("SELECT id, room_id, date, start_time, end_time FROM reservations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []Reservation
	for rows.Next() {
		var r Reservation
		if err := rows.Scan(&r.ID, &r.RoomID, &r.Date, &r.StartTime, &r.EndTime); err != nil {
			return nil, err
		}
		reservations = append(reservations, r)
	}
	return reservations, nil
}

func (db *Database) ViewReservations(roomID int, date string) ([]Reservation, error) {
	rows, err := db.Connection.Query("SELECT id, room_id, date, start_time, end_time FROM reservations WHERE room_id = ? AND date = ?", roomID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []Reservation
	for rows.Next() {
		var reservation Reservation
		if err = rows.Scan(&reservation.ID, &reservation.RoomID, &reservation.Date, &reservation.StartTime, &reservation.EndTime); err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}
	return reservations, nil
}

func (db *Database) CancelReservation(id int) error {
	res, err := db.Connection.Exec("DELETE FROM reservations WHERE id = ?", id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("no reservation found with id %d", id)
	}
	return nil
}

func (db *Database) ExportReservations(format string) error {
	switch format {
	case "json":
		return db.exportReservationsAsJSON()
	case "csv":
		return db.exportReservationsAsCSV()
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func (db *Database) exportReservationsAsJSON() error {
	reservations, err := db.ListReservations()
	if err != nil {
		return err
	}
	file, err := os.Create("reservations.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(reservations)
	if err != nil {
		return fmt.Errorf("error encoding reservations to JSON: %v", err)
	}

	return nil
}

func (db *Database) exportReservationsAsCSV() error {
	reservations, err := db.ListReservations()
	if err != nil {
		return err
	}
	file, err := os.Create("reservations.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"ID", "RoomID", "Date", "StartTime", "EndTime"}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("error writing CSV header: %v", err)
	}

	for _, res := range reservations {
		record := []string{
			strconv.Itoa(res.ID),
			strconv.Itoa(res.RoomID),
			res.Date,
			res.StartTime,
			res.EndTime,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing reservation to CSV: %v", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("error flushing CSV writer: %v", err)
	}

	return nil
}
