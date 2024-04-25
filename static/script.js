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
            roomList.innerHTML = '<tr><th>ID</th><th>Nom</th><th>Capacité</th><th>Actions</th></tr>';
            rooms.forEach(room => {
                roomList.innerHTML += `<tr><td>${room.id}</td><td>${room.name}</td><td>${room.capacity}</td><td><button onclick="deleteRoom(${room.id})">Supprimer</button></td></tr>`;
            });
        })
        .catch(error => console.error('Erreur lors de la récupération des salles:', error));
}

function fetchReservations() {
    fetch('/reservations')
        .then(response => response.json())
        .then(reservations => {
            const reservationList = document.getElementById('reservationList');
            reservationList.innerHTML = '<tr><th>ID</th><th>ID Salle</th><th>Date</th><th>Heure de début</th><th>Heure de fin</th><th>Actions</th></tr>';
            reservations.forEach(reservation => {
                reservationList.innerHTML += `<tr><td>${reservation.id}</td><td>${reservation.roomID}</td><td>${reservation.date}</td><td>${reservation.startTime}</td><td>${reservation.endTime}</td><td><button onclick="deleteReservation(${reservation.id})">Supprimer</button></td></tr>`;
            });
        })
        .catch(error => console.error('Erreur lors de la récupération des réservations:', error));
}

function addRoom(name, capacity) {
    fetch('/addRoom', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ Name: name, Capacity: parseInt(capacity, 10) })
    })
    .then(response => {
        if (!response.ok) throw new Error('Network response was not ok');
        fetchRooms();
    })
    .catch(error => console.error('Erreur lors de l\'ajout de la salle:', error));
}

function addReservation(roomId, date, startTime, endTime) {
    fetch('/createReservation', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ RoomID: parseInt(roomId, 10), Date: date, StartTime: startTime, EndTime: endTime })
    })
    .then(response => {
        if (!response.ok) throw new Error('Network response was not ok');
        fetchReservations();
    })
    .catch(error => console.error('Erreur lors de la création de la réservation:', error));
}

function deleteRoom(roomId) {
    fetch(`/deleteRoom?id=${roomId}`, { method: 'DELETE' })
    .then(response => {
        if (!response.ok) throw new Error('Network response was not ok');
        fetchRooms();
    })
    .catch(error => console.error('Erreur lors de la suppression de la salle:', error));
}

function deleteReservation(reservationId) {
    fetch(`/deleteReservation?id=${reservationId}`, { method: 'DELETE' })
    .then(response => {
        if (!response.ok) throw new Error('Network response was not ok');
        fetchReservations();
    })
    .catch(error => console.error('Erreur lors de la suppression de la réservation:', error));
}
