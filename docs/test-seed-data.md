# Test Seed Data

Datos de prueba para poblar la base con escenarios más realistas de Casa Torino.

## Orden recomendado de carga

1. Clientes
2. Productos
3. Ingredientes
4. Recetas
5. Ítems de receta
6. Pedidos
7. Ítems de pedido

## Notas

- Todos los `customer_type` usan `PERSON` o `COMPANY`
- Todos los `product_type` usan valores válidos del dominio
- Todos los valores monetarios van como string
- Las recetas deben crearse después de los productos
- Los pedidos deben crearse después de los clientes
- Los ítems de pedido deben crearse después de productos y pedidos

## Clientes

Usa `POST /customers`

```json
[
  {
    "alias": "customer_01",
    "full_name": "Laura Mendoza",
    "phone": "3001001001",
    "email": "laura.mendoza@example.com",
    "customer_type": "PERSON"
  },
  {
    "alias": "customer_02",
    "full_name": "Carlos Ramirez",
    "phone": "3001001002",
    "email": "carlos.ramirez@example.com",
    "customer_type": "PERSON"
  },
  {
    "alias": "customer_03",
    "full_name": "Sofia Herrera",
    "phone": "3001001003",
    "email": "sofia.herrera@example.com",
    "customer_type": "PERSON"
  },
  {
    "alias": "customer_04",
    "full_name": "Mateo Rojas",
    "phone": "3001001004",
    "email": "mateo.rojas@example.com",
    "customer_type": "PERSON"
  },
  {
    "alias": "customer_05",
    "full_name": "Valentina Castro",
    "phone": "3001001005",
    "email": "valentina.castro@example.com",
    "customer_type": "PERSON"
  },
  {
    "alias": "customer_06",
    "full_name": "Daniela Torres",
    "phone": "3001001006",
    "email": "daniela.torres@example.com",
    "customer_type": "PERSON"
  },
  {
    "alias": "customer_07",
    "full_name": "Andres Molina",
    "phone": "3001001007",
    "email": "andres.molina@example.com",
    "customer_type": "PERSON"
  },
  {
    "alias": "customer_08",
    "full_name": "Camila Vega",
    "phone": "3001001008",
    "email": "camila.vega@example.com",
    "customer_type": "PERSON"
  },
  {
    "alias": "customer_09",
    "full_name": "Oficinas Nova SAS",
    "phone": "3001001009",
    "email": "compras@oficinasnova.com",
    "customer_type": "COMPANY"
  },
  {
    "alias": "customer_10",
    "full_name": "Clinica Santa Maria",
    "phone": "3001001010",
    "email": "eventos@clinicasantamaria.com",
    "customer_type": "COMPANY"
  }
]
```

## Productos

Usa `POST /products`

```json
[
  {
    "alias": "product_01",
    "name": "Almuerzo Ejecutivo Pollo",
    "description": "Arroz, pechuga a la plancha, ensalada y bebida",
    "product_type": "LUNCH",
    "base_price": "18500",
    "cost_price": "9200"
  },
  {
    "alias": "product_02",
    "name": "Almuerzo Ejecutivo Carne",
    "description": "Arroz, carne asada, ensalada y bebida",
    "product_type": "LUNCH",
    "base_price": "19500",
    "cost_price": "9800"
  },
  {
    "alias": "product_03",
    "name": "Jugo Natural de Mango",
    "description": "Jugo natural sin azucar adicionada",
    "product_type": "JUICE",
    "base_price": "7000",
    "cost_price": "2800"
  },
  {
    "alias": "product_04",
    "name": "Jugo Natural de Fresa",
    "description": "Jugo natural con leche opcional",
    "product_type": "JUICE",
    "base_price": "7500",
    "cost_price": "3100"
  },
  {
    "alias": "product_05",
    "name": "Torta de Zanahoria",
    "description": "Porcion individual de torta casera",
    "product_type": "CAKE",
    "base_price": "8500",
    "cost_price": "3600"
  },
  {
    "alias": "product_06",
    "name": "Plan Fit Semanal",
    "description": "Cinco almuerzos balanceados para semana laboral",
    "product_type": "PLAN",
    "base_price": "89000",
    "cost_price": "47000"
  },
  {
    "alias": "product_07",
    "name": "Bowl Saludable",
    "description": "Proteina, quinoa, vegetales y salsa de la casa",
    "product_type": "LUNCH",
    "base_price": "21000",
    "cost_price": "10800"
  },
  {
    "alias": "product_08",
    "name": "Lasaña Familiar Empacada",
    "description": "Lasaña lista para calentar empacada al vacio",
    "product_type": "VACUUM_PACKED",
    "base_price": "32000",
    "cost_price": "16500"
  },
  {
    "alias": "product_09",
    "name": "Coffee Break Empresarial",
    "description": "Combo para reuniones de 10 personas",
    "product_type": "EVENT",
    "base_price": "120000",
    "cost_price": "65000"
  },
  {
    "alias": "product_10",
    "name": "Torta de Chocolate Entera",
    "description": "Torta para celebraciones de 12 porciones",
    "product_type": "CAKE",
    "base_price": "68000",
    "cost_price": "29000"
  }
]
```

## Ingredientes

Usa `POST /ingredients`

```json
[
  {
    "alias": "ingredient_01",
    "name": "Arroz Blanco",
    "unit": "KG",
    "average_cost": "4500",
    "stock": "25",
    "minimum_stock": "5"
  },
  {
    "alias": "ingredient_02",
    "name": "Pechuga de Pollo",
    "unit": "KG",
    "average_cost": "17000",
    "stock": "18",
    "minimum_stock": "4"
  },
  {
    "alias": "ingredient_03",
    "name": "Carne de Res",
    "unit": "KG",
    "average_cost": "24000",
    "stock": "14",
    "minimum_stock": "3"
  },
  {
    "alias": "ingredient_04",
    "name": "Lechuga",
    "unit": "KG",
    "average_cost": "6000",
    "stock": "8",
    "minimum_stock": "2"
  },
  {
    "alias": "ingredient_05",
    "name": "Tomate",
    "unit": "KG",
    "average_cost": "5200",
    "stock": "10",
    "minimum_stock": "2"
  },
  {
    "alias": "ingredient_06",
    "name": "Mango",
    "unit": "KG",
    "average_cost": "6800",
    "stock": "12",
    "minimum_stock": "3"
  },
  {
    "alias": "ingredient_07",
    "name": "Fresa",
    "unit": "KG",
    "average_cost": "9000",
    "stock": "9",
    "minimum_stock": "2"
  },
  {
    "alias": "ingredient_08",
    "name": "Leche",
    "unit": "L",
    "average_cost": "4200",
    "stock": "20",
    "minimum_stock": "5"
  },
  {
    "alias": "ingredient_09",
    "name": "Zanahoria",
    "unit": "KG",
    "average_cost": "3500",
    "stock": "11",
    "minimum_stock": "3"
  },
  {
    "alias": "ingredient_10",
    "name": "Huevos",
    "unit": "UNIT",
    "average_cost": "600",
    "stock": "120",
    "minimum_stock": "24"
  },
  {
    "alias": "ingredient_11",
    "name": "Harina de Trigo",
    "unit": "KG",
    "average_cost": "3200",
    "stock": "15",
    "minimum_stock": "4"
  },
  {
    "alias": "ingredient_12",
    "name": "Chocolate en Polvo",
    "unit": "KG",
    "average_cost": "14000",
    "stock": "6",
    "minimum_stock": "2"
  },
  {
    "alias": "ingredient_13",
    "name": "Quinoa",
    "unit": "KG",
    "average_cost": "18000",
    "stock": "7",
    "minimum_stock": "2"
  },
  {
    "alias": "ingredient_14",
    "name": "Queso Mozzarella",
    "unit": "KG",
    "average_cost": "23000",
    "stock": "8",
    "minimum_stock": "2"
  },
  {
    "alias": "ingredient_15",
    "name": "Pasta para Lasaña",
    "unit": "KG",
    "average_cost": "7800",
    "stock": "10",
    "minimum_stock": "3"
  }
]
```

## Recetas

Usa `POST /recipes`

```json
[
  {
    "alias": "recipe_01",
    "product_alias": "product_01",
    "name": "Receta Almuerzo Ejecutivo Pollo",
    "portions": 1
  },
  {
    "alias": "recipe_02",
    "product_alias": "product_02",
    "name": "Receta Almuerzo Ejecutivo Carne",
    "portions": 1
  },
  {
    "alias": "recipe_03",
    "product_alias": "product_03",
    "name": "Receta Jugo de Mango",
    "portions": 1
  },
  {
    "alias": "recipe_04",
    "product_alias": "product_04",
    "name": "Receta Jugo de Fresa",
    "portions": 1
  },
  {
    "alias": "recipe_05",
    "product_alias": "product_05",
    "name": "Receta Torta de Zanahoria",
    "portions": 1
  },
  {
    "alias": "recipe_06",
    "product_alias": "product_06",
    "name": "Receta Plan Fit Semanal",
    "portions": 5
  },
  {
    "alias": "recipe_07",
    "product_alias": "product_07",
    "name": "Receta Bowl Saludable",
    "portions": 1
  },
  {
    "alias": "recipe_08",
    "product_alias": "product_08",
    "name": "Receta Lasaña Familiar Empacada",
    "portions": 4
  },
  {
    "alias": "recipe_09",
    "product_alias": "product_09",
    "name": "Receta Coffee Break Empresarial",
    "portions": 10
  },
  {
    "alias": "recipe_10",
    "product_alias": "product_10",
    "name": "Receta Torta de Chocolate Entera",
    "portions": 12
  }
]
```

## Ítems de receta

Usa `POST /recipes/{id}/items`

```json
[
  {
    "recipe_alias": "recipe_01",
    "ingredient_alias": "ingredient_01",
    "quantity": "0.18",
    "unit": "KG"
  },
  {
    "recipe_alias": "recipe_01",
    "ingredient_alias": "ingredient_02",
    "quantity": "0.20",
    "unit": "KG"
  },
  {
    "recipe_alias": "recipe_01",
    "ingredient_alias": "ingredient_04",
    "quantity": "0.05",
    "unit": "KG"
  },
  {
    "recipe_alias": "recipe_01",
    "ingredient_alias": "ingredient_05",
    "quantity": "0.04",
    "unit": "KG"
  },
  {
    "recipe_alias": "recipe_02",
    "ingredient_alias": "ingredient_01",
    "quantity": "0.18",
    "unit": "KG"
  },
  {
    "recipe_alias": "recipe_02",
    "ingredient_alias": "ingredient_03",
    "quantity": "0.20",
    "unit": "KG"
  },
  {
    "recipe_alias": "recipe_02",
    "ingredient_alias": "ingredient_04",
    "quantity": "0.05",
    "unit": "KG"
  },
  {
    "recipe_alias": "recipe_03",
    "ingredient_alias": "ingredient_06",
    "quantity": "0.30",
    "unit": "KG"
  },
  {
    "recipe_alias": "recipe_04",
    "ingredient_alias": "ingredient_07",
    "quantity": "0.25",
    "unit": "KG"
  },
  {
    "recipe_alias": "recipe_04",
    "ingredient_alias": "ingredient_08",
    "quantity": "0.25",
    "unit": "L"
  },
  {
    "recipe_alias": "recipe_05",
    "ingredient_alias": "ingredient_09",
    "quantity": "0.12",
    "unit": "KG"
  },
  {
    "recipe_alias": "recipe_05",
    "ingredient_alias": "ingredient_10",
    "quantity": "2",
    "unit": "UNIT"
  },
  {
    "recipe_alias": "recipe_05",
    "ingredient_alias": "ingredient_11",
    "quantity": "0.10",
    "unit": "KG"
  },
  {
    "recipe_alias": "recipe_06",
    "ingredient_alias": "ingredient_02",
    "quantity": "1.00",
    "unit": "KG"
  },
  {
    "recipe_alias": "recipe_06",
    "ingredient_alias": "ingredient_01",
    "quantity": "0.90",
    "unit": "KG"
  },
  {
    "recipe_alias": "recipe_07",
    "ingredient_alias": "ingredient_13",
    "quantity": "0.12",
    "unit": "KG"
  },
  {
    "recipe_alias": "recipe_07",
    "ingredient_alias": "ingredient_02",
    "quantity": "0.18",
    "unit": "KG"
  },
  {
    "recipe_alias": "recipe_08",
    "ingredient_alias": "ingredient_14",
    "quantity": "0.35",
    "unit": "KG"
  },
  {
    "recipe_alias": "recipe_08",
    "ingredient_alias": "ingredient_15",
    "quantity": "0.40",
    "unit": "KG"
  },
  {
    "recipe_alias": "recipe_10",
    "ingredient_alias": "ingredient_12",
    "quantity": "0.20",
    "unit": "KG"
  }
]
```

## Pedidos

Usa `POST /orders`

```json
[
  {
    "alias": "order_01",
    "customer_alias": "customer_01",
    "discount": "0"
  },
  {
    "alias": "order_02",
    "customer_alias": "customer_02",
    "discount": "1000"
  },
  {
    "alias": "order_03",
    "customer_alias": "customer_09",
    "discount": "5000"
  },
  {
    "alias": "order_04",
    "customer_alias": "customer_05",
    "discount": "0"
  },
  {
    "alias": "order_05",
    "customer_alias": "customer_10",
    "discount": "8000"
  }
]
```

## Ítems de pedido

Usa `POST /orders/{id}/items`

```json
[
  {
    "order_alias": "order_01",
    "product_alias": "product_01",
    "quantity": 2
  },
  {
    "order_alias": "order_01",
    "product_alias": "product_03",
    "quantity": 2
  },
  {
    "order_alias": "order_02",
    "product_alias": "product_02",
    "quantity": 1
  },
  {
    "order_alias": "order_02",
    "product_alias": "product_05",
    "quantity": 1
  },
  {
    "order_alias": "order_03",
    "product_alias": "product_09",
    "quantity": 1
  },
  {
    "order_alias": "order_04",
    "product_alias": "product_07",
    "quantity": 2
  },
  {
    "order_alias": "order_04",
    "product_alias": "product_04",
    "quantity": 2
  },
  {
    "order_alias": "order_05",
    "product_alias": "product_10",
    "quantity": 1
  },
  {
    "order_alias": "order_05",
    "product_alias": "product_08",
    "quantity": 2
  }
]
```

## Estados sugeridos para pedidos

Después de crear los pedidos, puedes probar transiciones con `PATCH /orders/{id}/status`:

```json
[
  { "order_alias": "order_01", "status": "CONFIRMED" },
  { "order_alias": "order_01", "status": "IN_PREPARATION" },
  { "order_alias": "order_02", "status": "CONFIRMED" },
  { "order_alias": "order_03", "status": "CONFIRMED" },
  { "order_alias": "order_04", "status": "CANCELLED" }
]
```

## Pagos sugeridos

Usa `POST /payments`

```json
[
  {
    "order_alias": "order_01",
    "amount": "51000",
    "method": "TRANSFER",
    "status": "PAID"
  },
  {
    "order_alias": "order_02",
    "amount": "27000",
    "method": "NEQUI",
    "status": "PAID"
  },
  {
    "order_alias": "order_03",
    "amount": "115000",
    "method": "CARD",
    "status": "PENDING"
  },
  {
    "order_alias": "order_04",
    "amount": "57000",
    "method": "CASH",
    "status": "PAID"
  },
  {
    "order_alias": "order_05",
    "amount": "124000",
    "method": "TRANSFER",
    "status": "PARTIAL"
  }
]
```

## Recomendación práctica

Si vas a cargar esto en Postman:

1. Crea los clientes y guarda cada `id` en variables
2. Crea los productos y guarda cada `id`
3. Crea los ingredientes y guarda cada `id`
4. Crea las recetas y guarda cada `id`
5. Inserta ítems de receta usando los `id` ya guardados
6. Crea los pedidos
7. Inserta ítems de pedido
8. Registra pagos

Puedes complementar esta guía con:

- [postman-testing.md](/home/torodev/Documentos/workspace/backend/casa-torino-back-go/casa-torino-back-go/docs/postman-testing.md)
