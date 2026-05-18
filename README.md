# docker-time-analysis

Statistical analysis tool for Docker container startup times using regression modeling and interactive visualizations.

## Features

- CSV file upload and parsing
- Statistical analysis (correlation tables, descriptive statistics)
- Three regression model types: Linear, Exponential, Piecewise
- Statistical significance testing (R², Residual Variation, F-Stat Test, Pearson Coefficient)
- Interactive charts and visualizations with ECharts

## Collected Data

The data provided in [`/data`](/data) was collected inside a virtual machine (VMWare) running Ubuntu Server 24.04 LTS with the following configurations:

- Intel Core i5-10300H 4 Cores
- 8 GB RAM
- 80 GB SSD

System degradation at around 1000 Alpine Linux containers can be explained by resource exhaustion to create and/or sustain Docker containers.

Users are free to collect their own data by modifying/running [`/data/collect_data.sh`](/data/collect_data.sh)

## Screenshots

![overview](/README/overview.gif)
<p align="center">Overview of the application</p>

![chart](/README/chart.png)
<p align="center">Visualization of regression models</p>

![regression models data](/README/regression_models_data.png)
<p align="center">Regression models data</p>

![regression models validation](/README/regression_models_validation.png)
<p align="center">Regression models validation</p>

## Tech Stack

- **Backend**: Go, Gonum (statistical computing)
- **Frontend**: Svelte, TypeScript, ECharts

## Project Structure

```
backend/
├── cmd/              - Entry points
├── internal/
│   ├── parse/        - CSV parsing and data validation
│   ├── compute/      - Statistical analysis and regression models
│   └── server/       - HTTP API server
frontend/
├── src/              - Svelte components and TypeScript logic
├── dist/             - Built output
└── vite.config.ts    - Build configuration
data/                 - Sample CSV data files
```

## Getting Started

### Prerequisites

- Go 1.26.1 or later
- Node.js 18+ with npm

### Run as Docker Compose

```bash
git clone https://github.com/ItakawaM/docker-time-analysis
docker-compose up -d
```

### Run Manually

**Backend:**

```bash
cd backend
go build -o server.exe ./cmd
./server.exe
```

or

```bash
cd backend
make build
./server.exe
```

**Frontend:**

```bash
cd frontend
npm install
npm run dev    # Development server
npm run build  # Production build
```

The server runs on `http://localhost:8080` by default.

## Requirements

- **Go**: 1.26.1+
- **Node.js**: 18+

## License

MIT License — see [`LICENSE`](LICENSE)
