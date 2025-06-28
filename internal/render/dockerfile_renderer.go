package render

import (
	"io"
	"text/template"
)

// DockerfileTemplates contains all Dockerfile templates
const (
	SpringDockerfileTmpl = `FROM eclipse-temurin:{{.JDKVersion}} AS build
WORKDIR /app
COPY . .
RUN ./mvnw package {{if .DevMode}}-DskipTests{{end}}

FROM eclipse-temurin:{{.JDKVersion}}
WORKDIR /app
COPY --from=build /app/target/*.jar app.jar
EXPOSE {{.Port}}
CMD ["java", "-jar", "app.jar"]
`

	QuarkusDockerfileTmpl = `FROM registry.access.redhat.com/ubi8/openjdk-{{.JDKVersion}}:latest AS build
WORKDIR /app
COPY . .
RUN ./mvnw package {{if .DevMode}}-DskipTests{{end}} -Dquarkus.package.type=jar

FROM registry.access.redhat.com/ubi8/openjdk-{{.JDKVersion}}:latest
WORKDIR /app
COPY --from=build /app/target/quarkus-app/lib/ /deployments/lib/
COPY --from=build /app/target/quarkus-app/*.jar /deployments/
COPY --from=build /app/target/quarkus-app/app/ /deployments/app/
COPY --from=build /app/target/quarkus-app/quarkus/ /deployments/quarkus/
EXPOSE {{.Port}}
CMD ["java", "-jar", "/deployments/quarkus-run.jar"]
`

	MicronautDockerfileTmpl = `FROM eclipse-temurin:{{.JDKVersion}} AS build
WORKDIR /app
COPY . .
RUN ./gradlew build {{if .DevMode}}-x test{{end}}

FROM eclipse-temurin:{{.JDKVersion}}
WORKDIR /app
COPY --from=build /app/build/libs/*-all.jar app.jar
EXPOSE {{.Port}}
CMD ["java", "-jar", "app.jar"]
`

	JavaDockerfileTmpl = `FROM eclipse-temurin:{{.JDKVersion}}
WORKDIR /app
COPY . .
RUN javac Main.java
EXPOSE {{.Port}}
CMD ["java", "Main"]
`

	GenericDockerfileTmpl = `FROM alpine:latest
WORKDIR /app
COPY . .
EXPOSE {{.Port}}
CMD ["sh", "start.sh"]
`
)

// TemplateDockerfileRenderer is the implementation of DockerfileRenderer that uses templates
type TemplateDockerfileRenderer struct{}

// renderDockerfile executes a Dockerfile template with the specified options
func (r *TemplateDockerfileRenderer) renderDockerfile(w io.Writer, tmplContent string, name string, opts Options) error {
	tmpl, err := template.New(name).Parse(tmplContent)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, opts)
}

// RenderSpringDockerfile writes a Dockerfile for Spring Boot
func (r *TemplateDockerfileRenderer) RenderSpringDockerfile(w io.Writer, opts Options) error {
	return r.renderDockerfile(w, SpringDockerfileTmpl, "spring", opts)
}

// RenderQuarkusDockerfile writes a Dockerfile for Quarkus
func (r *TemplateDockerfileRenderer) RenderQuarkusDockerfile(w io.Writer, opts Options) error {
	return r.renderDockerfile(w, QuarkusDockerfileTmpl, "quarkus", opts)
}

// RenderMicronautDockerfile writes a Dockerfile for Micronaut
func (r *TemplateDockerfileRenderer) RenderMicronautDockerfile(w io.Writer, opts Options) error {
	return r.renderDockerfile(w, MicronautDockerfileTmpl, "micronaut", opts)
}

// RenderJavaDockerfile writes a Dockerfile for a generic Java application
func (r *TemplateDockerfileRenderer) RenderJavaDockerfile(w io.Writer, opts Options) error {
	return r.renderDockerfile(w, JavaDockerfileTmpl, "java", opts)
}

// RenderGenericDockerfile writes a generic Dockerfile for other application types
func (r *TemplateDockerfileRenderer) RenderGenericDockerfile(w io.Writer, opts Options) error {
	return r.renderDockerfile(w, GenericDockerfileTmpl, "generic", opts)
}

var defaultRenderer DockerfileRenderer = &TemplateDockerfileRenderer{}
