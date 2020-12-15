# XPTO Support Application
Essa aplicação tem como objetivo integrar as informações e escoar o fluxo de requisições da aplicação legada XPTO.

# Inicialização da aplicação

## Dependências
Para compilar o código da aplicação é necessário o go na versão > 1.13

Para gerar imagem docker e rodar a stack docker-compose é necessário o docker.

## Build
* Para compilar o código da aplicação utilize o comando: `make build`

## Docker
* Para gerar imagem docker da aplicação utilize o comando: `make docker-build`
* Para inicializar a stack completa (incluindo a imagem do xpto-support), rode o comando `make init-complete-stack`

## Stack de Serviços de Suporte
Toda stack de serviços para suporte da aplicação está configurada para rodar no docker, a partir do arquivo
docker-compose.

* Para inicializar a stack de serviços utilize o comando: `make init-stack`

# Testes
* Para rodar os testes da aplicação, rode o comando: `make run-test`
* Para rodar os testes de aceitação, rode o comando: `make run-integration-test` (é necessário estar rodando a stack
completa de serviços)

# Instrumentação
A aplicação é monitorada através do APM Server do Elasticsearch, a stack completa dispôe de uma UI com Kibana para 
monitrar os processos internos da aplicação. Para acessar a UI utilize o endereço http://localhost:5601/app/apm#/services/xpto-support

# Documentação da API:

## Sistema legado
Gestão da integração com sistema legado

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

## Negativados
Controle de débito de negativados

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
