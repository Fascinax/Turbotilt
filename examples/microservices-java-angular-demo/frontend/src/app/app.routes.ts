import { Routes } from '@angular/router';

export const routes: Routes = [
  { 
    path: 'products', 
    loadComponent: () => import('./features/products/product-list/product-list.component')
      .then(m => m.ProductListComponent) 
  },
  { 
    path: 'products/:id', 
    loadComponent: () => import('./features/products/product-details/product-details.component')
      .then(m => m.ProductDetailsComponent) 
  },
  { 
    path: 'users', 
    loadComponent: () => import('./features/users/user-list/user-list.component')
      .then(m => m.UserListComponent) 
  },
  { 
    path: 'users/:id', 
    loadComponent: () => import('./features/users/user-details/user-details.component')
      .then(m => m.UserDetailsComponent) 
  },
  { 
    path: 'orders', 
    loadComponent: () => import('./features/orders/order-list/order-list.component')
      .then(m => m.OrderListComponent) 
  },
  { 
    path: 'orders/:id', 
    loadComponent: () => import('./features/orders/order-details/order-details.component')
      .then(m => m.OrderDetailsComponent) 
  },
  { path: '', redirectTo: '/products', pathMatch: 'full' },
];
