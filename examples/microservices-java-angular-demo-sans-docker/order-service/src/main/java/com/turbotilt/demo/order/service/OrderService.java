package com.turbotilt.demo.order.service;

import com.turbotilt.demo.order.model.Order;
import com.turbotilt.demo.order.repository.OrderRepository;
import com.turbotilt.demo.order.producer.OrderProducer;
import jakarta.inject.Singleton;
import lombok.RequiredArgsConstructor;

import java.util.ArrayList;
import java.util.List;
import java.util.Optional;

@Singleton
@RequiredArgsConstructor
public class OrderService {
    
    private final OrderRepository orderRepository;
    private final OrderProducer orderProducer;
    
    public List<Order> findAllOrders() {
        List<Order> orders = new ArrayList<>();
        orderRepository.findAll().forEach(orders::add);
        return orders;
    }
    
    public Optional<Order> findOrderById(Long id) {
        return orderRepository.findById(id);
    }
    
    public List<Order> findOrdersByUserId(Long userId) {
        return orderRepository.findByUserId(userId);
    }
    
    public List<Order> findOrdersByStatus(String status) {
        return orderRepository.findByStatus(status);
    }
    
    public Order createOrder(Order order) {
        order.setStatus("CREATED");
        Order savedOrder = orderRepository.save(order);
        
        // Publish the order created event
        orderProducer.sendOrderCreatedEvent(savedOrder);
        
        return savedOrder;
    }
    
    public Optional<Order> updateOrderStatus(Long id, String status) {
        Optional<Order> orderOpt = orderRepository.findById(id);
        
        if (orderOpt.isPresent()) {
            Order order = orderOpt.get();
            order.setStatus(status);
            return Optional.of(orderRepository.update(order));
        }
        
        return Optional.empty();
    }
    
    public void deleteOrder(Long id) {
        orderRepository.deleteById(id);
    }
}
