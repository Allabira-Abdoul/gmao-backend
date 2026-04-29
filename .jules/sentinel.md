## 2025-02-28 - [API Gateway Defense-in-Depth]
**Vulnerability:** Missing security headers (`X-Content-Type-Options`, `X-Frame-Options`, `X-XSS-Protection`, `Strict-Transport-Security`) at the API Gateway level, leaving downstream microservices exposed to common web vulnerabilities (MIME sniffing, clickjacking, XSS) if their individual implementations lack these headers.
**Learning:** In a microservices architecture utilizing an API Gateway pattern (like with Consul/Gin in this repo), implementing global security middleware at the gateway layer is an effective defense-in-depth strategy. This ensures a consistent baseline of security across all reverse-proxied routes (`/api/*`), regardless of individual service configurations.
**Prevention:** Always apply common security headers globally at the API Gateway or edge reverse proxy level to provide baseline protection for all downstream services.

## 2025-02-28 - [API Gateway Unauthorized Internal Service Exposure]
**Vulnerability:** The API Gateway used a dynamic wildcard route (`/api/:service/*path`) to proxy requests to any service registered in Consul based entirely on user input. This allowed public attackers to access internal-only microservices that were not intended to be exposed to the internet, potentially leading to unauthorized data access or Server-Side Request Forgery (SSRF) style pivots.
**Learning:** In a dynamic proxy architecture utilizing service discovery mechanisms (like Consul), trusting user-provided routing inputs without validation is dangerous. Any service registered with the backend system becomes publicly accessible by default.
**Prevention:** Always implement an explicit whitelist at the gateway/proxy layer to restrict which downstream services are publicly accessible.
