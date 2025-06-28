package com.turbotilt.demo.product.resource;

import com.turbotilt.demo.product.model.Product;
import com.turbotilt.demo.product.service.ProductService;
import jakarta.inject.Inject;
import jakarta.ws.rs.*;
import jakarta.ws.rs.core.MediaType;
import jakarta.ws.rs.core.Response;
import org.eclipse.microprofile.openapi.annotations.Operation;
import org.eclipse.microprofile.openapi.annotations.tags.Tag;

import java.util.List;

@Path("/api/products")
@Produces(MediaType.APPLICATION_JSON)
@Consumes(MediaType.APPLICATION_JSON)
@Tag(name = "Product Resource", description = "Product management endpoints")
public class ProductResource {
    
    @Inject
    ProductService productService;
    
    @GET
    @Operation(summary = "Get all products")
    public Response getAllProducts() {
        List<Product> products = productService.findAllProducts();
        return Response.ok(products).build();
    }
    
    @GET
    @Path("/{id}")
    @Operation(summary = "Get a product by ID")
    public Response getProductById(@PathParam("id") Long id) {
        return productService.findProductById(id)
                .map(product -> Response.ok(product).build())
                .orElse(Response.status(Response.Status.NOT_FOUND).build());
    }
    
    @GET
    @Path("/category/{category}")
    @Operation(summary = "Get products by category")
    public Response getProductsByCategory(@PathParam("category") String category) {
        List<Product> products = productService.findProductsByCategory(category);
        return Response.ok(products).build();
    }
    
    @GET
    @Path("/search")
    @Operation(summary = "Search products by name")
    public Response searchProducts(@QueryParam("name") String name) {
        List<Product> products = productService.findProductsByName(name);
        return Response.ok(products).build();
    }
    
    @POST
    @Operation(summary = "Create a new product")
    public Response createProduct(Product product) {
        Product saved = productService.saveProduct(product);
        return Response.status(Response.Status.CREATED).entity(saved).build();
    }
    
    @PUT
    @Path("/{id}")
    @Operation(summary = "Update a product")
    public Response updateProduct(@PathParam("id") Long id, Product product) {
        return productService.findProductById(id)
                .map(existingProduct -> {
                    product.id = id;
                    Product updated = productService.saveProduct(product);
                    return Response.ok(updated).build();
                })
                .orElse(Response.status(Response.Status.NOT_FOUND).build());
    }
    
    @DELETE
    @Path("/{id}")
    @Operation(summary = "Delete a product")
    public Response deleteProduct(@PathParam("id") Long id) {
        boolean deleted = productService.deleteProduct(id);
        return deleted 
                ? Response.noContent().build() 
                : Response.status(Response.Status.NOT_FOUND).build();
    }
    
    @PATCH
    @Path("/{id}/stock")
    @Operation(summary = "Update product stock")
    public Response updateStock(@PathParam("id") Long id, @QueryParam("quantity") int quantity) {
        boolean updated = productService.updateStock(id, quantity);
        return updated 
                ? Response.ok().build() 
                : Response.status(Response.Status.NOT_FOUND).build();
    }
}
