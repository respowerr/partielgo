# PROJET SOUTENANCE GO

Projet de réservation de salle en GoLang en CLI et interface WEB.

## Bibliothèque

- Go-sqlite3

## Fonctionnalités

- Lister toutes les salles.
- Ajouter une salle.
- Supprimer une salle.
- Créer une réservation.
- Annuler une réservation.
- Lister toutes les réservations.
- Exporter les réservations sous format JSON & CSV.
- Lister les salles disponibles.
- Visualisation des réservations pour une salle et une date.

## Routes pour l'interface WEB

#### Lister toutes les chambres

```http
  GET /rooms
```

#### Ajouter une chambre

```http
  POST /addroom
```

#### Lister toutes les réservations

```http
  GET /reservations
```

#### Créer une réservation

```http
  POST /createreservation
```

#### Supprimer une réservation

```http
  DELETE /deletereservation
```
| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `id`      | `int`    | Id of the reservation      |

#### Supprimer une chambre

```http
  DELETE /deleteroom
```
| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `id`      | `int`    | Id of the room             |

## Installation

```bash
  cd partielgo
  go run main.go
```

## Auteurs

PIZZUTI Enzo, NOIRET Jules, MESLIN Adrien