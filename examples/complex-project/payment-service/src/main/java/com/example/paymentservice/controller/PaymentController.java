package com.example.paymentservice.controller;

import com.example.paymentservice.model.Payment;
import com.example.paymentservice.service.PaymentService;
import io.micronaut.http.HttpResponse;
import io.micronaut.http.annotation.*;
import lombok.RequiredArgsConstructor;

import java.util.List;

@Controller("/api/payments")
@RequiredArgsConstructor
public class PaymentController {
    
    private final PaymentService paymentService;
    
    @Get
    public HttpResponse<List<Payment>> getAllPayments() {
        return HttpResponse.ok(paymentService.findAllPayments());
    }
    
    @Get("/{id}")
    public HttpResponse<Payment> getPaymentById(@PathVariable Long id) {
        return paymentService.findPaymentById(id)
                .map(HttpResponse::ok)
                .orElse(HttpResponse.notFound());
    }
    
    @Get("/user/{userId}")
    public HttpResponse<List<Payment>> getPaymentsByUserId(@PathVariable Long userId) {
        return HttpResponse.ok(paymentService.findPaymentsByUserId(userId));
    }
    
    @Post
    public HttpResponse<Payment> createPayment(@Body Payment payment) {
        return HttpResponse.created(paymentService.processPayment(payment));
    }
    
    @Put("/{id}")
    public HttpResponse<Payment> updatePayment(@PathVariable Long id, @Body Payment payment) {
        return HttpResponse.ok(paymentService.updatePayment(id, payment));
    }
    
    @Delete("/{id}")
    public HttpResponse<Void> deletePayment(@PathVariable Long id) {
        paymentService.deletePayment(id);
        return HttpResponse.noContent();
    }
}
