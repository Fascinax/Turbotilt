package com.turbotilt.demo.order.model;

import io.micronaut.data.annotation.GeneratedValue;
import io.micronaut.data.annotation.Id;
import io.micronaut.data.annotation.MappedEntity;
import io.micronaut.data.annotation.Relation;
import io.micronaut.serde.annotation.Serdeable;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.math.BigDecimal;

@Serdeable
@MappedEntity
@Data
@NoArgsConstructor
@AllArgsConstructor
public class OrderItem {
    
    @Id
    @GeneratedValue
    private Long id;
    
    @Relation(Relation.Kind.MANY_TO_ONE)
    private Order order;
    
    private Long productId;
    private String productName;
    private int quantity;
    private BigDecimal price;
}
