package cli

import (
	"bufio"
	"fmt"
	"os"
	"partielgo/pkg/db"
	"strconv"
	"strings"
	"time"
)

func Run(database *db.Database) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\nMenu principal:")
		fmt.Println("1. Lister toutes les salles")
		fmt.Println("2. Ajouter une salle")
		fmt.Println("3. Supprimer une salle")
		fmt.Println("4. Créer une réservation")
		fmt.Println("5. Annuler une réservation")
		fmt.Println("6. Lister toutes les réservations")
		fmt.Println("7. Exporter les réservations en JSON")
		fmt.Println("8. Exporter les réservations en CSV")
		fmt.Println("9. Lister les salles disponibles")
		fmt.Println("10. Visualiser les réservations pour une salle et une date")
		fmt.Println("11. Quitter")

		fmt.Print("Veuillez choisir une option: ")
		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			listRooms(database)
		case "2":
			addRoom(database, scanner)
		case "3":
			removeRoom(database, scanner)
		case "4":
			createReservation(database, scanner)
		case "5":
			cancelReservation(database, scanner)
		case "6":
			listReservations(database)
		case "7":
			exportReservations(database, "json")
		case "8":
			exportReservations(database, "csv")
		case "9":
			listAvailableRooms(database, scanner)
		case "10":
			viewReservationsForRoomAndDate(database, scanner)
		case "11":
			fmt.Println("Au revoir !")
			os.Exit(0)
		default:
			fmt.Println("Option non valide.")
		}
	}
}

func listAvailableRooms(database *db.Database, scanner *bufio.Scanner) {
	fmt.Print("Entrez la date (YYYY-MM-DD): ")
	scanner.Scan()
	date := scanner.Text()

	fmt.Print("Entrez l'heure de début (HH:MM): ")
	scanner.Scan()
	startTime := scanner.Text()

	fmt.Print("Entrez l'heure de fin (HH:MM): ")
	scanner.Scan()
	endTime := scanner.Text()

	availableRooms, err := database.ListAvailableRooms(date, startTime, endTime)
	if err != nil {
		fmt.Printf("Erreur lors de la récupération des salles disponibles: %v\n", err)
		return
	}
	fmt.Println("\nSalles disponibles:")
	for _, room := range availableRooms {
		fmt.Printf("ID: %d, Nom: %s, Capacité: %d\n", room.ID, room.Name, room.Capacity)
	}
}

func listRooms(database *db.Database) {
	rooms, err := database.ListRooms()
	if err != nil {
		fmt.Printf("Erreur lors de la liste des salles: %v\n", err)
		return
	}
	fmt.Println("\nListe des salles:")
	for _, room := range rooms {
		fmt.Printf("ID: %d, Nom: %s, Capacité: %d\n", room.ID, room.Name, room.Capacity)
	}
}

func addRoom(database *db.Database, scanner *bufio.Scanner) {
	fmt.Print("Entrez le nom de la salle: ")
	scanner.Scan()
	name := scanner.Text()

	fmt.Print("Entrez la capacité de la salle: ")
	scanner.Scan()
	capacity, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Erreur: La capacité doit être un nombre.")
		return
	}

	if err := database.AddRoom(name, capacity); err != nil {
		fmt.Printf("Erreur lors de l'ajout de la salle: %v\n", err)
	} else {
		fmt.Println("Salle ajoutée avec succès.")
	}
}

func removeRoom(database *db.Database, scanner *bufio.Scanner) {
	fmt.Print("Entrez l'ID de la salle à supprimer: ")
	scanner.Scan()
	roomID, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Erreur: L'ID de la salle doit être un nombre.")
		return
	}

	if err := database.RemoveRoom(roomID); err != nil {
		fmt.Printf("Erreur lors de la suppression de la salle: %v\n", err)
	} else {
		fmt.Println("Salle supprimée avec succès.")
	}
}

func createReservation(database *db.Database, scanner *bufio.Scanner) {
	roomID := getValidIntInput(scanner, "Entrez l'ID de la salle: ")
	date := getValidDateInput(scanner, "Entrez la date de la réservation (YYYY-MM-DD): ")
	startTime := getValidTimeInput(scanner, "Entrez l'heure de début (HH:MM): ")
	endTime := getValidTimeInput(scanner, "Entrez l'heure de fin (HH:MM): ")

	if err := database.AddReservation(roomID, date, startTime, endTime); err != nil {
		fmt.Printf("Erreur lors de la création de la réservation: %v\n", err)
	} else {
		fmt.Println("Réservation créée avec succès.")
	}
}

func getValidIntInput(scanner *bufio.Scanner, prompt string) int {
	var input int
	var err error
	for {
		fmt.Print(prompt)
		scanner.Scan()
		input, err = strconv.Atoi(strings.TrimSpace(scanner.Text()))
		if err == nil {
			break
		}
		fmt.Println("Entrée non valide. Veuillez entrer un nombre entier.")
	}
	return input
}

func getValidDateInput(scanner *bufio.Scanner, prompt string) string {
	var date time.Time
	var err error
	for {
		fmt.Print(prompt)
		scanner.Scan()
		date, err = time.Parse("2006-01-02", strings.TrimSpace(scanner.Text()))
		if err == nil {
			break
		}
		fmt.Println("Entrée non valide. Veuillez entrer une date au format YYYY-MM-DD.")
	}
	return date.Format("2006-01-02")
}

func getValidTimeInput(scanner *bufio.Scanner, prompt string) string {
	var t time.Time
	var err error
	for {
		fmt.Print(prompt)
		scanner.Scan()
		t, err = time.Parse("15:04", strings.TrimSpace(scanner.Text()))
		if err == nil {
			break
		}
		fmt.Println("Entrée non valide. Veuillez entrer une heure au format HH:MM.")
	}
	return t.Format("15:04")
}

func cancelReservation(database *db.Database, scanner *bufio.Scanner) {
	fmt.Print("Entrez l'ID de la réservation à annuler: ")
	scanner.Scan()
	reservationID, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Erreur: L'ID de la réservation doit être un nombre.")
		return
	}

	if err := database.CancelReservation(reservationID); err != nil {
		fmt.Printf("Erreur lors de l'annulation de la réservation: %v\n", err)
	} else {
		fmt.Println("Réservation annulée avec succès.")
	}
}

func listReservations(database *db.Database) {
	reservations, err := database.ListReservations()
	if err != nil {
		fmt.Printf("Erreur lors de la liste des réservations: %v\n", err)
		return
	}
	fmt.Println("\nListe des réservations:")
	for _, r := range reservations {
		fmt.Printf("ID: %d, Salle ID: %d, Date: %s, Heure de début: %s, Heure de fin: %s\n", r.ID, r.RoomID, r.Date, r.StartTime, r.EndTime)
	}
}

func exportReservations(database *db.Database, format string) {
	if err := database.ExportReservations(format); err != nil {
		fmt.Printf("Erreur lors de l'exportation des réservations: %v\n", err)
	} else {
		fmt.Printf("Réservations exportées au format %s avec succès.\n", format)
	}
}
func viewReservationsForRoomAndDate(database *db.Database, scanner *bufio.Scanner) {
	fmt.Print("Entrez l'ID de la salle: ")
	scanner.Scan()
	roomID, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Entrée non valide pour l'ID de la salle.")
		return
	}

	fmt.Print("Entrez la date (YYYY-MM-DD): ")
	scanner.Scan()
	date := scanner.Text()

	reservations, err := database.ViewReservations(roomID, date)
	if err != nil {
		fmt.Printf("Erreur lors de l'affichage des réservations: %v\n", err)
		return
	}
	fmt.Println("\nRéservations pour la salle et la date choisies:")
	for _, reservation := range reservations {
		fmt.Printf("ID: %d, Room ID: %d, Date: %s, Start Time: %s, End Time: %s\n", reservation.ID, reservation.RoomID, reservation.Date, reservation.StartTime, reservation.EndTime)
	}
}
