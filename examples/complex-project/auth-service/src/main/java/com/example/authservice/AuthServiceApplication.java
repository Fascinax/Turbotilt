package com.example.authservice;

import io.quarkus.runtime.Quarkus;
import io.quarkus.runtime.annotations.QuarkusMain;

@QuarkusMain
public class AuthServiceApplication {
    public static void main(String[] args) {
        Quarkus.run(args);
    }
}
