# weatherCheckOTL
Desafio Observabilidade e OTL - Go Express

## Índice
1. [Pré-requisitos](#pré-requisitos)
2. [Configurações de ambiente](#configurações-de-ambiente)
3. [Executando o projeto](#executando-o-projeto)
4. [Descrição](#descrição)

## Pré-requisitos
Assegure-se de ter as seguintes ferramentas instaladas:
- [Golang](https://go.dev/doc/install)
- [Docker](https://docs.docker.com/compose/install/)

## Configurações de ambiente
Para executar as consultas de temperatura é necessária a obtençao de uma chave de acesso à API WeatherAPI. Para isso, siga os passos abaixo:    
1. Acesse o site [WeatherAPI](https://www.weatherapi.com/) e crie uma conta.
2. Após a criação da conta, acesse o painel de controle e copie a chave de acesso.
3. No arquivo ".env" na pasta cmd, adicione a chave de acesso obtida no passo anterior.

Nota: o projeto está entregando um arquivo Dockerfile que permite o deploy na Google Cloud Platform no serviço Google Cloud Run.

## Executando o projeto
O projeto possui um Makefile com comandos utilitários para execução do projeto listados abaixo:

- Constroi a imagem e sobe o container do projeto build:
```
$ make build-run
```

- Sobe o container do projeto sem build:
```
$ make run
```

## Descrição
Este projeto consite na execução de dois serviços em paralelo:

# ServicoA
O ServicoA é responsável por:
- Validar o CEP 
- Fazer uma chamada para o ServicoB
- Retornar a respota obtida pelo ServicoB

# ServicoB
O ServicoB é responsável por:
- Receber o CEP do ServicoA e utilizar o mesmo para:
    - Consultar a localidade do CEP
    - Consultar a temperatura na localidade do CEP
- Retornar ao ServicoA os dados de localidade e temperatura

