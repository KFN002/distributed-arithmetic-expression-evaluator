package models

import (
	"sync"
	"time"
)

var ServersQuantity = 4

var Servers = ServersData{Servers: InitServers()}

// Server Данные сервера
type Server struct {
	ID       int
	Status   string
	Tasks    string
	LastPing string
}

// ServersData данные о серверах
type ServersData struct {
	Mu      sync.Mutex
	Servers map[int]*Server
}

func UpdateServers(id int, operation string, status string) {
	server := Server{ID: id, Status: status, Tasks: operation, LastPing: time.Now().Format("02-01-2006 15:04:05")}
	Servers.Mu.Lock()
	defer Servers.Mu.Unlock()
	Servers.Servers[id] = &server
}

func InitServers() map[int]*Server {
	servers := make(map[int]*Server)
	for serverID := 1; serverID <= ServersQuantity; serverID++ {
		server := &Server{ID: serverID, Status: "Stand By", Tasks: "", LastPing: time.Now().Format("02-01-2006 15:04:05")}
		servers[serverID] = server
	}
	return servers
}
