# Postman Testing Guide

## Base URL

Use:

```text
http://localhost:8080
```

## Common Response Shapes

Success:

```json
{ "data": { }, "message": "ok" }
```

Error:

```json
{ "error": "mensaje descriptivo", "code": "ERROR_CODE" }
```

## Recommended Order

1. `POST /customers`
2. `POST /products`
3. `POST /ingredients`
4. `POST /recipes`
5. `POST /recipes/{id}/items`
6. `GET /recipes/product/{product_id}`
7. `GET /recipes/{id}/cost`
8. `POST /orders`
9. `POST /orders/{id}/items`
10. `PATCH /orders/{id}/status`
11. `POST /payments`
12. `GET /orders/{id}/payments`
13. `DELETE` endpoints

## Customers

### POST `/customers`

```json
{
  "full_name": "Cliente Demo",
  "phone": "3001234567",
  "email": "cliente.demo@example.com",
  "customer_type": "PERSON"
}
```

### GET `/customers`

No body.

### GET `/customers/{id}`

Use the `id` returned by create.

### PUT `/customers/{id}`

```json
{
  "full_name": "Cliente Demo Editado",
  "phone": "3009990000",
  "email": "cliente.demo.editado@example.com",
  "customer_type": "COMPANY"
}
```

### DELETE `/customers/{id}`

No body.

After delete:

- `GET /customers` should no longer list it
- `GET /customers/{id}` should return an inactive/not-found style error depending on flow

## Products

### POST `/products`

```json
{
  "name": "Almuerzo Ejecutivo",
  "description": "Arroz, proteina y ensalada",
  "product_type": "LUNCH",
  "base_price": "18000",
  "cost_price": "9000"
}
```

### GET `/products`

Optional query:

```text
/products?product_type=LUNCH
```

### PUT `/products/{id}`

```json
{
  "name": "Almuerzo Premium",
  "description": "Actualizado",
  "product_type": "LUNCH",
  "base_price": "22000",
  "cost_price": "11000"
}
```

### DELETE `/products/{id}`

No body.

## Ingredients

### POST `/ingredients`

```json
{
  "name": "Arroz",
  "unit": "KG",
  "average_cost": "4500",
  "stock": "20",
  "minimum_stock": "5"
}
```

### PUT `/ingredients/{id}`

```json
{
  "name": "Arroz Diana",
  "unit": "KG",
  "average_cost": "4800",
  "stock": "18",
  "minimum_stock": "4"
}
```

### DELETE `/ingredients/{id}`

No body.

## Recipes

### POST `/recipes`

```json
{
  "product_id": "{{product_id}}",
  "name": "Receta Almuerzo Ejecutivo",
  "portions": 1
}
```

### POST `/recipes/{id}/items`

```json
{
  "ingredient_id": "{{ingredient_id}}",
  "quantity": "0.25",
  "unit": "KG"
}
```

### GET `/recipes/product/{product_id}`

No body.

### GET `/recipes/{id}/cost`

No body.

## Orders

### POST `/orders`

```json
{
  "customer_id": "{{customer_id}}",
  "discount": "1000"
}
```

### GET `/orders`

Optional query:

```text
/orders?customer_id={{customer_id}}
```

### GET `/orders/{id}`

No body.

### POST `/orders/{id}/items`

```json
{
  "product_id": "{{product_id}}",
  "quantity": 2
}
```

### PATCH `/orders/{id}/status`

Example valid transitions:

- `PENDING` -> `CONFIRMED`
- `CONFIRMED` -> `IN_PREPARATION`
- `IN_PREPARATION` -> `READY`
- `READY` -> `DELIVERED`

Body:

```json
{
  "status": "CONFIRMED"
}
```

## Payments

### POST `/payments`

```json
{
  "order_id": "{{order_id}}",
  "amount": "35000",
  "method": "CASH",
  "status": "PENDING"
}
```

### GET `/orders/{id}/payments`

No body.

### PATCH `/payments/{id}/status`

```json
{
  "status": "PAID"
}
```

## Suggested Postman Variables

Create these collection or environment variables as you go:

- `base_url`
- `customer_id`
- `product_id`
- `ingredient_id`
- `recipe_id`
- `order_id`
- `payment_id`

## Common Error Checks

- Invalid UUID should return `400`
- Invalid enum/value object should return `400`
- Invalid status transition should return `422`
- Missing resource should return `404`

## Quick cURL Example

```bash
curl -X POST http://localhost:8080/customers \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Cliente Demo",
    "phone": "3001234567",
    "email": "cliente.demo@example.com",
    "customer_type": "PERSON"
  }'
```
