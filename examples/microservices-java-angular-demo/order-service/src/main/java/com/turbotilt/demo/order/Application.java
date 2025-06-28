package com.turbotilt.demo.order;

import io.micronaut.runtime.Micronaut;
import io.swagger.v3.oas.annotations.OpenAPIDefinition;
import io.swagger.v3.oas.annotations.info.Contact;
import io.swagger.v3.oas.annotations.info.Info;

@OpenAPIDefinition(
    info = @Info(
        title = "Order Service",
        version = "1.0",
        description = "Order Service API using Micronaut",
        contact = @Contact(name = "Turbotilt", email = "contact@turbotilt.example")
    )
)
public class Application {
    public static void main(String[] args) {
        Micronaut.run(Application.class, args);
    }
}
