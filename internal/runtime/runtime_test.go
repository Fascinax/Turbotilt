package runtime

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// mockCmd est utilisé pour simuler les commandes externes
type mockCmd struct {
	executed bool
	command  string
	args     []string
}

var lastMockCmd mockCmd

func mockExecCommand(command string, args ...string) *exec.Cmd {
	// Enregistrer la commande appelée
	lastMockCmd = mockCmd{
		executed: true,
		command:  command,
		args:     args,
	}
	
	// Créer une commande qui ne fait rien (en mode test)
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

// TestHelperProcess n'est pas un vrai test, c'est un helper pour simuler des commandes externes
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// Cette fonction simule une commande externe qui se termine avec succès
	os.Exit(0)
}

func TestTiltUp(t *testing.T) {
	// Sauvegarder la fonction exec.Command originale et la restaurer après le test
	origExecCommand := execCommand
	defer func() { execCommand = origExecCommand }()
	
	// Remplacer par notre mock
	execCommand = mockExecCommand
	
	// Réinitialiser l'état du mock
	lastMockCmd = mockCmd{}
	
	// Force la détection de Tilt comme installé
	// (le mock se chargera de simuler l'exécution)
	orig := isTiltInstalled
	defer func() { isTiltInstalled = orig }()
	isTiltInstalled = func() bool { return true }
	
	// Test avec Tilt
	opts := RunOptions{
		UseTilt:  true,
		Detached: false,
		Debug:    true,
	}
	
	err := TiltUp(opts)
	if err != nil {
		t.Errorf("TiltUp a retourné une erreur: %v", err)
	}
	
	if !lastMockCmd.executed {
		t.Error("La commande n'a pas été exécutée")
	}
	
	if lastMockCmd.command != "tilt" {
		t.Errorf("La commande devrait être 'tilt', mais c'est '%s'", lastMockCmd.command)
	}
	
	// Vérifier que l'option debug a été correctement passée
	found := false
	for _, arg := range lastMockCmd.args {
		if arg == "--debug" {
			found = true
			break
		}
	}
	if !found {
		t.Error("L'option --debug n'a pas été passée à la commande tilt")
	}
}

func TestComposeUp(t *testing.T) {
	// Sauvegarder la fonction exec.Command originale et la restaurer après le test
	origExecCommand := execCommand
	defer func() { execCommand = origExecCommand }()
	
	// Remplacer par notre mock
	execCommand = mockExecCommand
	
	// Réinitialiser l'état du mock
	lastMockCmd = mockCmd{}
	
	// Test avec Docker Compose en mode détaché et un service spécifique
	opts := RunOptions{
		Detached:    true,
		ServiceName: "api",
	}
	
	err := ComposeUp(opts)
	if err != nil {
		t.Errorf("ComposeUp a retourné une erreur: %v", err)
	}
	
	if !lastMockCmd.executed {
		t.Error("La commande n'a pas été exécutée")
	}
	
	if lastMockCmd.command != "docker" {
		t.Errorf("La commande devrait être 'docker', mais c'est '%s'", lastMockCmd.command)
	}
	
	// Vérifier que les arguments corrects ont été passés
	expectedArgs := []string{"compose", "up", "-d", "api"}
	for i, arg := range expectedArgs {
		if i >= len(lastMockCmd.args) || lastMockCmd.args[i] != arg {
			t.Errorf("L'argument %d devrait être '%s', mais c'est '%s'", 
				i, arg, lastMockCmd.args[i])
		}
	}
}

func TestDryRun(t *testing.T) {
	// Sauvegarder la fonction exec.Command originale et la restaurer après le test
	origExecCommand := execCommand
	defer func() { execCommand = origExecCommand }()
	
	// Remplacer par notre mock
	execCommand = mockExecCommand
	
	// Test en mode dry-run
	opts := RunOptions{
		UseTilt: true,
		DryRun:  true,
	}
	
	lastMockCmd = mockCmd{}
	err := TiltUp(opts)
	if err != nil {
		t.Errorf("TiltUp en mode dry-run a retourné une erreur: %v", err)
	}
	
	if lastMockCmd.executed {
		t.Error("En mode dry-run, la commande ne devrait pas être exécutée")
	}
	
	// Test pour ComposeUp aussi
	lastMockCmd = mockCmd{}
	opts.UseTilt = false
	err = ComposeUp(opts)
	if err != nil {
		t.Errorf("ComposeUp en mode dry-run a retourné une erreur: %v", err)
	}
	
	if lastMockCmd.executed {
		t.Error("En mode dry-run, la commande ne devrait pas être exécutée")
	}
}

func TestSetupCleanup(t *testing.T) {
	// Créer des fichiers temporaires pour le test
	tempDir := t.TempDir()
	tempFiles := []string{
		filepath.Join(tempDir, "file1.tmp"),
		filepath.Join(tempDir, "file2.tmp"),
	}
	
	for _, file := range tempFiles {
		if err := os.WriteFile(file, []byte("test content"), 0644); err != nil {
			t.Fatalf("Impossible de créer le fichier temporaire: %v", err)
		}
	}
	
	// Configurer le nettoyage
	SetupCleanup(tempFiles)
	
	// Vérifier que les fichiers existent toujours
	for _, file := range tempFiles {
		if _, err := os.Stat(file); err != nil {
			t.Errorf("Le fichier %s devrait exister: %v", file, err)
		}
	}
	
	// Note: On ne peut pas tester directement le handler de signal SIGINT/SIGTERM
	// dans un test unitaire, donc on vérifie seulement que la fonction ne crash pas
}
