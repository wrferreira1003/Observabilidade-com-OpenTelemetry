# Projeto: Observabilidade com OpenTelemetry

## Objetivo do Projeto

Este projeto implementa dois serviços, **Serviço A** e **Serviço B**, que se comunicam entre si utilizando tracing distribuído com OpenTelemetry. A aplicação demonstra como rastrear requisições de ponta a ponta entre dois serviços, incluindo a propagação de contexto e a medição de tempo de resposta para operações críticas.

Os serviços estão configurados para exportar os spans gerados para um **OTEL Collector**, que os encaminha para o **Zipkin**. O Zipkin é usado para visualizar os traces.

## Estrutura do Projeto

- **Serviço A**: Recebe requisições com um CEP e encaminha para o Serviço B.
- **Serviço B**: Processa o CEP, busca a localização e a temperatura associada ao CEP.

## Requisitos

- Docker e Docker Compose instalados na máquina.
- Arquivo `cep.http` na raiz do projeto para realizar requisições de teste.

## Como Rodar o Projeto

1. Clone o repositório e entre no diretório do projeto:
   ```bash
   git clone <url-do-repositorio>
   cd Observabilidade-open-telemetry

2. Execute o Docker Compose para iniciar os serviços:
   ```bash
   docker-compose up -d

3. Execute a requisição no arquivo `cep.http`

4. Acesse o Zipkin no navegador: http://localhost:9411/zipkin/