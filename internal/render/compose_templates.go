package render

// ComposeTemplates contient les templates pour docker-compose.yml
const (
	// ComposeTemplateWithEnvFile est le template de docker-compose.yml avec un fichier d'environnement
	ComposeTemplateWithEnvFile = `version: '3'
services:
  app:
    build: .
    ports:
      - '{{.Port}}:{{.Port}}'
    volumes:
      - './src:/app/src'
    env_file:
      - {{.EnvFile}}
    environment:
{{if eq .Framework "spring"}}      - SPRING_PROFILES_ACTIVE={{if .DevMode}}dev{{else}}prod{{end}}
{{else if eq .Framework "quarkus"}}      - QUARKUS_PROFILE={{if .DevMode}}dev{{else}}prod{{end}}
{{else if eq .Framework "micronaut"}}      - MICRONAUT_ENVIRONMENTS={{if .DevMode}}dev{{else}}prod{{end}}
{{else}}      # Ajoutez vos variables d'environnement spécifiques ici
{{end}}
{{range .Services}}
  {{.Name}}:
    image: {{.Image}}
    ports:
      - '{{.Port}}'
    environment:
      - SPRING_PROFILES_ACTIVE={{if $.DevMode}}dev{{else}}prod{{end}}
{{end}}`

	// ComposeTemplate est le template de docker-compose.yml sans fichier d'environnement
	ComposeTemplate = `version: '3'
services:
  app:
    build: .
    ports:
      - '{{.Port}}:{{.Port}}'
    volumes:
      - './src:/app/src'
    environment:
{{if eq .Framework "spring"}}      - SPRING_PROFILES_ACTIVE={{if .DevMode}}dev{{else}}prod{{end}}
{{else if eq .Framework "quarkus"}}      - QUARKUS_PROFILE={{if .DevMode}}dev{{else}}prod{{end}}
{{else if eq .Framework "micronaut"}}      - MICRONAUT_ENVIRONMENTS={{if .DevMode}}dev{{else}}prod{{end}}
{{else}}      # Ajoutez vos variables d'environnement spécifiques ici
{{end}}
{{range .Services}}
  {{.Name}}:
    image: {{.Image}}
    ports:
      - '{{.Port}}'
    environment:
      - SPRING_PROFILES_ACTIVE={{if $.DevMode}}dev{{else}}prod{{end}}
{{end}}`
)
