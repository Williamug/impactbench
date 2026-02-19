# ImpactBench 🚀

[![Test and Build](https://github.com/williamug/impactbench/actions/workflows/test.yml/badge.svg)](https://github.com/williamug/impactbench/actions/workflows/test.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/williamug/impactbench)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**ImpactBench** is a global CLI tool designed for engineering teams to quantify the impact of their changes. It combines benchmarking, load testing, and regression-aware performance analysis to ensure that every refactor or optimization is backed by measurable data.

---

## 🌟 Key Features

- **Global CLI**: Single binary implementation for seamless performance auditing.
- **Multi-Level Config**: Intelligent merging of global, project-specific, and default configurations.
- **Regression Detection**: Automated threshold evaluation to block performance regressions in CI/CD.
- **Load Testing**: High-concurrency engine using Goroutines and worker pools with percentile analytics (P95, P99).
- **Premium Visualization**: Professional, high-fidelity terminal tables using `go-pretty`.
- **Flexible Storage**: Dual-backend support for SQLite and JSON.

---

## 🚀 Installation

### From Source
```bash
git clone https://github.com/williamug/impactbench.git
cd impactbench
go build -o impactbench ./cmd/impactbench
sudo mv impactbench /usr/local/bin/
```

---

## 🛠 Usage

### 1. Simple Benchmark
Capture performance metrics for a single URL.
```bash
impactbench run --url=http://localhost:8080 --label=baseline
```

### 2. Load Testing
Stress test your system with concurrent users.
```bash
impactbench loadtest --url=http://localhost:8080 --users=50 --duration=30
```

### 3. Comparison
Compare two snapshots to see the delta of your improvements.
```bash
impactbench compare --baseline=run_1 --current=run_2
```

### 4. Automated Review
Evaluate the latest benchmark against previous results using configurable thresholds.
```bash
impactbench review --fail-on-regression
```

---

## ⚙️ Configuration

ImpactBench uses a hierarchical configuration system (Viper):
1. **Default**: Hardcoded fallbacks.
2. **Global**: `~/.impactbench/config.yaml`
3. **Project**: `./.impactbench/config.yaml`

### Example Config
```yaml
base_url: http://localhost:8080
framework: laravel

thresholds:
  response_time: 10.0 # % change allowed
  error_rate: 2.0
  memory: 15.0

load_test:
  default_users: 100
  default_duration: 60
```

---

## 🏛 Architecture

Built with a modular, adapter-based architecture:
- **CLI Layer**: Cobra-based command handling.
- **Runner Engine**: Business logic for benchmark orchestration.
- **Load Engine**: Concurrency pool for high-traffic simulation.
- **Adapter Layer**: Framework-agnostic interfaces (HTTP, Laravel, Django).
- **Storage Layer**: Immutable snapshot management.

---

## 🧪 Testing

Run the comprehensive test suite to verify internal logic:
```bash
go test -v ./test/...
```

The CI/CD pipeline (GitHub Actions) automatically verifies tests and builds on every push to `main`.

---

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

---

*ImpactBench is an internal tool designed by Nugsoft Engineering.*
