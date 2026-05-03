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

## 2024-05-24 - API Gateway Bypass to Internal Endpoints
**Vulnerability:** The API Gateway allowed external clients to access `/internal/` service-to-service endpoints by directly specifying the path in the proxied URL. Furthermore, it automatically injected the `X-Gateway-Service` header on all proxied requests and allowed external clients to spoof the `X-Internal-Service` header, fully satisfying downstream `RequireInternalService` middleware checks.
**Learning:** Downstream microservices trust headers injected by the API Gateway to authorize internal traffic. The gateway must act as a strict firewall, blocking access to `/internal/` endpoints and scrubbing internal headers from incoming external requests.
**Prevention:** Always ensure the API Gateway properly sanitizes URL paths (e.g. `path.Clean()`) and blocks external access to internal routes before reverse proxying. Always delete internal metadata headers (like `X-Internal-Service`) from incoming requests before forwarding them to downstream services.
