package cli

import (
	"bufio"
	"fmt"
	"os"
	"partielgo/pkg/db"
	"strconv"
)

var database = db.New()
var scanner = bufio.NewScanner(os.Stdin)

func Run() {
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
			listRooms()
		case "2":
			addRoom()
		case "3":
			removeRoom()
		case "4":
			createReservation()
		case "5":
			cancelReservation()
		case "6":
			listReservations()
		case "7":
			exportReservations("json")
		case "8":
			exportReservations("csv")
		case "9":
			fmt.Println("Au revoir !")
			return
		default:
			fmt.Println("Option non valide. Veuillez choisir une option valide.")
		}
	}
}

func listRooms() {
	rooms := database.ListRooms()
	if len(rooms) == 0 {
		fmt.Println("Aucune salle disponible.")
	} else {
		for _, room := range rooms {
			fmt.Printf("Salle %d: %s (Capacité: %d)\n", room.ID, room.Name, room.Capacity)
		}
	}
}

func addRoom() {
	name := prompt("Entrez le nom de la salle : ")
	capacity := promptForInt("Entrez la capacité de la salle : ")
	if err := database.AddRoom(name, capacity); err != nil {
		fmt.Println("Erreur lors de l'ajout de la salle :", err)
	} else {
		fmt.Println("Salle ajoutée avec succès.")
	}
}

func removeRoom() {
	roomID := promptForInt("Entrez l'ID de la salle à supprimer : ")
	if err := database.RemoveRoom(roomID); err != nil {
		fmt.Println("Erreur lors de la suppression de la salle :", err)
	} else {
		fmt.Println("Salle supprimée avec succès.")
	}
}

func createReservation() {
	fmt.Println("Entrez les détails de la réservation :")
	roomID := promptForInt("ID de la salle : ")
	date := prompt("Date (YYYY-MM-DD) : ")
	startTime := prompt("Heure de début (HH:MM) : ")
	endTime := prompt("Heure de fin (HH:MM) : ")
	err := database.AddReservation(db.Reservation{
		RoomID:    roomID,
		Date:      date,
		StartTime: startTime,
		EndTime:   endTime,
	})
	if err != nil {
		fmt.Println("Erreur lors de la création de la réservation :", err)
	} else {
		fmt.Println("Réservation créée avec succès.")
	}
}

func cancelReservation() {
	reservationID := promptForInt("ID de la réservation à annuler : ")
	err := database.CancelReservation(reservationID)
	if err != nil {
		fmt.Println("Erreur lors de l'annulation de la réservation :", err)
	} else {
		fmt.Println("Réservation annulée avec succès.")
	}
}

func listReservations() {
	fmt.Println("Liste de toutes les réservations :")
	reservations := database.ListAllReservations()
	if len(reservations) == 0 {
		fmt.Println("Aucune réservation trouvée.")
	} else {
		for _, res := range reservations {
			fmt.Printf("Réservation ID %d pour la salle %d du %s de %s à %s\n", res.ID, res.RoomID, res.Date, res.StartTime, res.EndTime)
		}
	}
}

func exportReservations(format string) {
	filename := prompt("Entrez le nom du fichier pour l'export (ex : reservations.json ou reservations.csv) :")
	var err error
	if format == "json" {
		err = database.ExportJSON(filename)
	} else if format == "csv" {
		err = database.ExportCSV(filename)
	}
	if err != nil {
		fmt.Println("Erreur lors de l'export des réservations :", err)
	} else {
		fmt.Println("Réservations exportées avec succès.")
	}
}

func prompt(message string) string {
	fmt.Print(message)
	scanner.Scan()
	return scanner.Text()
}

func promptForInt(message string) int {
	response := prompt(message)
	number, err := strconv.Atoi(response)
	if err != nil {
		fmt.Println("Veuillez entrer un nombre valide.")
		return promptForInt(message) // Récursivité pour demander à nouveau
	}
	return number
}
