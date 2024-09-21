package notrhttp

import (
	"net/http"
	"testing"
	"time"
)

// Helper function to initialize a new server
func setupServer(port string) *Server {
	server := NewServer(Server{
		Name:    "Test Server",
		Port:    port,
		Version: "1.0.0",
	})

	return server
}

// Helper function to run the server in a goroutine
func runServer(t *testing.T, server *Server, stopChan chan struct{}) {
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
}

// Helper function to make a GET request and check the response status
func checkGetRequest(t *testing.T, url string, expectedStatus int) {
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		t.Errorf("Expected status '%d', got '%d'", expectedStatus, resp.StatusCode)
	}
}

// Test NewServer initializes correctly
func TestNewServer(t *testing.T) {
	server := setupServer(":0001")

	if server.Port != ":0001" {
		t.Errorf("Expected port ':0001', got '%s'", server.Port)
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
	server := setupServer(":0002")

	if server.Name != "Test Server" {
		t.Errorf("Expected name 'Test Server', got '%s'", server.Name)
	}
}

// Test adding a GET route
func TestGetRoute(t *testing.T) {
	server := setupServer(":0003")

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

// Test fetching home route
func TestGetHomeRoute(t *testing.T) {
	server := setupServer(":0004")

	server.Get("/", func(rw Writer, r *Request) {
		rw.RespondWithSuccess(map[string]string{"message": "Hello, world!"})
	})

	stopChan := make(chan struct{})
	runServer(t, server, stopChan)

	time.Sleep(1000 * time.Millisecond)

	checkGetRequest(t, "http://localhost:0004/", http.StatusOK)

	close(stopChan)
}

// Test fetching a non-existent route
func TestGetRouteNotFound(t *testing.T) {
	server := setupServer(":0005")

	server.Get("/test", func(rw Writer, r *Request) {
		rw.RespondWithSuccess(map[string]string{"message": "Hello, world!"})
	})

	stopChan := make(chan struct{})
	runServer(t, server, stopChan)

	time.Sleep(1000 * time.Millisecond)

	checkGetRequest(t, "http://localhost:0005/nonexistent", http.StatusNotFound)

	close(stopChan)
}
