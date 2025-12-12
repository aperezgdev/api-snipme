# SnipMe API


[![codecov](https://codecov.io/github/aperezgdev/api-snipme/graph/badge.svg?token=U1TL84B47F)](https://codecov.io/github/aperezgdev/api-snipme)


API para acortamiento de enlaces con autenticación OAuth, analítica avanzada, geolocalización y monitoreo. Desarrollada en Go, lista para producción y despliegue con Docker.

## Características

- **Acortamiento de enlaces** con gestión de usuarios.
- **Autenticación OAuth** (Google, GitHub).
- **Analítica de enlaces**: vistas totales, únicas, por país, etc.
- **Geolocalización** de visitas (MaxMind GeoLite2).
- **Persistencia** en PostgreSQL y Redis.
- **Monitoreo y métricas** con Prometheus, Loki y Grafana.
- **Contenedores Docker** listos para desarrollo y producción.
- **Cobertura de tests** y utilidades para pruebas automáticas.

## Estructura del Proyecto

```
.
├── db/                # Esquemas SQL, seeds y modelos generados
├── src/               # Código fuente principal
│   ├── cmd/           # Entrypoint y bootstrap
│   └── internal/      # Contextos de dominio, infraestructura y aplicación
├── test/              # Pruebas de integración y helpers
├── Dockerfile         # Build multi-stage para Go y runtime seguro
├── compose.yaml       # Orquestación de servicios (API, DB, Redis, Prometheus, Loki, Grafana)
├── openapi.yaml       # Especificación OpenAPI de la API
├── prometheus.yml     # Configuración de Prometheus
├── loki.yml           # Configuración de Loki
└── README.md
```

## Requisitos

- Go 1.24+
- Docker y Docker Compose
- Make (opcional para scripts)
- Acceso a MaxMind GeoLite2 (archivo incluido en `db/geo/`)

## Instalación y Ejecución

### 1. Variables de entorno

Copia y ajusta tu archivo `.env`:

```bash
cp .env.example .env
```

Configura los valores para DB, Redis, OAuth, JWT, etc.

### 2. Levantar el entorno con Docker Compose

```bash
docker compose up --build
```

Esto levantará:
- API (Go)
- PostgreSQL
- Redis
- Prometheus
- Loki
- Grafana

La API estará disponible en el puerto configurado (`8081` por defecto).

### 3. Acceso a servicios

- **API**: http://localhost:8081
- **Grafana**: http://localhost:3001 (user/pass: admin/admin)
- **Prometheus**: http://localhost:9090
- **Loki**: http://localhost:3100

## Uso

Consulta la documentación OpenAPI en `openapi.yaml` o accede a `/swagger` si está habilitado.

### Endpoints principales

- `POST /public-short-links` - Crear enlace corto
- `GET /{code}` - Redireccionar enlace corto
- `GET /link-analytics/{code}` - Analítica de un enlace
- `POST /auth/google` - Login con Google
- `POST /auth/github` - Login con GitHub

## Pruebas

Ejecuta los tests con:

```bash
go test ./...
```

O usando los scripts de test en el directorio `test/`.

## Monitoreo y métricas

- **Prometheus** recolecta métricas de la API.
- **Loki** almacena logs estructurados.
- **Grafana** permite visualizar dashboards de métricas y logs.

## Licencia

Este proyecto está licenciado bajo la licencia GNU GPL v3. Consulta el archivo `LICENSE` para más detalles.
