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
		fmt.Println("\nMenu principal:")
		fmt.Println("1. Lister toutes les salles")
		fmt.Println("2. Ajouter une salle")
		fmt.Println("3. Supprimer une salle")
		fmt.Println("4. Créer une réservation")
		fmt.Println("5. Annuler une réservation")
		fmt.Println("6. Lister toutes les réservations")
		fmt.Println("7. Exporter les réservations en JSON")
		fmt.Println("8. Exporter les réservations en CSV")
		fmt.Println("9. Quitter")

		fmt.Print("Choisissez une option: ")
		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			database.ListRooms()
		case "2":
			addRoom(database, scanner)
		case "3":
			removeRoom(database, scanner)
		case "4":
			createReservation(database, scanner)
		case "5":
			cancelReservation(database, scanner)
		case "6":
			database.ListReservations()
		case "7":
			database.ExportReservations("json")
		case "8":
			database.ExportReservations("csv")
		case "9":
			fmt.Println("Au revoir !")
			os.Exit(0)
		default:
			fmt.Println("Option non valide. Veuillez choisir une option valide.")
		}
	}
}

func addRoom(database *db.Database, scanner *bufio.Scanner) {
	fmt.Print("Entrez le nom de la salle: ")
	scanner.Scan()
	name := scanner.Text()
	fmt.Print("Entrez la capacité de la salle: ")
	scanner.Scan()
	capacity, _ := strconv.Atoi(scanner.Text())
	if err := database.AddRoom(name, capacity); err != nil {
		fmt.Printf("Erreur lors de l'ajout de la salle: %v\n", err)
	} else {
		fmt.Println("Salle ajoutée avec succès.")
	}
}

func removeRoom(database *db.Database, scanner *bufio.Scanner) {
	fmt.Print("Entrez l'ID de la salle à supprimer: ")
	scanner.Scan()
	roomID, _ := strconv.Atoi(scanner.Text())
	if err := database.RemoveRoom(roomID); err != nil {
		fmt.Printf("Erreur lors de la suppression de la salle: %v\n", err)
	} else {
		fmt.Println("Salle supprimée avec succès.")
	}
}

func createReservation(database *db.Database, scanner *bufio.Scanner) {
	fmt.Print("Entrez l'ID de la salle: ")
	scanner.Scan()
	roomID, _ := strconv.Atoi(scanner.Text())
	fmt.Print("Entrez la date de la réservation (YYYY-MM-DD): ")
	scanner.Scan()
	date := scanner.Text()
	fmt.Print("Entrez l'heure de début (HH:MM): ")
	scanner.Scan()
	startTime := scanner.Text()
	fmt.Print("Entrez l'heure de fin (HH:MM): ")
	scanner.Scan()
	endTime := scanner.Text()

	if err := database.AddReservation(roomID, date, startTime, endTime); err != nil {
		fmt.Printf("Erreur lors de la création de la réservation: %v\n", err)
	} else {
		fmt.Println("Réservation créée avec succès.")
	}
}

func cancelReservation(database *db.Database, scanner *bufio.Scanner) {
	fmt.Print("Entrez l'ID de la réservation à annuler: ")
	scanner.Scan()
	reservationID, _ := strconv.Atoi(scanner.Text())
	if err := database.CancelReservation(reservationID); err != nil {
		fmt.Printf("Erreur lors de l'annulation de la réservation: %v\n", err)
	} else {
		fmt.Println("Réservation annulée avec succès.")
	}
}
