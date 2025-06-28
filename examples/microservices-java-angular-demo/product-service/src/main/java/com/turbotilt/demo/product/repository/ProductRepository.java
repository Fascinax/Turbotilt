package com.turbotilt.demo.product.repository;

import com.turbotilt.demo.product.model.Product;
import io.quarkus.hibernate.orm.panache.PanacheRepository;
import jakarta.enterprise.context.ApplicationScoped;

import java.util.List;

@ApplicationScoped
public class ProductRepository implements PanacheRepository<Product> {
    
    public List<Product> findByCategory(String category) {
        return list("category", category);
    }
    
    public List<Product> findByNameContaining(String name) {
        return list("name LIKE ?1", "%" + name + "%");
    }
}
