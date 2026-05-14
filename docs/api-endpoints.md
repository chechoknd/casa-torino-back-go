# API Endpoints

Base URL local:

```text
http://localhost:8080
```

## Customers

Requiere JWT con rol `ADMIN`.

- `POST http://localhost:8080/customers`
- `GET http://localhost:8080/customers`
- `GET http://localhost:8080/customers/{id}`
- `PUT http://localhost:8080/customers/{id}`
- `DELETE http://localhost:8080/customers/{id}`

## Products

Rutas admin. Requieren JWT con rol `ADMIN`.

- `POST http://localhost:8080/products`
- `GET http://localhost:8080/products`
- `GET http://localhost:8080/products/{id}`
- `PUT http://localhost:8080/products/{id}`
- `DELETE http://localhost:8080/products/{id}`

Campos adicionales para crear/editar producto:

- `image_url`
- `is_public`

## Public Catalog

No requiere JWT. Pensado para guest mode y catalogo cliente.

- `GET http://localhost:8080/public/products`
- `GET http://localhost:8080/public/products/{id}`
- `GET http://localhost:8080/public/product-categories`

`GET /public/products` soporta:

- `product_type`

La respuesta publica de productos no expone `cost_price`.

## Customer Panel

Requiere JWT con rol `CUSTOMER`.

- `GET http://localhost:8080/customer/profile`
- `GET http://localhost:8080/customer/orders`

El customer se resuelve por coincidencia entre el email del usuario autenticado y el email del customer existente.

## Ingredients

Requiere JWT con rol `ADMIN`.

- `POST http://localhost:8080/ingredients`
- `GET http://localhost:8080/ingredients`
- `GET http://localhost:8080/ingredients/{id}`
- `PUT http://localhost:8080/ingredients/{id}`
- `DELETE http://localhost:8080/ingredients/{id}`

## Recipes

Requiere JWT con rol `ADMIN`.

- `POST http://localhost:8080/recipes`
- `GET http://localhost:8080/recipes`
- `POST http://localhost:8080/recipes/{id}/items`
- `GET http://localhost:8080/recipes/product/{product_id}`
- `GET http://localhost:8080/recipes/{id}/cost`

## Orders

Requiere JWT con rol `ADMIN`.

- `POST http://localhost:8080/orders`
- `GET http://localhost:8080/orders`
- `GET http://localhost:8080/orders/{id}`
- `POST http://localhost:8080/orders/{id}/items`
- `PATCH http://localhost:8080/orders/{id}/status`
- `GET http://localhost:8080/orders/{id}/payments`

## Payments

Requiere JWT con rol `ADMIN`.

- `GET http://localhost:8080/payments`
- `POST http://localhost:8080/payments`
- `PATCH http://localhost:8080/payments/{id}/status`

## Campos relevantes de visualización

- `GET /recipes` y `GET /recipes/product/{product_id}` retornan `product_name`
- `GET /payments` y `GET /orders/{id}/payments` retornan `products[].product_name`
- `GET /orders` y `GET /orders/{id}` retornan `order_number`, `order_label`, `customer_name` y `items[].product_name`

## Query Params Disponibles

### Products

- `GET http://localhost:8080/products?product_type=LUNCH`
- `GET http://localhost:8080/public/products?product_type=LUNCH`

### Orders

- `GET http://localhost:8080/orders?customer_id={customer_id}`

## Formato de respuesta

Éxito:

```json
{ "data": { }, "message": "ok" }
```

Error:

```json
{ "error": "mensaje descriptivo", "code": "ERROR_CODE" }
```
