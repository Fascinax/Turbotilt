package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Vérifie installation & config",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔍 Vérification de l'environnement Turbotilt...")

		// Tableau pour stocker les résultats
		results := make(map[string]bool)

		// Vérifier Docker
		fmt.Print("⏳ Docker : ")
		if version, err := execCommand("docker", "--version"); err == nil {
			fmt.Printf("✅ %s\n", version)
			results["docker"] = true
		} else {
			fmt.Println("❌ Non installé ou inaccessible")
			results["docker"] = false
		}

		// Vérifier Docker Compose
		fmt.Print("⏳ Docker Compose : ")
		if version, err := execCommand("docker", "compose", "version"); err == nil {
			fmt.Printf("✅ %s\n", strings.Split(version, "\n")[0])
			results["docker-compose"] = true
		} else {
			fmt.Println("❌ Non installé ou inaccessible")
			results["docker-compose"] = false
		}

		// Vérifier Tilt
		fmt.Print("⏳ Tilt : ")
		if version, err := execCommand("tilt", "version"); err == nil {
			fmt.Printf("✅ %s\n", strings.Split(version, "\n")[0])
			results["tilt"] = true
		} else {
			fmt.Println("❌ Non installé ou inaccessible")
			fmt.Println("   👉 Pour installer Tilt: https://docs.tilt.dev/install.html")
			results["tilt"] = false
		}

		// Vérifier Java
		fmt.Print("⏳ Java : ")
		if version, err := execCommand("java", "-version"); err == nil {
			fmt.Printf("✅ %s\n", strings.Split(version, "\n")[0])
			results["java"] = true
		} else {
			fmt.Println("❌ Non installé ou inaccessible")
			results["java"] = false
		}

		// Vérifier Maven
		fmt.Print("⏳ Maven : ")
		if version, err := execCommand("mvn", "--version"); err == nil {
			fmt.Printf("✅ %s\n", strings.Split(version, "\n")[0])
			results["maven"] = true
		} else {
			fmt.Println("⚠️ Non installé (requis pour certains projets)")
			results["maven"] = false
		}

		// Vérifier Gradle
		fmt.Print("⏳ Gradle : ")
		if version, err := execCommand("gradle", "--version"); err == nil {
			fmt.Printf("✅ %s\n", strings.Split(strings.Split(version, "\n")[0], "----")[0])
			results["gradle"] = true
		} else {
			fmt.Println("⚠️ Non installé (requis pour certains projets)")
			results["gradle"] = false
		}

		// Vérifier le projet courant
		fmt.Println("\n📋 Projet courant :")

		// Vérifier si les fichiers requis existent
		checkProjectFiles()

		// Afficher les recommandations
		fmt.Println("\n📋 Recommandations :")
		if !results["docker"] {
			fmt.Println("❗ Docker est requis : https://docs.docker.com/get-docker/")
		}
		if !results["docker-compose"] {
			fmt.Println("❗ Docker Compose est requis : https://docs.docker.com/compose/install/")
		}
		if !results["tilt"] {
			fmt.Println("❗ Tilt est recommandé : https://docs.tilt.dev/install.html")
		}

		fmt.Println("\n🔧 Commandes disponibles :")
		fmt.Println("▶️ turbotilt init   : Initialiser un projet")
		fmt.Println("▶️ turbotilt up     : Démarrer l'environnement")
		fmt.Println("▶️ turbotilt doctor : Vérifier la configuration")
	},
}

func execCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func checkProjectFiles() {
	fileChecks := []struct {
		path string
		desc string
	}{
		{"pom.xml", "Maven"},
		{"build.gradle", "Gradle"},
		{"Dockerfile", "Docker"},
		{"docker-compose.yml", "Docker Compose"},
		{"Tiltfile", "Tilt"},
	}

	for _, check := range fileChecks {
		if _, err := os.Stat(check.path); err == nil {
			fmt.Printf("✅ %s trouvé (%s)\n", check.path, check.desc)
		} else {
			fmt.Printf("❌ %s non trouvé\n", check.path)
		}
	}
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
