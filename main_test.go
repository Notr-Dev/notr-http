package notrhttp

import (
	"net/http"
	"testing"
	"time"
)

// Test NewServer initializes correctly
func TestNewServer(t *testing.T) {
	server := NewServer(
		Server{
			Name:    "Test Server",
			Port:    ":8080",
			Version: "1.0.0",
		},
	)

	if server.Port != ":8080" {
		t.Errorf("Expected port ':8080', got '%s'", server.Port)
	}

	if server.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", server.Version)
	}

	if len(server.routes) != 0 {
		t.Errorf("Expected no routes, got '%d'", len(server.routes))
	}
}

// Test setting server name
func TestSetName(t *testing.T) {
	server := NewServer(
		Server{
			Name:    "Test Server",
			Port:    ":8080",
			Version: "1.0.0",
		},
	)

	if server.Name != "Test Server" {
		t.Errorf("Expected name 'Test Server', got '%s'", server.Name)
	}
}

// Test adding a GET route
func TestGetRoute(t *testing.T) {
	server := NewServer(
		Server{
			Name:    "Test Server",
			Port:    ":8080",
			Version: "1.0.0",
		},
	)

	server.Get("/test", func(rw Writer, r *Request) {
		rw.RespondWithSuccess(map[string]string{"message": "Hello, world!"})
	})

	if len(server.routes) != 1 {
		t.Errorf("Expected 1 route, got '%d'", len(server.routes))
	}

	route := server.routes[0]
	if route.Path != "/test" {
		t.Errorf("Expected path '/test', got '%s'", route.Path)
	}

	if route.Method != http.MethodGet {
		t.Errorf("Expected method 'GET', got '%s'", route.Method)
	}
}

// Test fetching home route
func TestGetHomeRoute(t *testing.T) {
	server := NewServer(
		Server{
			Name:    "Test Server",
			Port:    ":8080",
			Version: "1.0.0",
		},
	)

	server.Get("/", func(rw Writer, r *Request) {
		rw.RespondWithSuccess(map[string]string{"message": "Hello, world!"})
	})

	stopChan := make(chan struct{})
	errChan := make(chan error, 1)
	go func() {
		select {
		case <-stopChan:
			return
		default:
			if err := server.Run(); err != nil && err != http.ErrServerClosed {
				errChan <- err
			}
		}
	}()

	select {
	case err := <-errChan:
		t.Fatalf("Failed to start server: %v", err)
	case <-time.After(100 * time.Millisecond):
	}

	time.Sleep(1000 * time.Millisecond)

	resp, err := http.Get("http://localhost:8080/")
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status '200 OK', got '%d'", resp.StatusCode)
	}

	time.Sleep(100 * time.Millisecond)
	close(stopChan)
}

// Test fetching a non-existent route
func TestGetRouteNotFound(t *testing.T) {
	server := NewServer(
		Server{
			Name:    "Test Server",
			Port:    ":8080",
			Version: "1.0.0",
		},
	)

	server.Get("/test", func(rw Writer, r *Request) {
		rw.RespondWithSuccess(map[string]string{"message": "Hello, world!"})
	})

	stopChan := make(chan struct{})
	errChan := make(chan error, 1)
	go func() {
		select {
		case <-stopChan:
			return
		default:
			if err := server.Run(); err != nil && err != http.ErrServerClosed {
				errChan <- err
			}
		}
	}()

	select {
	case err := <-errChan:
		t.Fatalf("Failed to start server: %v", err)
	case <-time.After(100 * time.Millisecond):
	}

	time.Sleep(1000 * time.Millisecond)

	resp, err := http.Get("http://localhost:8080/nonexistent")
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status '404 Not Found', got '%d'", resp.StatusCode)
	}

	time.Sleep(100 * time.Millisecond)

	close(stopChan)
}
