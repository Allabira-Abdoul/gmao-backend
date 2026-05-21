# Authentication Service API Specification

The **Authentication Service** is responsible for managing session lifecycles. It registers with Consul under the name `authentication-service` and operates on port `8081`. 

All external requests are routed via the API Gateway using the prefix `/api/authentication/*`.

---

## 📡 Endpoints Specification

### 🟢 Validate Session (Public)

#### Endpoint Description
Verifies if a session token is active, valid, and not expired. Used by other microservices or the gateway to perform session checks.

- **HTTP Method**: `POST`
- **Path**: `/api/authentication/sessions/validate` (routed to `POST /sessions/validate` downstream)
- **Authentication**: None (Public / Internal)

#### Request Query Parameters
| Parameter | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `token` | `string` | Yes | The active session token to validate. |

#### Request Example
`POST http://localhost:8080/api/authentication/sessions/validate?token=v4.local.eyJleHAiOiIyMDI2...`

#### Success Response
- **Status Code**: `200 OK`
- **Headers**: `Content-Type: application/json`
- **Body**:
```json
{
  "id": "e9c15ad2-f673-455b-86d7-4632b50fe732",
  "user_id": "c3b99db1-d419-48e0-bb15-081079d38bb1",
  "token": "v4.local.eyJleHAiOiIyMDI2...",
  "expired_at": "2026-05-20T11:34:05Z",
  "created_at": "2026-05-19T11:34:05Z"
}
```

#### Failure Response (Token Missing)
- **Status Code**: `400 Bad Request`
- **Body**:
```json
{
  "error": "Query parameter token is required"
}
```

#### Failure Response (Token Invalid or Expired)
- **Status Code**: `401 Unauthorized`
- **Body**:
```json
{
  "error": "session expired or invalid"
}
```

---

### 🟢 Create Session (Authenticated)

#### Endpoint Description
Registers a new session in the system for a user.

- **HTTP Method**: `POST`
- **Path**: `/api/authentication/sessions` (routed to `POST /sessions` downstream)
- **Authentication**: Bearer JWT (Requires any valid authentication token)

#### Request Headers
| Header | Value | Required | Description |
| :--- | :--- | :--- | :--- |
| `Authorization` | `Bearer <token>` | Yes | The JWT of the authenticated system components or users. |

#### Request Body
| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `user_id` | `uuid` | Yes | The unique identifier of the user to generate a session for. |

##### Example
```json
{
  "user_id": "c3b99db1-d419-48e0-bb15-081079d38bb1"
}
```

#### Success Response
- **Status Code**: `201 Created`
- **Headers**: `Content-Type: application/json`
- **Body**:
```json
{
  "id": "e9c15ad2-f673-455b-86d7-4632b50fe732",
  "user_id": "c3b99db1-d419-48e0-bb15-081079d38bb1",
  "token": "v4.local.eyJleHAiOiIyMDI2...",
  "expired_at": "2026-05-20T11:34:05Z",
  "created_at": "2026-05-19T11:34:05Z"
}
```

#### Failure Response (Malformed Body or Missing user_id)
- **Status Code**: `400 Bad Request`
- **Body**:
```json
{
  "error": "Key: 'CreateSessionRequest.UserID' Error:Field validation for 'UserID' failed on the 'required' tag"
}
```

---

### 🟢 Revoke Session (Authenticated)

#### Endpoint Description
Explicitly terminates a session. Used when a user logs out.

- **HTTP Method**: `DELETE`
- **Path**: `/api/authentication/sessions` (routed to `DELETE /sessions` downstream)
- **Authentication**: Bearer JWT

#### Request Headers
| Header | Value | Required | Description |
| :--- | :--- | :--- | :--- |
| `Authorization` | `Bearer <token>` | Yes | The JWT of the authenticated session owner. |

#### Request Query Parameters
| Parameter | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `token` | `string` | Yes | The session token to revoke/delete. |

#### Success Response
- **Status Code**: `200 OK`
- **Body**:
```json
{
  "message": "Session revoked successfully"
}
```

#### Failure Response (Token Missing)
- **Status Code**: `400 Bad Request`
- **Body**:
```json
{
  "error": "Query parameter token is required"
}
```
