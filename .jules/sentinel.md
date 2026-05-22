## 2024-05-22 - Hardcoded JWT Secrets

**Vulnerability:** The codebase had fallback hardcoded JWT secrets (`"gmao-dev-secret-change-in-production"`) in multiple microservices for `JWT_SECRET` when not provided in the environment.

**Learning:** Fallback secrets mean that if production environment configuration is forgotten, the service starts with a publicly known secret, enabling attackers to forge valid JWTs and gain full unauthorized access.

**Prevention:** Always enforce strict configuration requirements. If a critical security variable like `JWT_SECRET` is missing, the service must explicitly fail to start (`log.Fatal`) rather than silently using an insecure fallback.
