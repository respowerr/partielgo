document.addEventListener('DOMContentLoaded', function() {
    fetchRooms();
    fetchReservations();

    document.getElementById('addRoomForm').addEventListener('submit', function(e) {
        e.preventDefault();
        const roomName = document.getElementById('roomName').value;
        const roomCapacity = document.getElementById('roomCapacity').value;
        addRoom(roomName, roomCapacity);
    });

    document.getElementById('addReservationForm').addEventListener('submit', function(e) {
        e.preventDefault();
        const reservationRoomId = document.getElementById('reservationRoomId').value;
        const reservationDate = document.getElementById('reservationDate').value;
        const reservationStartTime = document.getElementById('reservationStartTime').value;
        const reservationEndTime = document.getElementById('reservationEndTime').value;
        addReservation(reservationRoomId, reservationDate, reservationStartTime, reservationEndTime);
    });
});

function fetchRooms() {
    fetch('/rooms')
    .then(response => response.json())
    .then(rooms => {
        const roomList = document.getElementById('roomList');
        roomList.innerHTML = '<tr><th>ID</th><th>Nom</th><th>Capacité</th></tr>';
        rooms.forEach(room => {
            roomList.innerHTML += `<tr><td>${room.ID}</td><td>${room.Name}</td><td>${room.Capacity}</td></tr>`;
        });
    })
    .catch(error => console.error('Erreur lors de la récupération des salles:', error));
}

function fetchReservations() {
    fetch('/reservations')
    .then(response => response.json())
    .then(reservations => {
        const reservationList = document.getElementById('reservationList');
        reservationList.innerHTML = '<tr><th>ID</th><th>ID Salle</th><th>Date</th><th>Heure de début</th><th>Heure de fin</th></tr>';
        reservations.forEach(res => {
            reservationList.innerHTML += `<tr><td>${res.ID}</td><td>${res.RoomID}</td><td>${res.Date}</td><td>${res.StartTime}</td><td>${res.EndTime}</td></tr>`;
        });
    })
    .catch(error => console.error('Erreur lors de la récupération des réservations:', error));
}

function addRoom(name, capacity) {
    fetch('/addRoom', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name: name, capacity: parseInt(capacity, 10) })
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();
    })
    .then(data => {
        console.log('Salle ajoutée:', data);
        fetchRooms(); 
    })
    .catch(error => {
        console.error('Erreur lors de l\'ajout de la salle:', error);
    });
    
}


function addReservation(roomId, date, startTime, endTime) {
    fetch('/createReservation', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ roomID: parseInt(roomId, 10), date: date, startTime: startTime, endTime: endTime })
    })
    .then(response => response.json())
    .then(data => {
        console.log('Réservation créée:', data);
        fetchReservations(); 
    })
    .catch(error => console.error('Erreur lors de la création de la réservation:', error));
}
