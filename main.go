package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Room struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var rooms = make(map[string]Room)

func main() {
	http.HandleFunc("/rooms", getRoomsHandler)
	http.HandleFunc("/rooms/create", createRoomHandler)
	http.HandleFunc("/rooms/update", updateRoomHandler)
	http.HandleFunc("/rooms/delete", deleteRoomHandler)
	fmt.Println("Server is listening on :9001")
	http.ListenAndServe(":9001", nil)
}
func getRoomsHandler(w http.ResponseWriter, r *http.Request) {
	var roomList []Room
	for _, room := range rooms {
		roomList = append(roomList, room)
	}
	response, err := json.Marshal(roomList)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
func createRoomHandler(w http.ResponseWriter, r *http.Request) {
	var newRoom Room
	err := json.NewDecoder(r.Body).Decode(&newRoom)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}
	newRoom.ID = generateID()
	rooms[newRoom.ID] = newRoom
	response, err := json.Marshal(newRoom)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
func updateRoomHandler(w http.ResponseWriter, r *http.Request) {
	var updatedRoom Room
	err := json.NewDecoder(r.Body).Decode(&updatedRoom)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}
	if _, ok := rooms[updatedRoom.ID]; !ok {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}
	rooms[updatedRoom.ID] = updatedRoom
	response, err := json.Marshal(updatedRoom)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
func deleteRoomHandler(w http.ResponseWriter, r *http.Request) {
	var deleteRequest struct {
		ID string `json:"id"`
	}
	err := json.NewDecoder(r.Body).Decode(&deleteRequest)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}
	if _, ok := rooms[deleteRequest.ID]; !ok {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}
	delete(rooms, deleteRequest.ID)
	w.WriteHeader(http.StatusNoContent)
}
func generateID() string {
	return fmt.Sprintf("%d", len(rooms)+1)
}
