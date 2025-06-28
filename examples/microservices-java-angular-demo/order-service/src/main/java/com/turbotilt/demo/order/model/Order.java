package com.turbotilt.demo.order.model;

import io.micronaut.data.annotation.DateCreated;
import io.micronaut.data.annotation.GeneratedValue;
import io.micronaut.data.annotation.Id;
import io.micronaut.data.annotation.MappedEntity;
import io.micronaut.serde.annotation.Serdeable;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.List;

@Serdeable
@MappedEntity
@Data
@NoArgsConstructor
@AllArgsConstructor
public class Order {
    
    @Id
    @GeneratedValue
    private Long id;
    
    private Long userId;
    private String status;
    private BigDecimal totalAmount;
    
    @DateCreated
    private LocalDateTime createdAt;
    
    private List<OrderItem> items = new ArrayList<>();
}
