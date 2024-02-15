package models

import (
	"log"
	"sync"
	"time"
)

var Servers = NewServersManager(4)

type Server struct {
	ID       int
	Status   string
	Tasks    string
	LastPing string
}

type ServersData struct {
	Mu      sync.Mutex
	Servers map[int]*Server
}

type ServersManager struct {
	ServersQuantity int
	Servers         ServersData
}

func NewServersManager(quantity int) *ServersManager {
	return &ServersManager{
		ServersQuantity: quantity,
		Servers:         ServersData{Servers: make(map[int]*Server)},
	}
}

func (sm *ServersManager) InitServers() {
	sm.Servers.Mu.Lock()
	defer sm.Servers.Mu.Unlock()
	for serverID := 1; serverID <= sm.ServersQuantity; serverID++ {
		server := &Server{ID: serverID, Status: "Stand By", Tasks: "", LastPing: time.Now().Format("02-01-2006 15:04:05")}
		sm.Servers.Servers[serverID] = server
	}
}

func (sm *ServersManager) UpdateServers(id int, operation string, status string) {
	server := Server{ID: id, Status: status, Tasks: operation, LastPing: time.Now().Format("02-01-2006 15:04:05")}
	sm.Servers.Mu.Lock()
	defer sm.Servers.Mu.Unlock()
	sm.Servers.Servers[id] = &server
}

func (sm *ServersManager) SendHeartbeat(id int) {
	sm.Servers.Mu.Lock()
	defer sm.Servers.Mu.Unlock()

	server, exists := sm.Servers.Servers[id]
	if !exists {
		log.Println("Server has tripped over a cable, pal!")
	}

	server.LastPing = time.Now().Format("02-01-2006 15:04:05")
}
