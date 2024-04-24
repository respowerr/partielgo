package cli

import (
	"bufio"
	"fmt"
	"os"
	"partielgo/pkg/db"
	"strconv"
)

func Run(database *db.Database) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\nMenu principal :")
		fmt.Println("1. Lister toutes les salles")
		fmt.Println("2. Ajouter une salle")
		fmt.Println("3. Supprimer une salle")
		fmt.Println("4. Créer une réservation")
		fmt.Println("5. Annuler une réservation")
		fmt.Println("6. Lister toutes les réservations")
		fmt.Println("7. Exporter les réservations en JSON")
		fmt.Println("8. Exporter les réservations en CSV")
		fmt.Println("9. Quitter")

		fmt.Print("Choisissez une option : ")
		if !scanner.Scan() {
			continue
		}
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
			fmt.Println("Au revoir !")
			return
		default:
			fmt.Println("Option non valide. Veuillez choisir une option valide.")
		}
	}
}

func listRooms(database *db.Database) {
	rooms := database.ListRooms()
	if len(rooms) == 0 {
		fmt.Println("Aucune salle disponible.")
	} else {
		for _, room := range rooms {
			fmt.Printf("Salle %d: %s (Capacité: %d)\n", room.ID, room.Name, room.Capacity)
		}
	}
}

func addRoom(database *db.Database, scanner *bufio.Scanner) {
	fmt.Print("Entrez le nom de la salle : ")
	scanner.Scan()
	name := scanner.Text()
	fmt.Print("Entrez la capacité de la salle : ")
	scanner.Scan()
	capacity, _ := strconv.Atoi(scanner.Text())

	if err := database.AddRoom(name, capacity); err != nil {
		fmt.Println("Erreur lors de l'ajout de la salle :", err)
	} else {
		fmt.Println("Salle ajoutée avec succès.")
	}
}

func removeRoom(database *db.Database, scanner *bufio.Scanner) {
	fmt.Print("Entrez l'ID de la salle à supprimer : ")
	scanner.Scan()
	roomID, _ := strconv.Atoi(scanner.Text())

	if err := database.RemoveRoom(roomID); err != nil {
		fmt.Println("Erreur lors de la suppression de la salle :", err)
	} else {
		fmt.Println("Salle supprimée avec succès.")
	}
}

func createReservation(database *db.Database, scanner *bufio.Scanner) {
	fmt.Print("ID de la salle : ")
	scanner.Scan()
	roomID, _ := strconv.Atoi(scanner.Text())
	fmt.Print("Date (AAAA-MM-JJ) : ")
	scanner.Scan()
	date := scanner.Text()
	fmt.Print("Heure de début (HH:MM) : ")
	scanner.Scan()
	startTime := scanner.Text()
	fmt.Print("Heure de fin (HH:MM) : ")
	scanner.Scan()
	endTime := scanner.Text()

	reservation := db.Reservation{
		RoomID:    roomID,
		Date:      date,
		StartTime: startTime,
		EndTime:   endTime,
	}
	if err := database.AddReservation(reservation); err != nil {
		fmt.Println("Erreur lors de la création de la réservation :", err)
	} else {
		fmt.Println("Réservation créée avec succès.")
	}
}

func cancelReservation(database *db.Database, scanner *bufio.Scanner) {
	fmt.Print("ID de la réservation à annuler : ")
	scanner.Scan()
	reservationID, _ := strconv.Atoi(scanner.Text())

	if err := database.CancelReservation(reservationID); err != nil {
		fmt.Println("Erreur lors de l'annulation de la réservation :", err)
	} else {
		fmt.Println("Réservation annulée avec succès.")
	}
}

func listReservations(database *db.Database) {
	reservations := database.ListAllReservations()
	if len(reservations) == 0 {
		fmt.Println("Aucune réservation trouvée.")
	} else {
		for _, res := range reservations {
			fmt.Printf("Réservation ID %d pour la salle %d du %s de %s à %s\n", res.ID, res.RoomID, res.Date, res.StartTime, res.EndTime)
		}
	}
}

func exportReservations(database *db.Database, format string) {
	filename := ""
	if format == "json" {
		filename = "reservations.json"
	} else if format == "csv" {
		filename = "reservations.csv"
	}

	if err := database.ExportReservations(filename, format); err != nil {
		fmt.Println("Erreur lors de l'export des réservations :", err)
	} else {
		fmt.Println("Réservations exportées avec succès.")
	}
}
