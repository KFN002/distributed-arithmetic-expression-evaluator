package models

import (
	"log"
	"sync"
	"time"
)

var Servers = NewServersManager(1) // Серверов - 1, меняя переменную будет больше и меньше серверов

// Server структура данных сервера (воркера)
type Server struct {
	ID       int
	Status   string
	Tasks    string
	LastPing string
}

// ServersData структура данных серверов (воркеров)
type ServersData struct {
	Mu      sync.Mutex
	Servers map[int]*Server
}

// ServersManager структура менеджера серверов (воркеров)
type ServersManager struct {
	ServersQuantity int
	Servers         ServersData
}

// NewServersManager создание посредника между всеми серверами (воркерами), для удобной работы
func NewServersManager(quantity int) *ServersManager {
	return &ServersManager{
		ServersQuantity: quantity,
		Servers:         ServersData{Servers: make(map[int]*Server)},
	}
}

// InitServers Добавление воркеров (серверов) исходя из их количества - переменная окружения
func (sm *ServersManager) InitServers() {
	sm.Servers.Mu.Lock()
	defer sm.Servers.Mu.Unlock()
	for serverID := 1; serverID <= sm.ServersQuantity; serverID++ {
		server := &Server{ID: serverID, Status: "Online, standing by", Tasks: "", LastPing: time.Now().Format("02-01-2006 15:04:05")}
		sm.Servers.Servers[serverID] = server
	}
}

// UpdateServers изменение данных воркера
func (sm *ServersManager) UpdateServers(id int, operation string, status string) {
	sm.Servers.Mu.Lock()
	defer sm.Servers.Mu.Unlock()
	server, exists := sm.Servers.Servers[id]
	if !exists {
		log.Println("Server with ID", id, "not found")
		return
	}
	server.Status = status
	server.Tasks = operation
	server.LastPing = time.Now().Format("02-01-2006 15:04:05")
}

// SendHeartbeat Посыл ответа от воркера
func (sm *ServersManager) SendHeartbeat(id int) {
	sm.Servers.Mu.Lock()
	defer sm.Servers.Mu.Unlock()

	server, exists := sm.Servers.Servers[id]
	if !exists {
		log.Println("Server with ID", id, "not found")
		return
	}

	server.LastPing = time.Now().Format("02-01-2006 15:04:05")
}

// RunServers запуск работы посыла ответа от воркеров
func (sm *ServersManager) RunServers() {
	sm.Servers.Mu.Lock()
	defer sm.Servers.Mu.Unlock()

	for id := range sm.Servers.Servers {
		go func(id int) {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					Servers.SendHeartbeat(id)
				}
			}
		}(id)
	}
}
