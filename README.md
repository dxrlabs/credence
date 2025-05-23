# ![Go Test](https://github.com/dxrlabs/credence/actions/workflows/go.yml/badge.svg)
![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)
![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)

# Credence

**Credence** is a Go library that manages authentication for external APIs, handling token acquisition, refresh, and multi-credential orchestration. It wraps the standard `golang.org/x/oauth2` library with a thread-safe, extensible interface.

---

## 🚀 Features

- 🔐 OAuth2 Client Credentials Grant support
- 🔁 Automatic token caching and refreshing
- 🧵 Safe for concurrent use
- 🔑 Supports multiple credentials with isolated token lifecycles
- 📦 Clean, minimal API

---

## 📦 Install

```bash
go get github.com/dxrlabs/credence
```

---

## 🛠️ Usage

### 1. Register Credentials
```go
credence.Register("stripe", credence.Config{
	ClientID:     "your-client-id",
	ClientSecret: "your-client-secret",
	TokenURL:     "https://api.stripe.com/oauth/token",
	Scopes:       []string{"read", "write"},
})
```

### 2. Get Token
```go
token, err := credence.Token(context.Background(), "stripe")
if err != nil {
	log.Fatal(err)
}

req.Header.Set("Authorization", "Bearer "+token)
```

---

## 🔧 Roadmap
- TokenStore interface (in-memory, Redis, BoltDB, etc.)
- Retry-aware HTTP client middleware
- Jitter/backoff for refresh cycles
- OpenTelemetry / Prometheus integration

---

## 🤝 Contributing
PRs welcome! Open an issue if you have a feature request or bug.

---

## 📄 License
MIT
