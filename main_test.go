package notrhttp

import (
	"net/http"
	"testing"
)

// Test NewServer initializes correctly
func TestNewServer(t *testing.T) {
	server := NewServer("8080", "1.0.0")

	if server.Port != ":8080" {
		t.Errorf("Expected port ':8080', got '%s'", server.Port)
	}

	if server.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", server.Version)
	}

	if len(server.Routes) != 0 {
		t.Errorf("Expected no routes, got '%d'", len(server.Routes))
	}
}

// Test setting server name
func TestSetName(t *testing.T) {
	server := NewServer("8080", "1.0.0")
	server.SetName("Test Server")

	if server.Name != "Test Server" {
		t.Errorf("Expected name 'Test Server', got '%s'", server.Name)
	}
}

// Test adding a GET route
func TestGetRoute(t *testing.T) {
	server := NewServer("8080", "1.0.0")
	server.Get("/test", func(rw Writer, r *Request) {
		rw.RespondWithSuccess(map[string]string{"message": "Hello, world!"})
	})

	if len(server.Routes) != 1 {
		t.Errorf("Expected 1 route, got '%d'", len(server.Routes))
	}

	route := server.Routes[0]
	if route.Path != "/test" {
		t.Errorf("Expected path '/test', got '%s'", route.Path)
	}

	if route.Method != http.MethodGet {
		t.Errorf("Expected method 'GET', got '%s'", route.Method)
	}
}


