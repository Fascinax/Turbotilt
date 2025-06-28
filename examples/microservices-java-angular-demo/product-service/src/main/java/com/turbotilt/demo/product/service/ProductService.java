package com.turbotilt.demo.product.service;

import com.turbotilt.demo.product.model.Product;
import com.turbotilt.demo.product.repository.ProductRepository;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import jakarta.transaction.Transactional;

import java.util.List;
import java.util.Optional;

@ApplicationScoped
public class ProductService {
    
    @Inject
    ProductRepository productRepository;
    
    public List<Product> findAllProducts() {
        return productRepository.listAll();
    }
    
    public Optional<Product> findProductById(Long id) {
        return Optional.ofNullable(productRepository.findById(id));
    }
    
    public List<Product> findProductsByCategory(String category) {
        return productRepository.findByCategory(category);
    }
    
    public List<Product> findProductsByName(String name) {
        return productRepository.findByNameContaining(name);
    }
    
    @Transactional
    public Product saveProduct(Product product) {
        productRepository.persist(product);
        return product;
    }
    
    @Transactional
    public boolean deleteProduct(Long id) {
        return productRepository.deleteById(id);
    }
    
    @Transactional
    public boolean updateStock(Long id, int quantity) {
        Product product = productRepository.findById(id);
        if (product == null) {
            return false;
        }
        
        product.setStockQuantity(product.getStockQuantity() + quantity);
        return true;
    }
}
