package scan

import (
	"os"
)

func DetectFramework() string {
	if _, err := os.Stat("pom.xml"); err == nil {
		// TODO: Lire pom.xml et détecter Spring ou Quarkus
		return "spring"
	}
	// TODO: Ajouter détection gradle, quarkus.*
	return "unknown"
}
