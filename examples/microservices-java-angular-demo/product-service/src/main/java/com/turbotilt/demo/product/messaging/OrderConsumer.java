package com.turbotilt.demo.product.messaging;

import com.turbotilt.demo.product.model.Product;
import com.turbotilt.demo.product.service.ProductService;
import io.quarkus.logging.Log;
import io.smallrye.reactive.messaging.annotations.Blocking;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import org.eclipse.microprofile.reactive.messaging.Incoming;

import java.util.concurrent.CompletionStage;

@ApplicationScoped
public class OrderConsumer {
    
    @Inject
    ProductService productService;
    
    @Incoming("orders-in")
    @Blocking
    public void processOrder(OrderEvent orderEvent) {
        Log.info("Received order event: " + orderEvent);
        
        // Update product stock based on order items
        for (OrderItem item : orderEvent.getItems()) {
            productService.updateStock(item.getProductId(), -item.getQuantity());
            Log.info("Updated stock for product " + item.getProductId() + " with quantity " + (-item.getQuantity()));
        }
    }
    
    public static class OrderEvent {
        private String orderId;
        private OrderItem[] items;
        
        public String getOrderId() {
            return orderId;
        }
        
        public void setOrderId(String orderId) {
            this.orderId = orderId;
        }
        
        public OrderItem[] getItems() {
            return items;
        }
        
        public void setItems(OrderItem[] items) {
            this.items = items;
        }
    }
    
    public static class OrderItem {
        private Long productId;
        private int quantity;
        
        public Long getProductId() {
            return productId;
        }
        
        public void setProductId(Long productId) {
            this.productId = productId;
        }
        
        public int getQuantity() {
            return quantity;
        }
        
        public void setQuantity(int quantity) {
            this.quantity = quantity;
        }
    }
}
