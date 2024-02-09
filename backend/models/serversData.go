package models

import (
	"sync"
	"time"
)

var Servers = ServersData{Servers: map[int]*Server{
	1: {ID: 1, Status: "Stand By", Tasks: "", LastPing: time.Now().Format("02-01-2006 15:04:05")},
	2: {ID: 2, Status: "Stand By", Tasks: "", LastPing: time.Now().Format("02-01-2006 15:04:05")},
	3: {ID: 3, Status: "Stand By", Tasks: "", LastPing: time.Now().Format("02-01-2006 15:04:05")},
	4: {ID: 4, Status: "Stand By", Tasks: "", LastPing: time.Now().Format("02-01-2006 15:04:05")},
}}

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
