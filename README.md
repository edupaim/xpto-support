# XPTO Support Application

Essa aplicação tem como objetivo integrar as informações e escoar o fluxo de requisições da aplicação legada XPTO.

## Documentação da API:

## Sistema legado

### Integrar dados

Gera a integração dos dados do sistema legado com a aplicação

### Integrar (integrate) [POST]

+ Request (application/json)

    + Headers

            x-auth-key: [access_token]

+ Response 200 (application/json)

          {
              "status": "success",
              "data": {
                  "total": 10
              }
          }

## Sistema legado

### Integrar dados

Gera a integração dos dados do sistema legado com a aplicação

### Buscar Negativados (Negatives Query) [GET]

+ Request (application/json)

    + Headers

            x-auth-key: [access_token]

+ Response 200 (application/json)

          {
              "status": "success",
              "data": [{
                  "companydocument": "04843574000182",
                  "companyname": "dbz s.a.",
                  "customerdocument": "26658236674",
                  "value": 59.99,
                  "contract": "3132f136-3889-4efb-bf92-e1efbb3fe15e",
                  "debtdate": "2015-09-11t23:32:51z",
                  "inclusiondate": "2020-09-11t23:32:51z"
              }]
          }
