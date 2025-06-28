package render

// ComposeTemplates contains templates for docker-compose.yml
const (
	// ComposeTemplateWithEnvFile is the docker-compose.yml template with an environment file
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

	// ComposeTemplate is the docker-compose.yml template without an environment file
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
