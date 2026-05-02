## 2025-02-28 - [API Gateway Defense-in-Depth]
**Vulnerability:** Missing security headers (`X-Content-Type-Options`, `X-Frame-Options`, `X-XSS-Protection`, `Strict-Transport-Security`) at the API Gateway level, leaving downstream microservices exposed to common web vulnerabilities (MIME sniffing, clickjacking, XSS) if their individual implementations lack these headers.
**Learning:** In a microservices architecture utilizing an API Gateway pattern (like with Consul/Gin in this repo), implementing global security middleware at the gateway layer is an effective defense-in-depth strategy. This ensures a consistent baseline of security across all reverse-proxied routes (`/api/*`), regardless of individual service configurations.
**Prevention:** Always apply common security headers globally at the API Gateway or edge reverse proxy level to provide baseline protection for all downstream services.

## 2025-02-28 - [API Gateway Unauthorized Internal Service Exposure]
**Vulnerability:** The API Gateway used a dynamic wildcard route (`/api/:service/*path`) to proxy requests to any service registered in Consul based entirely on user input. This allowed public attackers to access internal-only microservices that were not intended to be exposed to the internet, potentially leading to unauthorized data access or Server-Side Request Forgery (SSRF) style pivots.
**Learning:** In a dynamic proxy architecture utilizing service discovery mechanisms (like Consul), trusting user-provided routing inputs without validation is dangerous. Any service registered with the backend system becomes publicly accessible by default.
**Prevention:** Always implement an explicit whitelist at the gateway/proxy layer to restrict which downstream services are publicly accessible.

## 2024-04-30 - Path Traversal in API Gateway Reverse Proxy
**Vulnerability:** API Gateway takes an arbitrary URL parameter and directly passes it to backend services without sanitization. This allowed path traversal sequences like `../../../` to be passed directly to backend services via reverse proxy, creating SSRF risks.
**Learning:** Default proxy handlers might forward raw, unescaped, or unnormalized paths. Go's Gin router doesn't automatically normalize `c.Param("path")` values against directory traversal sequences if they are passed dynamically into downstream proxies.
**Prevention:** Always normalize and validate external inputs that manipulate file paths or internal URL routing. Use `path.Clean("/" + targetPath)` for proxy target paths.

## 2026-05-02 - [API Gateway Internal Endpoint Exposure]
**Vulnerability:** The API Gateway permitted external requests to reach paths containing `/internal/` on downstream microservices. Because downstream services rely on `X-Gateway-Service` (added by the gateway) to identify authorized internal traffic, external requests forwarded by the gateway inherently bypassed `RequireInternalService()` middlewares.
**Learning:** In architectures where downstream services trust headers injected by the API gateway to authenticate internal traffic, the gateway must act as a strict barrier. It must actively block external requests to internal namespaces because the act of proxying them grants them internal trust by default.
**Prevention:** Explicitly block access to internal endpoints (e.g., paths starting with `/internal/`) in the API Gateway's reverse proxy logic.
