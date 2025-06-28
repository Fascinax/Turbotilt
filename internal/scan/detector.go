package scan

// BaseDetector contains common functionality for all detectors
type BaseDetector struct {
	// Common fields for all detectors
}

// DetectionResult contains the results of the framework detection
type DetectionResult struct {
	Framework    string            // The detected framework name (spring, quarkus, micronaut, etc.)
	BuildSystem  string            // The build system (maven, gradle, etc.)
	JavaVersion  string            // Java version
	Port         string            // Application port
	Dependencies []string          // Project dependencies
	Properties   map[string]string // Configuration properties
	Detected     bool              // Indicates if the framework was detected
}
