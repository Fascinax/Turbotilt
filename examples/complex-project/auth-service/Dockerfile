FROM registry.access.redhat.com/ubi8/openjdk-17:latest

ENV LANGUAGE='en_US:en'

WORKDIR /deployments

COPY --chown=jboss:jboss target/quarkus-app/lib/ /deployments/lib/
COPY --chown=jboss:jboss target/quarkus-app/*.jar /deployments/
COPY --chown=jboss:jboss target/quarkus-app/app/ /deployments/app/
COPY --chown=jboss:jboss target/quarkus-app/quarkus/ /deployments/quarkus/

EXPOSE 8082

CMD ["java", "-jar", "quarkus-run.jar"]
