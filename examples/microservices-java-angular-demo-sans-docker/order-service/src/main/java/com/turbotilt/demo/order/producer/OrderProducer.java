package com.turbotilt.demo.order.producer;

import com.turbotilt.demo.order.model.Order;
import io.micronaut.rabbitmq.annotation.Binding;
import io.micronaut.rabbitmq.annotation.RabbitClient;

@RabbitClient(value = "order-exchange")
public interface OrderProducer {
    
    @Binding("order.created")
    void sendOrderCreatedEvent(Order order);
}
