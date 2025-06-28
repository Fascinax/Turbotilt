package com.turbotilt.demo.order.controller;

import com.turbotilt.demo.order.model.Order;
import com.turbotilt.demo.order.service.OrderService;
import io.micronaut.http.HttpResponse;
import io.micronaut.http.annotation.*;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.Parameter;
import io.swagger.v3.oas.annotations.media.Content;
import io.swagger.v3.oas.annotations.media.Schema;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;

import java.util.List;

@Controller("/api/orders")
@RequiredArgsConstructor
@Tag(name = "Order", description = "Order management API")
public class OrderController {
    
    private final OrderService orderService;
    
    @Get
    @Operation(summary = "Get all orders", description = "Returns a list of all orders")
    public HttpResponse<List<Order>> getAllOrders() {
        return HttpResponse.ok(orderService.findAllOrders());
    }
    
    @Get("/{id}")
    @Operation(summary = "Get order by ID", description = "Returns a specific order by its ID")
    public HttpResponse<Order> getOrderById(@PathVariable Long id) {
        return orderService.findOrderById(id)
                .map(HttpResponse::ok)
                .orElse(HttpResponse.notFound());
    }
    
    @Get("/user/{userId}")
    @Operation(summary = "Get orders by user ID", description = "Returns all orders for a specific user")
    public HttpResponse<List<Order>> getOrdersByUserId(@PathVariable Long userId) {
        List<Order> orders = orderService.findOrdersByUserId(userId);
        return HttpResponse.ok(orders);
    }
    
    @Get("/status/{status}")
    @Operation(summary = "Get orders by status", description = "Returns all orders with a specific status")
    public HttpResponse<List<Order>> getOrdersByStatus(@PathVariable String status) {
        List<Order> orders = orderService.findOrdersByStatus(status);
        return HttpResponse.ok(orders);
    }
    
    @Post
    @Operation(
        summary = "Create a new order", 
        description = "Creates a new order and sends a message to RabbitMQ",
        responses = {
            @ApiResponse(
                responseCode = "201", 
                description = "Order created successfully",
                content = @Content(schema = @Schema(implementation = Order.class))
            )
        }
    )
    public HttpResponse<Order> createOrder(@Body Order order) {
        Order createdOrder = orderService.createOrder(order);
        return HttpResponse.created(createdOrder);
    }
    
    @Put("/{id}/status")
    @Operation(summary = "Update order status", description = "Updates the status of an existing order")
    public HttpResponse<Order> updateOrderStatus(
            @PathVariable Long id,
            @Parameter(description = "New status value") @QueryValue String status) {
        return orderService.updateOrderStatus(id, status)
                .map(HttpResponse::ok)
                .orElse(HttpResponse.notFound());
    }
    
    @Delete("/{id}")
    @Operation(summary = "Delete an order", description = "Deletes an order by its ID")
    public HttpResponse<?> deleteOrder(@PathVariable Long id) {
        orderService.deleteOrder(id);
        return HttpResponse.noContent();
    }
}
