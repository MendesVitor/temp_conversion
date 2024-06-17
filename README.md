### README

Este projeto consiste em um serviço web desenvolvido em Go que recebe um CEP como parâmetro, consulta a API do ViaCEP para obter a localização correspondente, e então utiliza a API do WeatherAPI para retornar informações climáticas dessa localização, incluindo a temperatura atual em Celsius, Fahrenheit e Kelvin. O serviço está hospedado no Google Cloud Run.

#### Funcionamento

1. **Requisitos:**

    - O sistema espera um CEP válido de 8 dígitos como parâmetro na URL.
    - O serviço consulta o ViaCEP para obter a cidade correspondente ao CEP.
    - Em seguida, consulta a WeatherAPI para obter as condições climáticas atuais da cidade.

2. **Endpoints:**

    - **GET /clima?cep={cep}**
        - Retorna um JSON com as temperaturas atuais em Celsius, Fahrenheit e Kelvin para o CEP especificado.
        - Exemplo de URL: [https://temp-conversion-bvmcjo6y2a-uc.a.run.app/clima?cep=18017005](https://temp-conversion-bvmcjo6y2a-uc.a.run.app/clima?cep=18017005)

3. **Respostas:**

    - **Sucesso (200 OK):**
        ```json
        {
            "temp_C": "28.5",
            "temp_F": "83.3",
            "temp_K": "301.6"
        }
        ```
    - **Falha (422 Unprocessable Entity):**

        - Quando o CEP não tem o formato correto.
        - Exemplo:
            ```json
            "invalid zipcode"
            ```

    - **Falha (404 Not Found):**
        - Quando o CEP não é encontrado.
        - Exemplo:
            ```json
            "can not find zipcode"
            ```

#### Deploy no Google Cloud Run

Este projeto está hospedado no Google Cloud Run e pode ser acessado através da URL [https://temp-conversion-bvmcjo6y2a-uc.a.run.app/clima?cep=18017005](https://temp-conversion-bvmcjo6y2a-uc.a.run.app/clima?cep=18017005).

#### Como Executar Localmente

Para executar este projeto localmente, você pode escolher entre duas abordagens: diretamente com o Go ou utilizando Docker.

##### Executar com Go

1. **Clone o repositório:**

    ```
    git clone <URL do Repositório>
    ```

2. **Instale as dependências:**

    ```
    go mod tidy
    ```

3. **Execute o servidor local:**

    ```
    go run main.go
    ```

4. **Acesse o serviço localmente através de:**

    ```
    http://localhost:8080/clima?cep=18017005
    ```

##### Executar com Docker

1. **Clone o repositório:**

    ```
    git clone <URL do Repositório>
    ```

2. **Execute com Docker:**

    ```
    docker compose up
    ```

3. **Acesse o serviço localmente através de:**

    ```
    http://localhost:8080/clima?cep=18017005
    ```
