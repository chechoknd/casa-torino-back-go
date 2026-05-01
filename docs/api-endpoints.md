# API Endpoints

Base URL local:

```text
http://localhost:8080
```

## Health

- `GET http://localhost:8080/health`
- `POST http://localhost:8080/health`

## Customers

- `POST http://localhost:8080/customers`
- `GET http://localhost:8080/customers`
- `GET http://localhost:8080/customers/{id}`
- `PUT http://localhost:8080/customers/{id}`
- `DELETE http://localhost:8080/customers/{id}`

## Products

- `POST http://localhost:8080/products`
- `GET http://localhost:8080/products`
- `GET http://localhost:8080/products/{id}`
- `PUT http://localhost:8080/products/{id}`
- `DELETE http://localhost:8080/products/{id}`

## Ingredients

- `POST http://localhost:8080/ingredients`
- `GET http://localhost:8080/ingredients`
- `GET http://localhost:8080/ingredients/{id}`
- `PUT http://localhost:8080/ingredients/{id}`
- `DELETE http://localhost:8080/ingredients/{id}`

## Recipes

- `POST http://localhost:8080/recipes`
- `POST http://localhost:8080/recipes/{id}/items`
- `GET http://localhost:8080/recipes/product/{product_id}`
- `GET http://localhost:8080/recipes/{id}/cost`

## Orders

- `POST http://localhost:8080/orders`
- `GET http://localhost:8080/orders`
- `GET http://localhost:8080/orders/{id}`
- `POST http://localhost:8080/orders/{id}/items`
- `PATCH http://localhost:8080/orders/{id}/status`
- `GET http://localhost:8080/orders/{id}/payments`

## Payments

- `POST http://localhost:8080/payments`
- `PATCH http://localhost:8080/payments/{id}/status`

## Query Params Disponibles

### Products

- `GET http://localhost:8080/products?product_type=LUNCH`

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
