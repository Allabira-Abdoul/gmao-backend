## 2024-06-12 - Inter-Service Communication should be abstracted behind Secondary Ports
**Learning:** Hardcoding `http.Client` and service registry lookups directly into application service methods (like `AuthService`) tightly couples the business logic to specific networking protocols and discovery tools.
**Action:** Always create a secondary port (e.g., `UserClient` interface) in the domain layer and implement an HTTP adapter for it in `adapters/secondary/http/`. This enforces DIP, making the application logic cleaner and highly testable.
