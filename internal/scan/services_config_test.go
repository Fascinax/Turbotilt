package scan

import (
	"testing"
)

// TestServiceConfig tests the ServiceConfig struct
func TestServiceConfig(t *testing.T) {
	// Test creating a service config
	service := ServiceConfig{
		Type:    MySQL,
		Version: "8.0",
		Port:    "3306",
		Credentials: map[string]string{
			"username": "root",
			"password": "password",
		},
	}

	// Verify the service config
	if service.Type != MySQL {
		t.Errorf("Expected service type %s, got %s", MySQL, service.Type)
	}

	if service.Version != "8.0" {
		t.Errorf("Expected service version %s, got %s", "8.0", service.Version)
	}

	if service.Port != "3306" {
		t.Errorf("Expected service port %s, got %s", "3306", service.Port)
	}

	if service.Credentials["username"] != "root" {
		t.Errorf("Expected username %s, got %s", "root", service.Credentials["username"])
	}

	if service.Credentials["password"] != "password" {
		t.Errorf("Expected password %s, got %s", "password", service.Credentials["password"])
	}
}
