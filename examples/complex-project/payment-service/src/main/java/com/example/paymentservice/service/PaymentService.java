package com.example.paymentservice.service;

import com.example.paymentservice.model.Payment;
import com.example.paymentservice.repository.PaymentRepository;
import jakarta.inject.Singleton;
import lombok.RequiredArgsConstructor;

import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Singleton
@RequiredArgsConstructor
public class PaymentService {
    
    private final PaymentRepository paymentRepository;
    
    public List<Payment> findAllPayments() {
        List<Payment> payments = new ArrayList<>();
        paymentRepository.findAll().forEach(payments::add);
        return payments;
    }
    
    public Optional<Payment> findPaymentById(Long id) {
        return paymentRepository.findById(id);
    }
    
    public List<Payment> findPaymentsByUserId(Long userId) {
        return paymentRepository.findByUserId(userId);
    }
    
    public Payment processPayment(Payment payment) {
        // Generate transaction ID
        payment.setTransactionId(UUID.randomUUID().toString());
        
        // Set creation timestamp
        payment.setCreatedAt(LocalDateTime.now());
        
        // Process payment logic would go here in a real application
        // For demo purposes, we'll just set the status to PROCESSED
        payment.setStatus("PROCESSED");
        
        return paymentRepository.save(payment);
    }
    
    public Payment updatePayment(Long id, Payment paymentDetails) {
        Payment payment = paymentRepository.findById(id)
                .orElseThrow(() -> new RuntimeException("Payment not found with id: " + id));
        
        payment.setAmount(paymentDetails.getAmount());
        payment.setCurrency(paymentDetails.getCurrency());
        payment.setStatus(paymentDetails.getStatus());
        payment.setPaymentMethod(paymentDetails.getPaymentMethod());
        
        return paymentRepository.update(payment);
    }
    
    public void deletePayment(Long id) {
        paymentRepository.deleteById(id);
    }
}
