# Security and Credential Protection Rule

- Do NOT include remote passwords, secrets, or production connection strings in source code (e.g., Go files, fallback arguments) or `compose.yaml`.
- All remote credentials must strictly be loaded from environment variables (e.g., Lambda environment variables injected by IaC/CI, local `.env`) and never committed or written into source files.
- `compose.yaml` is strictly for local container development (`ocg_db_local`) and must always use local dummy credentials (`postgres`/`postgres`).
