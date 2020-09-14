# Login and Consent app for Hydra tutorial

This app is a porting of https://github.com/ory/hydra-login-consent-node in Go.

## Quickstart

Create a network:

`docker network create hydraguide`

Create database for Hydra:

```sh
docker run \
  -p 5432:5432 \
  --network hydraguide \
  --name ory-hydra-example--postgres \
  -e POSTGRES_USER=hydra \
  -e POSTGRES_PASSWORD=secret \
  -e POSTGRES_DB=hydra \
  -d postgres:9.6
  
export SECRETS_SYSTEM=$(export LC_CTYPE=C; cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)
export DSN=postgres://hydra:secret@ory-hydra-example--postgres:5432/hydra?sslmode=disable
docker pull oryd/hydra:v1.7.4
docker run -it --rm \
  --network hydraguide \
  oryd/hydra:v1.7.4 \
  migrate sql --yes $DSN
```

Run Hydra (configured to route login and consent to our Login and Consent app):

```sh
docker run -d \
  --name ory-hydra-example--hydra \
  --network hydraguide \
  -p 9000:4444 \
  -p 9001:4445 \
  -e SECRETS_SYSTEM=$SECRETS_SYSTEM \
  -e DSN=$DSN \
  -e URLS_SELF_ISSUER=https://localhost:9000/ \
  -e URLS_CONSENT=http://localhost:9020/consent \
  -e URLS_LOGIN=http://localhost:9020/login \
  -e STRATEGIES_ACCESS_TOKEN=jwt \
  -e SECRETS_SYSTEM=prontopro-jwt-secret \
  oryd/hydra:v1.7.4 serve all
  ```
  
 Start Login and Consent app:
 
 ```sh
go run ./...
 ```
 
 Register an example client in Hydra and start it using Hydra demo utilities:
  
 ```sh
docker run --rm -it \
  -e HYDRA_ADMIN_URL=https://ory-hydra-example--hydra:4445 \
  --network hydraguide \
  oryd/hydra:v1.7.4 \
  clients create --skip-tls-verify \
    --id facebook-photo-backup \
    --secret some-secret \
    --grant-types authorization_code,refresh_token,client_credentials,implicit \
    --response-types token,code,id_token \
    --scope openid,offline,photos.read \
    --callbacks http://127.0.0.1:9010/callback
    
docker run --rm -it \
  --network hydraguide \
  -p 9010:9010 \
  oryd/hydra:v1.7.4 \
  token user --skip-tls-verify \
    --port 9010 \
    --auth-url https://localhost:9000/oauth2/auth \
    --token-url https://ory-hydra-example--hydra:4444/oauth2/token \
    --client-id facebook-photo-backup \
    --client-secret some-secret \
    --scope openid,offline,photos.read
```
  
If you want to play with https://github.com/lorenzoranucci/hydra-client-go client:
   
  ```sh
docker run --rm -it \
   -e HYDRA_ADMIN_URL=https://ory-hydra-example--hydra:4445 \
   --network hydraguide \
   oryd/hydra:v1.7.4 \
   clients create --skip-tls-verify \
     --id client-frontend.localhost \
     --secret some-secret \
     --grant-types authorization_code,refresh_token,client_credentials,implicit \
     --response-types token,code,id_token \
     --scope openid,offline \
     --callbacks http://client-frontend.localhost:9011/login-callback
     
cd ../hydra-client-go
go run ./...
 ```
