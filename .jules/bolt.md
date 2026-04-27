## 2024-05-23 - API Gateway Consul Discovery Bottleneck
**Learning:** The API Gateway queries Consul (`client.Health().Service()`) synchronously on *every single incoming request* to route traffic. In a microservices architecture, this adds at least 1-10ms of network latency per request and puts heavy load on the Consul server, severely limiting gateway throughput.
**Action:** When implementing service discovery routing in gateways, always cache the discovery results (even with a short TTL like 5-15 seconds) to avoid blocking the main request path on network calls to the service registry.
