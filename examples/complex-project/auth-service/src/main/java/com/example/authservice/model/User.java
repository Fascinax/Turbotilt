package com.example.authservice.model;

import io.quarkus.hibernate.orm.panache.PanacheEntity;
import jakarta.persistence.Entity;
import lombok.Getter;
import lombok.Setter;

@Entity
@Getter
@Setter
public class User extends PanacheEntity {
    private String username;
    private String password;
    private String email;
    private boolean active;
    private String role;
}
