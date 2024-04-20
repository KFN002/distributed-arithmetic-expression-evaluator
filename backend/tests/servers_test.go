package tests

import (
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/pkg/models"
	"testing"
	"time"
)

func TestNewServersManager(t *testing.T) {
	quantity := 3
	sm := models.NewServersManager(quantity)

	if sm.ServersQuantity != quantity {
		t.Errorf("NewServersManager() failed to set ServersQuantity. Got: %d, Expected: %d", sm.ServersQuantity, quantity)
	}

	if sm.Servers.Servers == nil {
		t.Error("NewServersManager() failed to initialize Servers map")
	}
}

func TestInitServers(t *testing.T) {
	quantity := 3
	sm := models.NewServersManager(quantity)

	sm.InitServers()

	if len(sm.Servers.Servers) != quantity {
		t.Errorf("InitServers() failed to initialize correct number of servers. Got: %d, Expected: %d", len(sm.Servers.Servers), quantity)
	}
}

func TestUpdateServers(t *testing.T) {
	sm := models.NewServersManager(1)
	sm.InitServers()

	id := 1
	operation := "test operation"
	status := "test status"

	sm.UpdateServers(id, operation, status)

	server, exists := sm.Servers.Servers[id]
	if !exists {
		t.Fatalf("UpdateServers() failed to find server with ID: %d", id)
	}

	if server.Tasks != operation {
		t.Errorf("UpdateServers() failed to set correct operation. Got: %s, Expected: %s", server.Tasks, operation)
	}

	if server.Status != status {
		t.Errorf("UpdateServers() failed to set correct status. Got: %s, Expected: %s", server.Status, status)
	}

	lastPing := time.Now().Format("02-01-2006 15:04:05")
	if server.LastPing != lastPing {
		t.Errorf("UpdateServers() failed to update LastPing. Got: %s, Expected: %s", server.LastPing, lastPing)
	}
}

func TestSendHeartbeat(t *testing.T) {
	sm := models.NewServersManager(1)
	sm.InitServers()

	id := 1

	sm.SendHeartbeat(id)

	server, exists := sm.Servers.Servers[id]
	if !exists {
		t.Fatalf("SendHeartbeat() failed to find server with ID: %d", id)
	}

	lastPing := time.Now().Format("02-01-2006 15:04:05")
	if server.LastPing != lastPing {
		t.Errorf("SendHeartbeat() failed to update LastPing. Got: %s, Expected: %s", server.LastPing, lastPing)
	}
}

func TestRunServers(t *testing.T) {
	sm := models.NewServersManager(1)
	sm.InitServers()

	go sm.RunServers()

	time.Sleep(3 * time.Second)

	for _, server := range sm.Servers.Servers {
		lastPing := time.Now().Add(-3 * time.Second).Format("02-01-2006 15:04:05")
		if server.LastPing != lastPing {
			t.Errorf("RunServers() failed to update LastPing for server %d. Got: %s, Expected: %s", server.ID, server.LastPing, lastPing)
		}
	}
}
