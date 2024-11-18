# Account data storing immudb vault demo
This is my demo application for utilizing immudb Vault API for storing banking account like data. <br/>
It consists of two parts backend created in golang and frontend created in React with typescript.

Application is storing accounting information within immudb Vault
with the following structure:
account number (unique), account name, iban, address,
amount, type (sending, receiving). <br/>
It has an backend API to add and retrieve accounting information. <br/>
It has a frontend that displays accounting information and allows to create new records. <br/>
Open `./screen.png` for frontend screnshot.

## Running
Get API token for your `https://vault.immudb.io/docs/api/v1` <br/>
Create `.env` file or replace `environment.API_PRIVATE_KEY` in `docker-compose.yaml` <br/>
Sample `.env` file: <br/>
```sh
export API_PRIVATE_KEY="default.your.private-apiToken"
```
Install docker and docker-compose and run
```sh
docker-compose up
```
You can now access frontend page at http://localhost:3000  <br/>
You can play with backend via http://localhost:1323.

## API 
Backend server exposes following two endpoints: <br/>
`GET /api/v1/account` - lists all created accounts.  <br/>
`POST /api/v1/account` - create new account based on JSON data input body. <br/>
Sample JSON data:
```json
{
  "number": 6,
  "name": "John Doe",
  "iban": "US10105097603123",
  "address": "123 Sesame Street",
  "amount": 1234,
  "type": "sending"
}
```

## immudb vault usage
Backend will store account data under `default` ledger and `default` collection. <br/>
If `default` collection does not exist it will be created by backend on startup.

## Code structure
- ./front - react frontend app
- ./cmd/server.go - backend echo server for API communication 
- ./pkg/logger - global zap logger
- ./pkg/client - custom immudb client for HTTP requests and some common operations wrapper (api.go) like listing collections in ledger.
- ./pkg/account - account model and account manager responsible for all database like operations

## Developing
### Backend:
Import key and run standalone backend with test data creation:
Currently server runs at 1323 port by default. 
```sh
. .env
ADD_TEST_DATA=true go run cmd/server.go 
```
### Frontend
See `front/README.md`
Run dev front on port 3000:
```sh
npm start dev
```

