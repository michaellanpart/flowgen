# Makefile for FlowGen
# Usage: make <target>

SHELL := /bin/bash

# Tools (override if needed)
GO ?= go
NPM ?= npm
PORT ?= 3001

# Binary output
BACKEND_BIN := backend/flowgen-backend

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help
	@echo "Available targets:" && \
	grep -E '^[a-zA-Z0-9_-]+:.*?## ' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS=":.*?## "}; {printf "\033[36m%-22s\033[0m %s\n", $$1, $$2}'

# -----------------------------
# Backend (Go)
# -----------------------------
.PHONY: backend-build
backend-build: ## Build backend binary
	@echo "[backend] building..."
	@cd backend && $(GO) build -o ../$(BACKEND_BIN) ./cmd

.PHONY: backend-run
backend-run: ## Run backend server locally (PORT=$(PORT))
	@echo "[backend] starting on :$(PORT) ..."
	@cd backend && PORT=$(PORT) $(GO) run ./cmd

.PHONY: backend-test
backend-test: ## Run backend tests
	@echo "[backend] running tests..."
	@cd backend && $(GO) test ./... -v

.PHONY: backend-lint
backend-lint: ## Lint backend (go fmt + go vet)
	@echo "[backend] formatting..."
	@cd backend && $(GO) fmt ./...
	@echo "[backend] vetting..."
	@cd backend && $(GO) vet ./...

# -----------------------------
# Frontend (Node)
# -----------------------------
# Frontend steps are optional; they gracefully skip if npm is unavailable

.PHONY: frontend-install
frontend-install: ## Install frontend dependencies (optional)
	@if [ -f frontend/package.json ]; then \
	  if command -v $(NPM) >/dev/null 2>&1; then \
	    echo "[frontend] installing deps..."; \
	    cd frontend && ($(NPM) ci || $(NPM) install); \
	  else \
	    echo "[frontend] npm not found; skipping install"; \
	  fi; \
	else \
	  echo "[frontend] no package.json; skipping"; \
	fi

.PHONY: frontend-build
frontend-build: ## Build frontend (optional, if scripts exist)
	@if [ -f frontend/package.json ]; then \
	  if command -v $(NPM) >/dev/null 2>&1; then \
	    echo "[frontend] building..."; \
	    cd frontend && $(NPM) run build; \
	  else \
	    echo "[frontend] npm not found; skipping build"; \
	  fi; \
	else \
	  echo "[frontend] no package.json; skipping"; \
	fi

.PHONY: frontend-lint
frontend-lint: ## Lint frontend (optional, if scripts exist)
	@if [ -f frontend/package.json ]; then \
	  if command -v $(NPM) >/dev/null 2>&1; then \
	    echo "[frontend] linting..."; \
	    cd frontend && $(NPM) run lint; \
	  else \
	    echo "[frontend] npm not found; skipping lint"; \
	  fi; \
	else \
	  echo "[frontend] no package.json; skipping"; \
	fi

# -----------------------------
# Composite
# -----------------------------
.PHONY: build
build: backend-build frontend-build ## Build backend and frontend

.PHONY: test
test: backend-test ## Run all tests (currently backend only)

.PHONY: lint
lint: backend-lint frontend-lint ## Lint backend and frontend

.PHONY: run
run: backend-run ## Run backend (serves frontend static files)

.PHONY: clean
clean: ## Remove build artifacts
	@echo "[clean] removing artifacts..."
	@rm -rf $(BACKEND_BIN) frontend/dist frontend/build
