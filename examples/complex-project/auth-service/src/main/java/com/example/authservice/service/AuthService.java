package com.example.authservice.service;

import com.example.authservice.dto.AuthRequest;
import com.example.authservice.dto.AuthResponse;
import com.example.authservice.model.User;
import io.smallrye.jwt.build.Jwt;
import jakarta.enterprise.context.ApplicationScoped;
import org.eclipse.microprofile.config.inject.ConfigProperty;

import java.time.Duration;
import java.util.Arrays;
import java.util.HashSet;

@ApplicationScoped
public class AuthService {

    @ConfigProperty(name = "mp.jwt.verify.issuer")
    String issuer;

    public AuthResponse login(AuthRequest authRequest) {
        // In a real application, you would validate against a database
        // This is a simplified example
        User user = User.find("username", authRequest.getUsername()).firstResult();
        
        if (user == null || !authRequest.getPassword().equals(user.getPassword())) {
            throw new RuntimeException("Invalid username or password");
        }
        
        String token = generateToken(user);
        
        return new AuthResponse(token, user.getUsername(), user.getRole());
    }

    public void register(AuthRequest authRequest) {
        // Check if user already exists
        if (User.find("username", authRequest.getUsername()).count() > 0) {
            throw new RuntimeException("Username already exists");
        }
        
        User user = new User();
        user.setUsername(authRequest.getUsername());
        user.setPassword(authRequest.getPassword()); // In a real app, you would hash this
        user.setActive(true);
        user.setRole("USER");
        
        user.persist();
    }

    public boolean validateToken(String token) {
        // In a real application, you would validate the JWT token
        // This is a simplified example
        return token != null && token.startsWith("Bearer ");
    }

    private String generateToken(User user) {
        return Jwt.issuer(issuer)
                .upn(user.getUsername())
                .groups(new HashSet<>(Arrays.asList(user.getRole())))
                .expiresIn(Duration.ofHours(24))
                .sign();
    }
}
