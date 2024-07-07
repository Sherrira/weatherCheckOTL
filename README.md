# weatherCheckOTL
Desafio Observabilidade e OTL - Go Express

## Índice
1. [Objetivo](#objetivo)
2. [Descrição](#descrição)
3. [Configurações de ambiente](#configurações-de-ambiente)
4. [Pré-requisitos](#pré-requisitos)
5. [Executando o projeto](#executando-o-projeto)


## Objetivo
O objetivo deste projeto é avaliar através do Zipkin o tracing dos spans de execução do projeto usando o padrão de Open Telemetry.

## Descrição
Este projeto consiste na execução de dois serviços em paralelo (ServicoA e ServicoB), e executa também outros dois servidores para telemetria (Zipkin e OTEL collector).

### ServicoA
O ServicoA é responsável por:
- Validar o CEP 
- Fazer uma chamada para o ServicoB
- Retornar a resposta obtida pelo ServicoB
- A aplicação possui um endpoint HTTP POST que responde na porta 8080. 

### ServicoB
O ServicoB é responsável por:
- Receber o CEP do ServicoA e utilizar o mesmo para:
    - Consultar a localidade do CEP
    - Consultar a temperatura na localidade do CEP
- Retornar ao ServicoA os dados de localidade e temperatura
- A aplicação possui um endpoint HTTP GET que responde na porta 8081.

### Zipkin
O Zipkin é um sistema de tracing distribuído que permite rastrear os spans de execução do projeto.

### OTEL Collector
O OTEL collector é um container que executa o coletor de dados no padrão Open Telemetry.
Neste projeto foi configurado um pipeline que recebe dados via gRPC na porta 4317 e envia tais dados ao Zipkin em sua porta padrão (9411).
As configurações de pipeline são feitas no arquivo `.docker/otel-collector-config.yaml`.

## Pré-requisitos
Assegure-se de ter as seguintes ferramentas instaladas:
- [Golang](https://go.dev/doc/install)
- [Docker](https://docs.docker.com/compose/install/)

## Configurações de ambiente
Para executar as consultas de temperatura é necessária a obtenção de uma chave de acesso à API WeatherAPI. Para isso, siga os passos abaixo:    
1. Acesse o site [WeatherAPI](https://www.weatherapi.com/) e crie uma conta.
2. Após a criação da conta, acesse o painel de controle e copie a chave de acesso.
3. No arquivo ".env" na pasta cmd do servicoA, adicione a chave de acesso obtida no passo anterior.

Nota: o projeto está entregando um arquivo docker-compose.yaml para a execução do projeto com containers.

## Executando o projeto

O projeto possui um Makefile com comandos utilitários para execução do projeto listados abaixo:

- Constrói as imagens e sobe os containers do projeto:
```
$ make build-run
```

- Sobe os containers do projeto sem build:
```
$ make run
```

- Encerra a execução dos containers do projeto:
```
$ make stop
```

Após a criação das imagens e subida dos containers, é possível, realizar uma chamada através do arquivo  `ServicoA/api/validade_cep.http` que inicialmente fará a validação do número do CEP passado no body da request.
```
{
    "cep": "08223110"
}
```
Uma vez que o CEP informado seja válido, uma chamada será feita automaticamente ao ServicoB que identificará a localização do CEP indicado, e fará a consulta e cálculo da temperatura deste local em ºC, ºF e ºK, retornando o resultado no body da response.

Para a avaliação do tracing entre os serviços, acesse o console do Zipkin através da url `localhost:9411`.