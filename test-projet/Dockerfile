FROM quay.io/quarkus/ubi-quarkus-native-image:22.3-java21 AS build
WORKDIR /app
COPY . .
RUN ./mvnw package -Pnative -DskipTests

FROM quay.io/quarkus/quarkus-micro-image:2.0
WORKDIR /app
COPY --from=build /app/target/*-runner /app/application
EXPOSE 9090
CMD ["./application", "-Dquarkus.http.host=0.0.0.0"]
