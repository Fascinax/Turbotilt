package runtime

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// mockCmd is used to simulate external commands
type mockCmd struct {
	executed bool
	command  string
	args     []string
}

var lastMockCmd mockCmd

func mockExecCommand(command string, args ...string) *exec.Cmd {
	// Record the called command
	lastMockCmd = mockCmd{
		executed: true,
		command:  command,
		args:     args,
	}

	// Create a command that does nothing (in test mode)
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

// TestHelperProcess is not a real test, it's a helper to simulate external commands
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// This function simulates an external command that completes successfully
	os.Exit(0)
}

func TestTiltUp(t *testing.T) {
	// Save the original exec.Command function and restore it after the test
	origExecCommand := execCommand
	defer func() { execCommand = origExecCommand }()

	// Replace with our mock
	execCommand = mockExecCommand

	// Reset the mock state
	lastMockCmd = mockCmd{}

	// Force Tilt to be detected as installed
	// (the mock will handle the execution)
	orig := isTiltInstalled
	defer func() { isTiltInstalled = orig }()
	isTiltInstalled = func() bool { return true }

	// Test with Tilt
	opts := RunOptions{
		UseTilt:  true,
		Detached: false,
		Debug:    true,
	}

	err := TiltUp(opts)
	if err != nil {
		t.Errorf("TiltUp returned an error: %v", err)
	}

	if !lastMockCmd.executed {
		t.Error("The command was not executed")
	}

	if lastMockCmd.command != "tilt" {
		t.Errorf("The command should be 'tilt', but it's '%s'", lastMockCmd.command)
	}

	// Verify that the debug option was correctly passed
	found := false
	for _, arg := range lastMockCmd.args {
		if arg == "--debug" {
			found = true
			break
		}
	}
	if !found {
		t.Error("The --debug option was not passed to the tilt command")
	}
}

func TestComposeUp(t *testing.T) {
	// Save the original exec.Command function and restore it after the test
	origExecCommand := execCommand
	defer func() { execCommand = origExecCommand }()

	// Replace with our mock
	execCommand = mockExecCommand

	// Reset the mock state
	lastMockCmd = mockCmd{}

	// Test with Docker Compose in detached mode and a specific service
	opts := RunOptions{
		Detached:    true,
		ServiceName: "api",
	}

	err := ComposeUp(opts)
	if err != nil {
		t.Errorf("ComposeUp returned an error: %v", err)
	}

	if !lastMockCmd.executed {
		t.Error("The command was not executed")
	}

	if lastMockCmd.command != "docker" {
		t.Errorf("The command should be 'docker', but it's '%s'", lastMockCmd.command)
	}

	// Verify that the correct arguments were passed
	expectedArgs := []string{"compose", "up", "-d", "api"}
	for i, arg := range expectedArgs {
		if i >= len(lastMockCmd.args) || lastMockCmd.args[i] != arg {
			t.Errorf("Argument %d should be '%s', but it's '%s'",
				i, arg, lastMockCmd.args[i])
		}
	}
}

func TestDryRun(t *testing.T) {
	// Save the original exec.Command function and restore it after the test
	origExecCommand := execCommand
	defer func() { execCommand = origExecCommand }()

	// Replace with our mock
	execCommand = mockExecCommand

	// Test in dry-run mode
	opts := RunOptions{
		UseTilt: true,
		DryRun:  true,
	}

	lastMockCmd = mockCmd{}
	err := TiltUp(opts)
	if err != nil {
		t.Errorf("TiltUp in dry-run mode returned an error: %v", err)
	}

	if lastMockCmd.executed {
		t.Error("In dry-run mode, the command should not be executed")
	}

	// Test for ComposeUp too
	lastMockCmd = mockCmd{}
	opts.UseTilt = false
	err = ComposeUp(opts)
	if err != nil {
		t.Errorf("ComposeUp in dry-run mode returned an error: %v", err)
	}

	if lastMockCmd.executed {
		t.Error("In dry-run mode, the command should not be executed")
	}
}

func TestSetupCleanup(t *testing.T) {
	// Create temporary files for the test
	tempDir := t.TempDir()
	tempFiles := []string{
		filepath.Join(tempDir, "file1.tmp"),
		filepath.Join(tempDir, "file2.tmp"),
	}

	for _, file := range tempFiles {
		if err := os.WriteFile(file, []byte("test content"), 0644); err != nil {
			t.Fatalf("Unable to create temporary file: %v", err)
		}
	}

	// Configure cleanup
	SetupCleanup(tempFiles)

	// Verify that the files still exist
	for _, file := range tempFiles {
		if _, err := os.Stat(file); err != nil {
			t.Errorf("File %s should exist: %v", file, err)
		}
	}

	// Note: We can't directly test the SIGINT/SIGTERM signal handler
	// in a unit test, so we just verify that the function doesn't crash
}
