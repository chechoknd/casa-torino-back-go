#!/bin/bash
# Seed script for Casa Torino test data
# Usage: bash scripts/seed_data.sh
# Requires: API running at localhost:8080, jq installed
set -euo pipefail

API="http://localhost:8080"

declare -A CUSTOMER_IDS
declare -A PRODUCT_IDS
declare -A INGREDIENT_IDS
declare -A RECIPE_IDS
declare -A ORDER_IDS

echo "========================================"
echo " Seeding Casa Torino Test Data"
echo "========================================"

# ========== 1. CUSTOMERS ==========
echo ""
echo ">>> 1/7 Creating Customers..."
curl -s -X POST "$API/customers" \
  -H "Content-Type: application/json" \
  -d '{"full_name":"Laura Mendoza","phone":"3001001001","email":"laura.mendoza@example.com","customer_type":"PERSON"}' \
  | jq -r '.data.id' | read -r CUSTOMER_IDS["customer_01"]
CUSTOMER_IDS["customer_01"]=$(curl -s -X POST "$API/customers" -H "Content-Type: application/json" \
  -d '{"full_name":"Laura Mendoza","phone":"3001001001","email":"laura.mendoza@example.com","customer_type":"PERSON"}' | jq -r '.data.id')
echo "  customer_01 -> ${CUSTOMER_IDS["customer_01"]}"

CUSTOMER_IDS["customer_02"]=$(curl -s -X POST "$API/customers" -H "Content-Type: application/json" \
  -d '{"full_name":"Carlos Ramirez","phone":"3001001002","email":"carlos.ramirez@example.com","customer_type":"PERSON"}' | jq -r '.data.id')
echo "  customer_02 -> ${CUSTOMER_IDS["customer_02"]}"

CUSTOMER_IDS["customer_03"]=$(curl -s -X POST "$API/customers" -H "Content-Type: application/json" \
  -d '{"full_name":"Sofia Herrera","phone":"3001001003","email":"sofia.herrera@example.com","customer_type":"PERSON"}' | jq -r '.data.id')
echo "  customer_03 -> ${CUSTOMER_IDS["customer_03"]}"

CUSTOMER_IDS["customer_04"]=$(curl -s -X POST "$API/customers" -H "Content-Type: application/json" \
  -d '{"full_name":"Mateo Rojas","phone":"3001001004","email":"mateo.rojas@example.com","customer_type":"PERSON"}' | jq -r '.data.id')
echo "  customer_04 -> ${CUSTOMER_IDS["customer_04"]}"

CUSTOMER_IDS["customer_05"]=$(curl -s -X POST "$API/customers" -H "Content-Type: application/json" \
  -d '{"full_name":"Valentina Castro","phone":"3001001005","email":"valentina.castro@example.com","customer_type":"PERSON"}' | jq -r '.data.id')
echo "  customer_05 -> ${CUSTOMER_IDS["customer_05"]}"

CUSTOMER_IDS["customer_06"]=$(curl -s -X POST "$API/customers" -H "Content-Type: application/json" \
  -d '{"full_name":"Daniela Torres","phone":"3001001006","email":"daniela.torres@example.com","customer_type":"PERSON"}' | jq -r '.data.id')
echo "  customer_06 -> ${CUSTOMER_IDS["customer_06"]}"

CUSTOMER_IDS["customer_07"]=$(curl -s -X POST "$API/customers" -H "Content-Type: application/json" \
  -d '{"full_name":"Andres Molina","phone":"3001001007","email":"andres.molina@example.com","customer_type":"PERSON"}' | jq -r '.data.id')
echo "  customer_07 -> ${CUSTOMER_IDS["customer_07"]}"

CUSTOMER_IDS["customer_08"]=$(curl -s -X POST "$API/customers" -H "Content-Type: application/json" \
  -d '{"full_name":"Camila Vega","phone":"3001001008","email":"camila.vega@example.com","customer_type":"PERSON"}' | jq -r '.data.id')
echo "  customer_08 -> ${CUSTOMER_IDS["customer_08"]}"

CUSTOMER_IDS["customer_09"]=$(curl -s -X POST "$API/customers" -H "Content-Type: application/json" \
  -d '{"full_name":"Oficinas Nova SAS","phone":"3001001009","email":"compras@oficinasnova.com","customer_type":"COMPANY"}' | jq -r '.data.id')
echo "  customer_09 -> ${CUSTOMER_IDS["customer_09"]}"

CUSTOMER_IDS["customer_10"]=$(curl -s -X POST "$API/customers" -H "Content-Type: application/json" \
  -d '{"full_name":"Clinica Santa Maria","phone":"3001001010","email":"eventos@clinicasantamaria.com","customer_type":"COMPANY"}' | jq -r '.data.id')
echo "  customer_10 -> ${CUSTOMER_IDS["customer_10"]}"

# ========== 2. PRODUCTS ==========
echo ""
echo ">>> 2/7 Creating Products..."
PRODUCT_IDS["product_01"]=$(curl -s -X POST "$API/products" -H "Content-Type: application/json" \
  -d '{"name":"Almuerzo Ejecutivo Pollo","description":"Arroz, pechuga a la plancha, ensalada y bebida","product_type":"LUNCH","base_price":"18500","cost_price":"9200"}' | jq -r '.data.id')
echo "  product_01 -> ${PRODUCT_IDS["product_01"]}"

PRODUCT_IDS["product_02"]=$(curl -s -X POST "$API/products" -H "Content-Type: application/json" \
  -d '{"name":"Almuerzo Ejecutivo Carne","description":"Arroz, carne asada, ensalada y bebida","product_type":"LUNCH","base_price":"19500","cost_price":"9800"}' | jq -r '.data.id')
echo "  product_02 -> ${PRODUCT_IDS["product_02"]}"

PRODUCT_IDS["product_03"]=$(curl -s -X POST "$API/products" -H "Content-Type: application/json" \
  -d '{"name":"Jugo Natural de Mango","description":"Jugo natural sin azucar adicionada","product_type":"JUICE","base_price":"7000","cost_price":"2800"}' | jq -r '.data.id')
echo "  product_03 -> ${PRODUCT_IDS["product_03"]}"

PRODUCT_IDS["product_04"]=$(curl -s -X POST "$API/products" -H "Content-Type: application/json" \
  -d '{"name":"Jugo Natural de Fresa","description":"Jugo natural con leche opcional","product_type":"JUICE","base_price":"7500","cost_price":"3100"}' | jq -r '.data.id')
echo "  product_04 -> ${PRODUCT_IDS["product_04"]}"

PRODUCT_IDS["product_05"]=$(curl -s -X POST "$API/products" -H "Content-Type: application/json" \
  -d '{"name":"Torta de Zanahoria","description":"Porcion individual de torta casera","product_type":"CAKE","base_price":"8500","cost_price":"3600"}' | jq -r '.data.id')
echo "  product_05 -> ${PRODUCT_IDS["product_05"]}"

PRODUCT_IDS["product_06"]=$(curl -s -X POST "$API/products" -H "Content-Type: application/json" \
  -d '{"name":"Plan Fit Semanal","description":"Cinco almuerzos balanceados para semana laboral","product_type":"PLAN","base_price":"89000","cost_price":"47000"}' | jq -r '.data.id')
echo "  product_06 -> ${PRODUCT_IDS["product_06"]}"

PRODUCT_IDS["product_07"]=$(curl -s -X POST "$API/products" -H "Content-Type: application/json" \
  -d '{"name":"Bowl Saludable","description":"Proteina, quinoa, vegetales y salsa de la casa","product_type":"LUNCH","base_price":"21000","cost_price":"10800"}' | jq -r '.data.id')
echo "  product_07 -> ${PRODUCT_IDS["product_07"]}"

PRODUCT_IDS["product_08"]=$(curl -s -X POST "$API/products" -H "Content-Type: application/json" \
  -d '{"name":"Lasaña Familiar Empacada","description":"Lasaña lista para calentar empacada al vacio","product_type":"VACUUM_PACKED","base_price":"32000","cost_price":"16500"}' | jq -r '.data.id')
echo "  product_08 -> ${PRODUCT_IDS["product_08"]}"

PRODUCT_IDS["product_09"]=$(curl -s -X POST "$API/products" -H "Content-Type: application/json" \
  -d '{"name":"Coffee Break Empresarial","description":"Combo para reuniones de 10 personas","product_type":"EVENT","base_price":"120000","cost_price":"65000"}' | jq -r '.data.id')
echo "  product_09 -> ${PRODUCT_IDS["product_09"]}"

PRODUCT_IDS["product_10"]=$(curl -s -X POST "$API/products" -H "Content-Type: application/json" \
  -d '{"name":"Torta de Chocolate Entera","description":"Torta para celebraciones de 12 porciones","product_type":"CAKE","base_price":"68000","cost_price":"29000"}' | jq -r '.data.id')
echo "  product_10 -> ${PRODUCT_IDS["product_10"]}"

# ========== 3. INGREDIENTS ==========
echo ""
echo ">>> 3/7 Creating Ingredients..."
INGREDIENT_IDS["ingredient_01"]=$(curl -s -X POST "$API/ingredients" -H "Content-Type: application/json" \
  -d '{"name":"Arroz Blanco","unit":"KG","average_cost":"4500","stock":"25","minimum_stock":"5"}' | jq -r '.data.id')
echo "  ingredient_01 -> ${INGREDIENT_IDS["ingredient_01"]}"

INGREDIENT_IDS["ingredient_02"]=$(curl -s -X POST "$API/ingredients" -H "Content-Type: application/json" \
  -d '{"name":"Pechuga de Pollo","unit":"KG","average_cost":"17000","stock":"18","minimum_stock":"4"}' | jq -r '.data.id')
echo "  ingredient_02 -> ${INGREDIENT_IDS["ingredient_02"]}"

INGREDIENT_IDS["ingredient_03"]=$(curl -s -X POST "$API/ingredients" -H "Content-Type: application/json" \
  -d '{"name":"Carne de Res","unit":"KG","average_cost":"24000","stock":"14","minimum_stock":"3"}' | jq -r '.data.id')
echo "  ingredient_03 -> ${INGREDIENT_IDS["ingredient_03"]}"

INGREDIENT_IDS["ingredient_04"]=$(curl -s -X POST "$API/ingredients" -H "Content-Type: application/json" \
  -d '{"name":"Lechuga","unit":"KG","average_cost":"6000","stock":"8","minimum_stock":"2"}' | jq -r '.data.id')
echo "  ingredient_04 -> ${INGREDIENT_IDS["ingredient_04"]}"

INGREDIENT_IDS["ingredient_05"]=$(curl -s -X POST "$API/ingredients" -H "Content-Type: application/json" \
  -d '{"name":"Tomate","unit":"KG","average_cost":"5200","stock":"10","minimum_stock":"2"}' | jq -r '.data.id')
echo "  ingredient_05 -> ${INGREDIENT_IDS["ingredient_05"]}"

INGREDIENT_IDS["ingredient_06"]=$(curl -s -X POST "$API/ingredients" -H "Content-Type: application/json" \
  -d '{"name":"Mango","unit":"KG","average_cost":"6800","stock":"12","minimum_stock":"3"}' | jq -r '.data.id')
echo "  ingredient_06 -> ${INGREDIENT_IDS["ingredient_06"]}"

INGREDIENT_IDS["ingredient_07"]=$(curl -s -X POST "$API/ingredients" -H "Content-Type: application/json" \
  -d '{"name":"Fresa","unit":"KG","average_cost":"9000","stock":"9","minimum_stock":"2"}' | jq -r '.data.id')
echo "  ingredient_07 -> ${INGREDIENT_IDS["ingredient_07"]}"

INGREDIENT_IDS["ingredient_08"]=$(curl -s -X POST "$API/ingredients" -H "Content-Type: application/json" \
  -d '{"name":"Leche","unit":"L","average_cost":"4200","stock":"20","minimum_stock":"5"}' | jq -r '.data.id')
echo "  ingredient_08 -> ${INGREDIENT_IDS["ingredient_08"]}"

INGREDIENT_IDS["ingredient_09"]=$(curl -s -X POST "$API/ingredients" -H "Content-Type: application/json" \
  -d '{"name":"Zanahoria","unit":"KG","average_cost":"3500","stock":"11","minimum_stock":"3"}' | jq -r '.data.id')
echo "  ingredient_09 -> ${INGREDIENT_IDS["ingredient_09"]}"

INGREDIENT_IDS["ingredient_10"]=$(curl -s -X POST "$API/ingredients" -H "Content-Type: application/json" \
  -d '{"name":"Huevos","unit":"UNIT","average_cost":"600","stock":"120","minimum_stock":"24"}' | jq -r '.data.id')
echo "  ingredient_10 -> ${INGREDIENT_IDS["ingredient_10"]}"

INGREDIENT_IDS["ingredient_11"]=$(curl -s -X POST "$API/ingredients" -H "Content-Type: application/json" \
  -d '{"name":"Harina de Trigo","unit":"KG","average_cost":"3200","stock":"15","minimum_stock":"4"}' | jq -r '.data.id')
echo "  ingredient_11 -> ${INGREDIENT_IDS["ingredient_11"]}"

INGREDIENT_IDS["ingredient_12"]=$(curl -s -X POST "$API/ingredients" -H "Content-Type: application/json" \
  -d '{"name":"Chocolate en Polvo","unit":"KG","average_cost":"14000","stock":"6","minimum_stock":"2"}' | jq -r '.data.id')
echo "  ingredient_12 -> ${INGREDIENT_IDS["ingredient_12"]}"

INGREDIENT_IDS["ingredient_13"]=$(curl -s -X POST "$API/ingredients" -H "Content-Type: application/json" \
  -d '{"name":"Quinoa","unit":"KG","average_cost":"18000","stock":"7","minimum_stock":"2"}' | jq -r '.data.id')
echo "  ingredient_13 -> ${INGREDIENT_IDS["ingredient_13"]}"

INGREDIENT_IDS["ingredient_14"]=$(curl -s -X POST "$API/ingredients" -H "Content-Type: application/json" \
  -d '{"name":"Queso Mozzarella","unit":"KG","average_cost":"23000","stock":"8","minimum_stock":"2"}' | jq -r '.data.id')
echo "  ingredient_14 -> ${INGREDIENT_IDS["ingredient_14"]}"

INGREDIENT_IDS["ingredient_15"]=$(curl -s -X POST "$API/ingredients" -H "Content-Type: application/json" \
  -d '{"name":"Pasta para Lasaña","unit":"KG","average_cost":"7800","stock":"10","minimum_stock":"3"}' | jq -r '.data.id')
echo "  ingredient_15 -> ${INGREDIENT_IDS["ingredient_15"]}"

# ========== 4. RECIPES ==========
echo ""
echo ">>> 4/7 Creating Recipes..."
RECIPE_IDS["recipe_01"]=$(curl -s -X POST "$API/recipes" -H "Content-Type: application/json" \
  -d '{"product_id":"'"${PRODUCT_IDS["product_01"]}"'","name":"Receta Almuerzo Ejecutivo Pollo","portions":1}' | jq -r '.data.id')
echo "  recipe_01 -> ${RECIPE_IDS["recipe_01"]}"

RECIPE_IDS["recipe_02"]=$(curl -s -X POST "$API/recipes" -H "Content-Type: application/json" \
  -d '{"product_id":"'"${PRODUCT_IDS["product_02"]}"'","name":"Receta Almuerzo Ejecutivo Carne","portions":1}' | jq -r '.data.id')
echo "  recipe_02 -> ${RECIPE_IDS["recipe_02"]}"

RECIPE_IDS["recipe_03"]=$(curl -s -X POST "$API/recipes" -H "Content-Type: application/json" \
  -d '{"product_id":"'"${PRODUCT_IDS["product_03"]}"'","name":"Receta Jugo de Mango","portions":1}' | jq -r '.data.id')
echo "  recipe_03 -> ${RECIPE_IDS["recipe_03"]}"

RECIPE_IDS["recipe_04"]=$(curl -s -X POST "$API/recipes" -H "Content-Type: application/json" \
  -d '{"product_id":"'"${PRODUCT_IDS["product_04"]}"'","name":"Receta Jugo de Fresa","portions":1}' | jq -r '.data.id')
echo "  recipe_04 -> ${RECIPE_IDS["recipe_04"]}"

RECIPE_IDS["recipe_05"]=$(curl -s -X POST "$API/recipes" -H "Content-Type: application/json" \
  -d '{"product_id":"'"${PRODUCT_IDS["product_05"]}"'","name":"Receta Torta de Zanahoria","portions":1}' | jq -r '.data.id')
echo "  recipe_05 -> ${RECIPE_IDS["recipe_05"]}"

RECIPE_IDS["recipe_06"]=$(curl -s -X POST "$API/recipes" -H "Content-Type: application/json" \
  -d '{"product_id":"'"${PRODUCT_IDS["product_06"]}"'","name":"Receta Plan Fit Semanal","portions":5}' | jq -r '.data.id')
echo "  recipe_06 -> ${RECIPE_IDS["recipe_06"]}"

RECIPE_IDS["recipe_07"]=$(curl -s -X POST "$API/recipes" -H "Content-Type: application/json" \
  -d '{"product_id":"'"${PRODUCT_IDS["product_07"]}"'","name":"Receta Bowl Saludable","portions":1}' | jq -r '.data.id')
echo "  recipe_07 -> ${RECIPE_IDS["recipe_07"]}"

RECIPE_IDS["recipe_08"]=$(curl -s -X POST "$API/recipes" -H "Content-Type: application/json" \
  -d '{"product_id":"'"${PRODUCT_IDS["product_08"]}"'","name":"Receta Lasaña Familiar Empacada","portions":4}' | jq -r '.data.id')
echo "  recipe_08 -> ${RECIPE_IDS["recipe_08"]}"

RECIPE_IDS["recipe_09"]=$(curl -s -X POST "$API/recipes" -H "Content-Type: application/json" \
  -d '{"product_id":"'"${PRODUCT_IDS["product_09"]}"'","name":"Receta Coffee Break Empresarial","portions":10}' | jq -r '.data.id')
echo "  recipe_09 -> ${RECIPE_IDS["recipe_09"]}"

RECIPE_IDS["recipe_10"]=$(curl -s -X POST "$API/recipes" -H "Content-Type: application/json" \
  -d '{"product_id":"'"${PRODUCT_IDS["product_10"]}"'","name":"Receta Torta de Chocolate Entera","portions":12}' | jq -r '.data.id')
echo "  recipe_10 -> ${RECIPE_IDS["recipe_10"]}"

# ========== 5. RECIPE ITEMS ==========
echo ""
echo ">>> 5/7 Adding Recipe Items..."
curl -s -X POST "$API/recipes/${RECIPE_IDS["recipe_01"]}/items" -H "Content-Type: application/json" \
  -d '{"ingredient_id":"'"${INGREDIENT_IDS["ingredient_01"]}"'","quantity":"0.18","unit":"KG"}' | jq -r '.data.id' > /dev/null && echo "  recipe_01 + ingredient_01"
curl -s -X POST "$API/recipes/${RECIPE_IDS["recipe_01"]}/items" -H "Content-Type: application/json" \
  -d '{"ingredient_id":"'"${INGREDIENT_IDS["ingredient_02"]}"'","quantity":"0.20","unit":"KG"}' | jq -r '.data.id' > /dev/null && echo "  recipe_01 + ingredient_02"
curl -s -X POST "$API/recipes/${RECIPE_IDS["recipe_01"]}/items" -H "Content-Type: application/json" \
  -d '{"ingredient_id":"'"${INGREDIENT_IDS["ingredient_04"]}"'","quantity":"0.05","unit":"KG"}' | jq -r '.data.id' > /dev/null && echo "  recipe_01 + ingredient_04"
curl -s -X POST "$API/recipes/${RECIPE_IDS["recipe_01"]}/items" -H "Content-Type: application/json" \
  -d '{"ingredient_id":"'"${INGREDIENT_IDS["ingredient_05"]}"'","quantity":"0.04","unit":"KG"}' | jq -r '.data.id' > /dev/null && echo "  recipe_01 + ingredient_05"

curl -s -X POST "$API/recipes/${RECIPE_IDS["recipe_02"]}/items" -H "Content-Type: application/json" \
  -d '{"ingredient_id":"'"${INGREDIENT_IDS["ingredient_01"]}"'","quantity":"0.18","unit":"KG"}' | jq -r '.data.id' > /dev/null && echo "  recipe_02 + ingredient_01"
curl -s -X POST "$API/recipes/${RECIPE_IDS["recipe_02"]}/items" -H "Content-Type: application/json" \
  -d '{"ingredient_id":"'"${INGREDIENT_IDS["ingredient_03"]}"'","quantity":"0.20","unit":"KG"}' | jq -r '.data.id' > /dev/null && echo "  recipe_02 + ingredient_03"
curl -s -X POST "$API/recipes/${RECIPE_IDS["recipe_02"]}/items" -H "Content-Type: application/json" \
  -d '{"ingredient_id":"'"${INGREDIENT_IDS["ingredient_04"]}"'","quantity":"0.05","unit":"KG"}' | jq -r '.data.id' > /dev/null && echo "  recipe_02 + ingredient_04"

curl -s -X POST "$API/recipes/${RECIPE_IDS["recipe_03"]}/items" -H "Content-Type: application/json" \
  -d '{"ingredient_id":"'"${INGREDIENT_IDS["ingredient_06"]}"'","quantity":"0.30","unit":"KG"}' | jq -r '.data.id' > /dev/null && echo "  recipe_03 + ingredient_06"

curl -s -X POST "$API/recipes/${RECIPE_IDS["recipe_04"]}/items" -H "Content-Type: application/json" \
  -d '{"ingredient_id":"'"${INGREDIENT_IDS["ingredient_07"]}"'","quantity":"0.25","unit":"KG"}' | jq -r '.data.id' > /dev/null && echo "  recipe_04 + ingredient_07"
curl -s -X POST "$API/recipes/${RECIPE_IDS["recipe_04"]}/items" -H "Content-Type: application/json" \
  -d '{"ingredient_id":"'"${INGREDIENT_IDS["ingredient_08"]}"'","quantity":"0.25","unit":"L"}' | jq -r '.data.id' > /dev/null && echo "  recipe_04 + ingredient_08"

curl -s -X POST "$API/recipes/${RECIPE_IDS["recipe_05"]}/items" -H "Content-Type: application/json" \
  -d '{"ingredient_id":"'"${INGREDIENT_IDS["ingredient_09"]}"'","quantity":"0.12","unit":"KG"}' | jq -r '.data.id' > /dev/null && echo "  recipe_05 + ingredient_09"
curl -s -X POST "$API/recipes/${RECIPE_IDS["recipe_05"]}/items" -H "Content-Type: application/json" \
  -d '{"ingredient_id":"'"${INGREDIENT_IDS["ingredient_10"]}"'","quantity":"2","unit":"UNIT"}' | jq -r '.data.id' > /dev/null && echo "  recipe_05 + ingredient_10"
curl -s -X POST "$API/recipes/${RECIPE_IDS["recipe_05"]}/items" -H "Content-Type: application/json" \
  -d '{"ingredient_id":"'"${INGREDIENT_IDS["ingredient_11"]}"'","quantity":"0.10","unit":"KG"}' | jq -r '.data.id' > /dev/null && echo "  recipe_05 + ingredient_11"

curl -s -X POST "$API/recipes/${RECIPE_IDS["recipe_06"]}/items" -H "Content-Type: application/json" \
  -d '{"ingredient_id":"'"${INGREDIENT_IDS["ingredient_02"]}"'","quantity":"1.00","unit":"KG"}' | jq -r '.data.id' > /dev/null && echo "  recipe_06 + ingredient_02"
curl -s -X POST "$API/recipes/${RECIPE_IDS["recipe_06"]}/items" -H "Content-Type: application/json" \
  -d '{"ingredient_id":"'"${INGREDIENT_IDS["ingredient_01"]}"'","quantity":"0.90","unit":"KG"}' | jq -r '.data.id' > /dev/null && echo "  recipe_06 + ingredient_01"

curl -s -X POST "$API/recipes/${RECIPE_IDS["recipe_07"]}/items" -H "Content-Type: application/json" \
  -d '{"ingredient_id":"'"${INGREDIENT_IDS["ingredient_13"]}"'","quantity":"0.12","unit":"KG"}' | jq -r '.data.id' > /dev/null && echo "  recipe_07 + ingredient_13"
curl -s -X POST "$API/recipes/${RECIPE_IDS["recipe_07"]}/items" -H "Content-Type: application/json" \
  -d '{"ingredient_id":"'"${INGREDIENT_IDS["ingredient_02"]}"'","quantity":"0.18","unit":"KG"}' | jq -r '.data.id' > /dev/null && echo "  recipe_07 + ingredient_02"

curl -s -X POST "$API/recipes/${RECIPE_IDS["recipe_08"]}/items" -H "Content-Type: application/json" \
  -d '{"ingredient_id":"'"${INGREDIENT_IDS["ingredient_14"]}"'","quantity":"0.35","unit":"KG"}' | jq -r '.data.id' > /dev/null && echo "  recipe_08 + ingredient_14"
curl -s -X POST "$API/recipes/${RECIPE_IDS["recipe_08"]}/items" -H "Content-Type: application/json" \
  -d '{"ingredient_id":"'"${INGREDIENT_IDS["ingredient_15"]}"'","quantity":"0.40","unit":"KG"}' | jq -r '.data.id' > /dev/null && echo "  recipe_08 + ingredient_15"

curl -s -X POST "$API/recipes/${RECIPE_IDS["recipe_10"]}/items" -H "Content-Type: application/json" \
  -d '{"ingredient_id":"'"${INGREDIENT_IDS["ingredient_12"]}"'","quantity":"0.20","unit":"KG"}' | jq -r '.data.id' > /dev/null && echo "  recipe_10 + ingredient_12"

# ========== 6. ORDERS ==========
echo ""
echo ">>> 6/7 Creating Orders..."
ORDER_IDS["order_01"]=$(curl -s -X POST "$API/orders" -H "Content-Type: application/json" \
  -d '{"customer_id":"'"${CUSTOMER_IDS["customer_01"]}"'","discount":"0"}' | jq -r '.data.id')
echo "  order_01 -> ${ORDER_IDS["order_01"]}"

ORDER_IDS["order_02"]=$(curl -s -X POST "$API/orders" -H "Content-Type: application/json" \
  -d '{"customer_id":"'"${CUSTOMER_IDS["customer_02"]}"'","discount":"1000"}' | jq -r '.data.id')
echo "  order_02 -> ${ORDER_IDS["order_02"]}"

ORDER_IDS["order_03"]=$(curl -s -X POST "$API/orders" -H "Content-Type: application/json" \
  -d '{"customer_id":"'"${CUSTOMER_IDS["customer_09"]}"'","discount":"5000"}' | jq -r '.data.id')
echo "  order_03 -> ${ORDER_IDS["order_03"]}"

ORDER_IDS["order_04"]=$(curl -s -X POST "$API/orders" -H "Content-Type: application/json" \
  -d '{"customer_id":"'"${CUSTOMER_IDS["customer_05"]}"'","discount":"0"}' | jq -r '.data.id')
echo "  order_04 -> ${ORDER_IDS["order_04"]}"

ORDER_IDS["order_05"]=$(curl -s -X POST "$API/orders" -H "Content-Type: application/json" \
  -d '{"customer_id":"'"${CUSTOMER_IDS["customer_10"]}"'","discount":"8000"}' | jq -r '.data.id')
echo "  order_05 -> ${ORDER_IDS["order_05"]}"

# ========== 7. ORDER ITEMS ==========
echo ""
echo ">>> 7/7 Adding Order Items..."
curl -s -X POST "$API/orders/${ORDER_IDS["order_01"]}/items" -H "Content-Type: application/json" \
  -d '{"product_id":"'"${PRODUCT_IDS["product_01"]}"'","quantity":2}' | jq -r '.data.id' > /dev/null && echo "  order_01 + product_01"
curl -s -X POST "$API/orders/${ORDER_IDS["order_01"]}/items" -H "Content-Type: application/json" \
  -d '{"product_id":"'"${PRODUCT_IDS["product_03"]}"'","quantity":2}' | jq -r '.data.id' > /dev/null && echo "  order_01 + product_03"

curl -s -X POST "$API/orders/${ORDER_IDS["order_02"]}/items" -H "Content-Type: application/json" \
  -d '{"product_id":"'"${PRODUCT_IDS["product_02"]}"'","quantity":1}' | jq -r '.data.id' > /dev/null && echo "  order_02 + product_02"
curl -s -X POST "$API/orders/${ORDER_IDS["order_02"]}/items" -H "Content-Type: application/json" \
  -d '{"product_id":"'"${PRODUCT_IDS["product_05"]}"'","quantity":1}' | jq -r '.data.id' > /dev/null && echo "  order_02 + product_05"

curl -s -X POST "$API/orders/${ORDER_IDS["order_03"]}/items" -H "Content-Type: application/json" \
  -d '{"product_id":"'"${PRODUCT_IDS["product_09"]}"'","quantity":1}' | jq -r '.data.id' > /dev/null && echo "  order_03 + product_09"

curl -s -X POST "$API/orders/${ORDER_IDS["order_04"]}/items" -H "Content-Type: application/json" \
  -d '{"product_id":"'"${PRODUCT_IDS["product_07"]}"'","quantity":2}' | jq -r '.data.id' > /dev/null && echo "  order_04 + product_07"
curl -s -X POST "$API/orders/${ORDER_IDS["order_04"]}/items" -H "Content-Type: application/json" \
  -d '{"product_id":"'"${PRODUCT_IDS["product_04"]}"'","quantity":2}' | jq -r '.data.id' > /dev/null && echo "  order_04 + product_04"

curl -s -X POST "$API/orders/${ORDER_IDS["order_05"]}/items" -H "Content-Type: application/json" \
  -d '{"product_id":"'"${PRODUCT_IDS["product_10"]}"'","quantity":1}' | jq -r '.data.id' > /dev/null && echo "  order_05 + product_10"
curl -s -X POST "$API/orders/${ORDER_IDS["order_05"]}/items" -H "Content-Type: application/json" \
  -d '{"product_id":"'"${PRODUCT_IDS["product_08"]}"'","quantity":2}' | jq -r '.data.id' > /dev/null && echo "  order_05 + product_08"

echo ""
echo "========================================"
echo " Seed complete!"
echo "========================================"
