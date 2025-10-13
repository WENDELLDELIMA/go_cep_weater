# CEP → Clima (Go + Cloud Run)

Serviço em Go que recebe um **CEP (8 dígitos)**, resolve a **cidade/UF** via **ViaCEP** e retorna as temperaturas **Celsius, Fahrenheit e Kelvin** via **WeatherAPI**.

## Endpoint
`GET /weather?cep=<8_digitos>`

### Sucesso
- **HTTP 200**
```json
{ "temp_C": 28.5, "temp_F": 83.3, "temp_K": 301.5 }
```

### Erros
- **HTTP 422** — `invalid zipcode` (formato inválido, precisa ter 8 dígitos)
- **HTTP 404** — `can not find zipcode` (CEP não encontrado na ViaCEP)

## Execução local (Docker)
```bash
export WEATHERAPI_KEY=SEU_TOKEN_WEATHERAPI
docker compose up --build
# Teste
curl "http://localhost:8080/weather?cep=01001000"
```

## Execução local (sem Docker)
```bash
export WEATHERAPI_KEY=SEU_TOKEN_WEATHERAPI
go run ./cmd/server
curl "http://localhost:8080/weather?cep=01001000"
```

## Testes automatizados
Os testes usam servidores HTTP falsos (httptest) e variáveis `VIA_CEP_BASE` e `WEATHERAPI_BASE` para mock.
```bash
go test ./...
```

## Deploy no Google Cloud Run (Buildpacks)
Pré-requisitos: `gcloud` autenticado e projeto selecionado.

```bash
# 1) Subir variável de ambiente com a chave da WeatherAPI
gcloud run deploy cep-weather \
  --source . \
  --region southamerica-east1 \
  --allow-unauthenticated \
  --set-env-vars WEATHERAPI_KEY=SEU_TOKEN_WEATHERAPI

# Após o deploy, o comando imprime a URL pública (Cloud Run URL)
# Teste:
curl "https://cep-weather-619173290419.southamerica-east1.run.app/weather?cep=08583450"
```

> Também é possível construir a imagem Docker e publicar no Artifact Registry, depois `gcloud run deploy --image ...`.

## Observações
- `F = C * 1,8 + 32`
- `K = C + 273` 
