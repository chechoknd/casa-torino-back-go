#!/usr/bin/env bash

set -euo pipefail

BASE_URL="${1:-http://localhost:8080}"

declare -A CUSTOMERS
declare -A PRODUCTS
declare -A INGREDIENTS
declare -A RECIPES
declare -A ORDERS
declare -A PAYMENTS

post_json() {
  local path="$1"
  local payload="$2"
  curl -fsS -X POST "${BASE_URL}${path}" \
    -H "Content-Type: application/json" \
    -d "${payload}"
}

patch_json() {
  local path="$1"
  local payload="$2"
  curl -fsS -X PATCH "${BASE_URL}${path}" \
    -H "Content-Type: application/json" \
    -d "${payload}"
}

create_customer() {
  local alias="$1"
  local full_name="$2"
  local phone="$3"
  local email="$4"
  local customer_type="$5"
  local response

  response="$(post_json "/customers" "$(jq -nc \
    --arg full_name "$full_name" \
    --arg phone "$phone" \
    --arg email "$email" \
    --arg customer_type "$customer_type" \
    '{full_name:$full_name,phone:$phone,email:$email,customer_type:$customer_type}')" )"
  CUSTOMERS["$alias"]="$(jq -r '.data.id' <<<"$response")"
  echo "customer ${alias}: ${CUSTOMERS[$alias]}"
}

create_product() {
  local alias="$1"
  local name="$2"
  local description="$3"
  local product_type="$4"
  local base_price="$5"
  local cost_price="$6"
  local response

  response="$(post_json "/products" "$(jq -nc \
    --arg name "$name" \
    --arg description "$description" \
    --arg product_type "$product_type" \
    --arg base_price "$base_price" \
    --arg cost_price "$cost_price" \
    '{name:$name,description:$description,product_type:$product_type,base_price:$base_price,cost_price:$cost_price}')" )"
  PRODUCTS["$alias"]="$(jq -r '.data.id' <<<"$response")"
  echo "product ${alias}: ${PRODUCTS[$alias]}"
}

create_ingredient() {
  local alias="$1"
  local name="$2"
  local unit="$3"
  local average_cost="$4"
  local stock="$5"
  local minimum_stock="$6"
  local response

  response="$(post_json "/ingredients" "$(jq -nc \
    --arg name "$name" \
    --arg unit "$unit" \
    --arg average_cost "$average_cost" \
    --arg stock "$stock" \
    --arg minimum_stock "$minimum_stock" \
    '{name:$name,unit:$unit,average_cost:$average_cost,stock:$stock,minimum_stock:$minimum_stock}')" )"
  INGREDIENTS["$alias"]="$(jq -r '.data.id' <<<"$response")"
  echo "ingredient ${alias}: ${INGREDIENTS[$alias]}"
}

create_recipe() {
  local alias="$1"
  local product_alias="$2"
  local name="$3"
  local portions="$4"
  local response

  response="$(post_json "/recipes" "$(jq -nc \
    --arg product_id "${PRODUCTS[$product_alias]}" \
    --arg name "$name" \
    --argjson portions "$portions" \
    '{product_id:$product_id,name:$name,portions:$portions}')" )"
  RECIPES["$alias"]="$(jq -r '.data.id' <<<"$response")"
  echo "recipe ${alias}: ${RECIPES[$alias]}"
}

add_recipe_item() {
  local recipe_alias="$1"
  local ingredient_alias="$2"
  local quantity="$3"
  local unit="$4"

  post_json "/recipes/${RECIPES[$recipe_alias]}/items" "$(jq -nc \
    --arg ingredient_id "${INGREDIENTS[$ingredient_alias]}" \
    --arg quantity "$quantity" \
    --arg unit "$unit" \
    '{ingredient_id:$ingredient_id,quantity:$quantity,unit:$unit}')" >/dev/null
  echo "recipe item ${recipe_alias} -> ${ingredient_alias}"
}

create_order() {
  local alias="$1"
  local customer_alias="$2"
  local discount="$3"
  local response

  response="$(post_json "/orders" "$(jq -nc \
    --arg customer_id "${CUSTOMERS[$customer_alias]}" \
    --arg discount "$discount" \
    '{customer_id:$customer_id,discount:$discount}')" )"
  ORDERS["$alias"]="$(jq -r '.data.id' <<<"$response")"
  echo "order ${alias}: ${ORDERS[$alias]}"
}

add_order_item() {
  local order_alias="$1"
  local product_alias="$2"
  local quantity="$3"

  post_json "/orders/${ORDERS[$order_alias]}/items" "$(jq -nc \
    --arg product_id "${PRODUCTS[$product_alias]}" \
    --argjson quantity "$quantity" \
    '{product_id:$product_id,quantity:$quantity}')" >/dev/null
  echo "order item ${order_alias} -> ${product_alias}"
}

update_order_status() {
  local order_alias="$1"
  local status="$2"

  patch_json "/orders/${ORDERS[$order_alias]}/status" "$(jq -nc \
    --arg status "$status" \
    '{status:$status}')" >/dev/null
  echo "order status ${order_alias}: ${status}"
}

create_payment() {
  local alias="$1"
  local order_alias="$2"
  local amount="$3"
  local method="$4"
  local status="$5"
  local response

  response="$(post_json "/payments" "$(jq -nc \
    --arg order_id "${ORDERS[$order_alias]}" \
    --arg amount "$amount" \
    --arg method "$method" \
    --arg status "$status" \
    '{order_id:$order_id,amount:$amount,method:$method,status:$status}')" )"
  PAYMENTS["$alias"]="$(jq -r '.data.id' <<<"$response")"
  echo "payment ${alias}: ${PAYMENTS[$alias]}"
}

echo "Seeding API at ${BASE_URL}"

create_customer "customer_01" "Laura Mendoza" "3001001001" "laura.mendoza.seed@example.com" "PERSON"
create_customer "customer_02" "Carlos Ramirez" "3001001002" "carlos.ramirez.seed@example.com" "PERSON"
create_customer "customer_03" "Sofia Herrera" "3001001003" "sofia.herrera.seed@example.com" "PERSON"
create_customer "customer_04" "Mateo Rojas" "3001001004" "mateo.rojas.seed@example.com" "PERSON"
create_customer "customer_05" "Valentina Castro" "3001001005" "valentina.castro.seed@example.com" "PERSON"
create_customer "customer_06" "Daniela Torres" "3001001006" "daniela.torres.seed@example.com" "PERSON"
create_customer "customer_07" "Andres Molina" "3001001007" "andres.molina.seed@example.com" "PERSON"
create_customer "customer_08" "Camila Vega" "3001001008" "camila.vega.seed@example.com" "PERSON"
create_customer "customer_09" "Oficinas Nova SAS" "3001001009" "compras.seed@oficinasnova.com" "COMPANY"
create_customer "customer_10" "Clinica Santa Maria" "3001001010" "eventos.seed@clinicasantamaria.com" "COMPANY"

create_product "product_01" "Almuerzo Ejecutivo Pollo" "Arroz, pechuga a la plancha, ensalada y bebida" "LUNCH" "18500" "9200"
create_product "product_02" "Almuerzo Ejecutivo Carne" "Arroz, carne asada, ensalada y bebida" "LUNCH" "19500" "9800"
create_product "product_03" "Jugo Natural de Mango" "Jugo natural sin azucar adicionada" "JUICE" "7000" "2800"
create_product "product_04" "Jugo Natural de Fresa" "Jugo natural con leche opcional" "JUICE" "7500" "3100"
create_product "product_05" "Torta de Zanahoria" "Porcion individual de torta casera" "CAKE" "8500" "3600"
create_product "product_06" "Plan Fit Semanal" "Cinco almuerzos balanceados para semana laboral" "PLAN" "89000" "47000"
create_product "product_07" "Bowl Saludable" "Proteina, quinoa, vegetales y salsa de la casa" "LUNCH" "21000" "10800"
create_product "product_08" "Lasaña Familiar Empacada" "Lasaña lista para calentar empacada al vacio" "VACUUM_PACKED" "32000" "16500"
create_product "product_09" "Coffee Break Empresarial" "Combo para reuniones de 10 personas" "EVENT" "120000" "65000"
create_product "product_10" "Torta de Chocolate Entera" "Torta para celebraciones de 12 porciones" "CAKE" "68000" "29000"

create_ingredient "ingredient_01" "Arroz Blanco" "KG" "4500" "25" "5"
create_ingredient "ingredient_02" "Pechuga de Pollo" "KG" "17000" "18" "4"
create_ingredient "ingredient_03" "Carne de Res" "KG" "24000" "14" "3"
create_ingredient "ingredient_04" "Lechuga" "KG" "6000" "8" "2"
create_ingredient "ingredient_05" "Tomate" "KG" "5200" "10" "2"
create_ingredient "ingredient_06" "Mango" "KG" "6800" "12" "3"
create_ingredient "ingredient_07" "Fresa" "KG" "9000" "9" "2"
create_ingredient "ingredient_08" "Leche" "L" "4200" "20" "5"
create_ingredient "ingredient_09" "Zanahoria" "KG" "3500" "11" "3"
create_ingredient "ingredient_10" "Huevos" "UNIT" "600" "120" "24"
create_ingredient "ingredient_11" "Harina de Trigo" "KG" "3200" "15" "4"
create_ingredient "ingredient_12" "Chocolate en Polvo" "KG" "14000" "6" "2"
create_ingredient "ingredient_13" "Quinoa" "KG" "18000" "7" "2"
create_ingredient "ingredient_14" "Queso Mozzarella" "KG" "23000" "8" "2"
create_ingredient "ingredient_15" "Pasta para Lasaña" "KG" "7800" "10" "3"

create_recipe "recipe_01" "product_01" "Receta Almuerzo Ejecutivo Pollo" 1
create_recipe "recipe_02" "product_02" "Receta Almuerzo Ejecutivo Carne" 1
create_recipe "recipe_03" "product_03" "Receta Jugo de Mango" 1
create_recipe "recipe_04" "product_04" "Receta Jugo de Fresa" 1
create_recipe "recipe_05" "product_05" "Receta Torta de Zanahoria" 1
create_recipe "recipe_06" "product_06" "Receta Plan Fit Semanal" 5
create_recipe "recipe_07" "product_07" "Receta Bowl Saludable" 1
create_recipe "recipe_08" "product_08" "Receta Lasaña Familiar Empacada" 4
create_recipe "recipe_09" "product_09" "Receta Coffee Break Empresarial" 10
create_recipe "recipe_10" "product_10" "Receta Torta de Chocolate Entera" 12

add_recipe_item "recipe_01" "ingredient_01" "0.18" "KG"
add_recipe_item "recipe_01" "ingredient_02" "0.20" "KG"
add_recipe_item "recipe_01" "ingredient_04" "0.05" "KG"
add_recipe_item "recipe_01" "ingredient_05" "0.04" "KG"
add_recipe_item "recipe_02" "ingredient_01" "0.18" "KG"
add_recipe_item "recipe_02" "ingredient_03" "0.20" "KG"
add_recipe_item "recipe_02" "ingredient_04" "0.05" "KG"
add_recipe_item "recipe_03" "ingredient_06" "0.30" "KG"
add_recipe_item "recipe_04" "ingredient_07" "0.25" "KG"
add_recipe_item "recipe_04" "ingredient_08" "0.25" "L"
add_recipe_item "recipe_05" "ingredient_09" "0.12" "KG"
add_recipe_item "recipe_05" "ingredient_10" "2" "UNIT"
add_recipe_item "recipe_05" "ingredient_11" "0.10" "KG"
add_recipe_item "recipe_06" "ingredient_02" "1.00" "KG"
add_recipe_item "recipe_06" "ingredient_01" "0.90" "KG"
add_recipe_item "recipe_07" "ingredient_13" "0.12" "KG"
add_recipe_item "recipe_07" "ingredient_02" "0.18" "KG"
add_recipe_item "recipe_08" "ingredient_14" "0.35" "KG"
add_recipe_item "recipe_08" "ingredient_15" "0.40" "KG"
add_recipe_item "recipe_10" "ingredient_12" "0.20" "KG"

create_order "order_01" "customer_01" "0"
create_order "order_02" "customer_02" "1000"
create_order "order_03" "customer_09" "5000"
create_order "order_04" "customer_05" "0"
create_order "order_05" "customer_10" "8000"

add_order_item "order_01" "product_01" 2
add_order_item "order_01" "product_03" 2
add_order_item "order_02" "product_02" 1
add_order_item "order_02" "product_05" 1
add_order_item "order_03" "product_09" 1
add_order_item "order_04" "product_07" 2
add_order_item "order_04" "product_04" 2
add_order_item "order_05" "product_10" 1
add_order_item "order_05" "product_08" 2

update_order_status "order_01" "CONFIRMED"
update_order_status "order_01" "IN_PREPARATION"
update_order_status "order_02" "CONFIRMED"
update_order_status "order_03" "CONFIRMED"
update_order_status "order_04" "CANCELLED"

create_payment "payment_01" "order_01" "51000" "TRANSFER" "PAID"
create_payment "payment_02" "order_02" "27000" "NEQUI" "PAID"
create_payment "payment_03" "order_03" "115000" "CARD" "PENDING"
create_payment "payment_04" "order_04" "57000" "CASH" "PAID"
create_payment "payment_05" "order_05" "124000" "TRANSFER" "PARTIAL"

echo "Seed completed"
