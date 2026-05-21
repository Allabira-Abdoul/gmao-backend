# Audit Service API Specification

The **Audit Service** provides a secure, tamper-proof, write-once read-many (WORM) audit logger that records all user and automated service-to-service operations across the entire GMAO microservices ecosystem. It registers with Consul under the name `audit-service` and operates on port `8087`.

All external client requests are routed through the API Gateway using the prefix `/api/audit/*`. Downstream routes beginning with `/internal` are strictly isolated from direct external access.

---

## 🔒 Security & Privileges

The service enforces strict boundary controls depending on the endpoint type:

1. **Internal Inter-Service Logger (`POST /internal/*`)**: Requires a valid `X-Internal-Service` header signed by a trusted service. Any direct gateway request to this path is blocked by gateway policy.
2. **Auditor Inspector (`GET /audit-logs`)**: Requires user authentication (Bearer JWT) containing the specific `AUDITOR` privilege or the `SYSTEM_ADMIN` privilege.

---

## 📡 Endpoints Specification

### 🟢 List All Audit Logs (Privilege: `AUDITOR` or `SYSTEM_ADMIN`)

Retrieve a chronological log of all recorded actions in the system.

- **HTTP Method**: `GET`
- **Path**: `/api/audit/audit-logs` (routed to `GET /audit-logs` downstream)
- **Authentication**: Bearer JWT

#### Success Response
- **Status Code**: `200 OK`
- **Body**:
```json
{
  "status": "success",
  "data": [
    {
      "id": "e44d5c41-862d-45df-bb78-ecb18360d8ef",
      "service_name": "user-service",
      "action": "USER_CREATE",
      "details": "Registered new technician abdou@gmao.local with role TECHNICIAN",
      "user_id": "8fa1ad78-831e-4cb8-8c10-9bd74130be52",
      "performed_at": "2026-05-19T11:45:00Z"
    },
    {
      "id": "7ca9da22-e421-4fba-bb89-11c9d9d3000b",
      "service_name": "auth-service",
      "action": "USER_LOGIN",
      "details": "User session created via web interface",
      "user_id": "8fa1ad78-831e-4cb8-8c10-9bd74130be52",
      "performed_at": "2026-05-19T11:40:00Z"
    }
  ]
}
```

#### Error Responses
- **Status Code**: `401 Unauthorized` (Token missing or malformed)
```json
{
  "status": "error",
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Authorization header is required and must begin with Bearer"
  }
}
```
- **Status Code**: `403 Forbidden` (User does not possess dynamic `AUDITOR` privilege)
```json
{
  "status": "error",
  "error": {
    "code": "FORBIDDEN",
    "message": "Dynamic privilege 'AUDITOR' is required to perform this action"
  }
}
```

---

### 🟢 Write Audit Log (Internal Only)

Log a secure action token. This endpoint is dedicated to system microservices and cannot be reached externally.

- **HTTP Method**: `POST`
- **Path**: `/internal/audit-logs`
- **Authentication**: `X-Internal-Service` network token check

#### Request Headers
- `X-Internal-Service`: `<secret-service-token>`

#### Request Body
| Field | Type | Required | Validation | Description |
| :--- | :--- | :--- | :--- | :--- |
| `service_name` | `string` | Yes | Non-empty | Microservice registering the event (e.g. `auth-service`, `asset-service`). |
| `action` | `string` | Yes | Non-empty | Domain action (e.g. `USER_LOGIN`, `ASSET_UPDATE`). |
| `details` | `string` | No | Text description | Descriptive context of the action. |
| `user_id` | `string` | No | Valid UUID | Optional ID of the user triggering the action. |

##### Example
```json
{
  "service_name": "asset-service",
  "action": "ASSET_CREATE",
  "details": "Created asset CNC Milling Machine (Code: CNC-04)",
  "user_id": "8fa1ad78-831e-4cb8-8c10-9bd74130be52"
}
```

#### Success Response
- **Status Code**: `201 Created`
- **Body**:
```json
{
  "status": "success",
  "data": {
    "id": "cf63db24-f719-482a-a9bd-8311d9d93ee4",
    "service_name": "asset-service",
    "action": "ASSET_CREATE",
    "details": "Created asset CNC Milling Machine (Code: CNC-04)",
    "user_id": "8fa1ad78-831e-4cb8-8c10-9bd74130be52",
    "performed_at": "2026-05-19T12:05:00Z"
  }
}
```

#### Error Responses
- **Status Code**: `403 Forbidden` (Direct access via gateway or missing internal secret header)
```json
{
  "status": "error",
  "error": {
    "code": "FORBIDDEN",
    "message": "Only internal network calls are allowed to perform this operation"
  }
}
```
