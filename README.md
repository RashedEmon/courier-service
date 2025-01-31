# Project courier-service

This is a courier management system API application. Following are the features:
- Create Order.
- See Order List
- Cancelled order.
- Signup
- Login

## Project Dependency
- Go 1.23.2
- Gin
- Postgres 16
- Docker
- Docker-Compose


## Project Setup
- Clone git repository
```bash
git clone https://github.com/RashedEmon/courier-service.git
```
- Run following command to create config files and set appropriate values.
```bash
cp .env.sample .env
cp .config.yaml.sample config.yaml
```
- Run following command to download dependencies. (development)
```bash
go mod tidy
```
- Run Server in machine. (development)
```bash
go run cmd/api/main.py
```
- To run docker container. Run the following command. (Production)
```bash
docker-compose up -d
```

## API Documentation

#### Sign Up:
```bash
curl -x POST '{BASEURL}/api/v1/signup' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "01901901901@mailinator.com",
    "password": "321dsa"
}'
```

#### Log In:
```bash
curl -x POST '{BASEURL}/api/v1/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "01901901901@mailinator.com",
    "password": "321dsa"
}'
```

#### Create Order:
```bash
curl -X POST '{BASEURL}/api/v1/orders' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer your-token' \
--data '{
  "recipient_name": "Rashed",
  "recipient_phone": "01310589625",
  "recipient_address": "Banai, Dhaka",
  "store_id": 52560,
  "recipient_city": 7,
  "delivery_type": 2,
  "item_quantity": 10,
  "merchant_order_id": "65892",
  "item_weight": 15,
  "item_type": 2,
  "amount_to_collect": 15000,
  "status": "pending"
}
```

#### List Orders:
```bash
curl '{BASEURL}/api/v1/orders/all' \
--header 'Authorization: Bearer your-token'
```

#### Cancel Order:
```bash
curl -x PUT '{BASEURL}/api/v1/orders/SD310125IJGYAG/cancel'\
--header 'Authorization: Bearer your-token'
```

### Future Scope
- Unit Test
- CI/CD