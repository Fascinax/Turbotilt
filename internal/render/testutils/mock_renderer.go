package testutils

import "io"

// MockDockerfileRenderer is a mock implementation for testing
type MockDockerfileRenderer struct{}

// RenderSpringDockerfile is a mock for testing
func (m *MockDockerfileRenderer) RenderSpringDockerfile(w io.Writer, opts interface{}) error {
	// Mock implementation for tests
	return nil
}

// RenderQuarkusDockerfile is a mock for testing
func (m *MockDockerfileRenderer) RenderQuarkusDockerfile(w io.Writer, opts interface{}) error {
	// Mock implementation for tests
	return nil
}

// RenderMicronautDockerfile is a mock for testing
func (m *MockDockerfileRenderer) RenderMicronautDockerfile(w io.Writer, opts interface{}) error {
	// Mock implementation for tests
	return nil
}

// RenderJavaDockerfile is a mock for testing
func (m *MockDockerfileRenderer) RenderJavaDockerfile(w io.Writer, opts interface{}) error {
	// Mock implementation for tests
	return nil
}

// RenderGenericDockerfile is a mock for testing
func (m *MockDockerfileRenderer) RenderGenericDockerfile(w io.Writer, opts interface{}) error {
	// Mock implementation for tests
	return nil
}

// NewMockRenderer creates a new MockDockerfileRenderer
func NewMockRenderer() *MockDockerfileRenderer {
	return &MockDockerfileRenderer{}
}
