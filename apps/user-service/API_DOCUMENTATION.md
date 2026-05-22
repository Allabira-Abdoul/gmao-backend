# User & RBAC Service API Specification

The **User & RBAC Service** manages user accounts, user profiles, system roles, and associated capability privileges. It registers with Consul under the name `user-service` and operates on port `8082`.

All external client requests are routed through the API Gateway using the prefix `/api/user/*`.

---

## 🔒 Security & Privileges

This service enforces the privilege-based RBAC model. Users must present a Bearer JWT containing the required privilege for each endpoint:

- **User Read**: Requires `USER_VIEW`
- **User Create**: Requires `USER_CREATE`
- **User Update**: Requires `USER_UPDATE`
- **User Delete**: Requires `USER_DELETE`
- **Role Read**: Requires `ROLE_VIEW`
- **Role Create**: Requires `ROLE_CREATE`
- **Role Update**: Requires `ROLE_UPDATE`
- **Role Delete**: Requires `ROLE_DELETE`

---

## 📡 Endpoints Specification

### 🟢 Get Current User Profile (Authenticated)

Retrieve details of the currently logged-in user context.

- **HTTP Method**: `GET`
- **Path**: `/api/user/users/me` (routed to `GET /users/me` downstream)
- **Authentication**: Bearer JWT (Any active logged-in user)

#### Request Format
- **Headers**:
  - `Authorization: Bearer <token>`

#### Success Response
- **Status Code**: `200 OK`
- **Body**:
```json
{
  "id": "c3b99db1-d419-48e0-bb15-081079d38bb1",
  "full_name": "Auditor User",
  "email": "auditor@gmao.com",
  "status": "ACTIVE",
  "role": {
    "id": "a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d",
    "name": "Auditor",
    "description": "System auditor with logs inspection privilege",
    "privileges": [
      "USER_VIEW",
      "ROLE_VIEW",
      "SYSTEM_AUDIT_VIEW"
    ],
    "created_at": "2026-05-19T10:00:00Z",
    "updated_at": "2026-05-19T10:00:00Z"
  },
  "team": null,
  "created_at": "2026-05-19T10:15:00Z",
  "updated_at": "2026-05-19T10:15:00Z"
}
```

---

### 🟢 List Users (Privilege: `USER_VIEW`)

Paginated list of all users in the system.

- **HTTP Method**: `GET`
- **Path**: `/api/user/users`
- **Authentication**: Bearer JWT

#### Request Query Parameters
| Parameter | Type | Required | Default | Description |
| :--- | :--- | :--- | :--- | :--- |
| `page` | `integer` | No | `1` | The page number to fetch. |
| `per_page` | `integer` | No | `20` | The page size (limit). |

#### Success Response
- **Status Code**: `200 OK`
- **Body**:
```json
[
  {
    "id": "c3b99db1-d419-48e0-bb15-081079d38bb1",
    "full_name": "Auditor User",
    "email": "auditor@gmao.com",
    "status": "ACTIVE",
    "role": {
      "id": "a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d",
      "name": "Auditor",
      "description": "System auditor with logs inspection privilege",
      "privileges": ["USER_VIEW", "ROLE_VIEW", "SYSTEM_AUDIT_VIEW"],
      "created_at": "2026-05-19T10:00:00Z",
      "updated_at": "2026-05-19T10:00:00Z"
    },
    "team": null,
    "created_at": "2026-05-19T10:15:00Z",
    "updated_at": "2026-05-19T10:15:00Z"
  }
]
```

---

### 🟢 Get Specific User (Privilege: `USER_VIEW`)

Retrieve user record by unique UUID.

- **HTTP Method**: `GET`
- **Path**: `/api/user/users/:id`
- **Authentication**: Bearer JWT

#### Success Response
- **Status Code**: `200 OK`
- **Body**: (Similar structure to `GET /users/me`)

#### Failure Response (Not Found)
- **Status Code**: `404 Not Found`
- **Body**:
```json
{
  "error": "User not found"
}
```

---

### 🟢 Create User (Privilege: `USER_CREATE`)

Register a new user into the system.

- **HTTP Method**: `POST`
- **Path**: `/api/user/users`
- **Authentication**: Bearer JWT

#### Request Body
| Field | Type | Required | Validation | Description |
| :--- | :--- | :--- | :--- | :--- |
| `full_name` | `string` | Yes | Min 2, max 255 chars | Full name of the user. |
| `email` | `string` | Yes | Valid email syntax | System login email. |
| `password` | `string` | Yes | Min 8 chars | Plain-text password. |
| `role_id` | `string` | Yes | Valid UUID | Target role UUID. |

##### Example
```json
{
  "full_name": "Technician Joe",
  "email": "joe@gmao.com",
  "password": "securepassword123",
  "role_id": "b2f689e4-cc79-4d6d-85fa-7f8976a45612"
}
```

#### Success Response
- **Status Code**: `201 Created`
- **Body**: (Returns the newly created user without password)

#### Failure Response (Email Already Registered)
- **Status Code**: `409 Conflict`
- **Body**:
```json
{
  "error": "email already registered"
}
```

---

### 🟢 Update User (Privilege: `USER_UPDATE`)

Update an existing user's details, status, role, or team assignment.

- **HTTP Method**: `PUT`
- **Path**: `/api/user/users/:id`
- **Authentication**: Bearer JWT

#### Request Body
| Field | Type | Required | Validation / Allowed Values | Description |
| :--- | :--- | :--- | :--- | :--- |
| `full_name` | `string` | No | Min 2, max 255 chars | Updated full name. |
| `email` | `string` | No | Valid email syntax | Updated login email. |
| `status` | `string` | No | `ACTIVE`, `INACTIVE`, `LOCKED` | Account status. |
| `role_id` | `string` | No | Valid UUID | Updated role. |
| `team_id` | `string` | No | Valid UUID | Updated team assignment. |

---

### 🟢 Delete User (Privilege: `USER_DELETE`)

Permanently remove a user account.

- **HTTP Method**: `DELETE`
- **Path**: `/api/user/users/:id`
- **Authentication**: Bearer JWT

#### Success Response
- **Status Code**: `200 OK`
- **Body**:
```json
{
  "message": "User deleted successfully"
}
```

---

## 👥 Dynamic Role Management Endpoints

### 🟢 List Roles (Privilege: `ROLE_VIEW`)
- **HTTP Method**: `GET`
- **Path**: `/api/user/roles`

### 🟢 Get Specific Role (Privilege: `ROLE_VIEW`)
- **HTTP Method**: `GET`
- **Path**: `/api/user/roles/:id`

### 🟢 Create Role (Privilege: `ROLE_CREATE`)
- **HTTP Method**: `POST`
- **Path**: `/api/user/roles`
- **Request Body**:
```json
{
  "name": "Custom Technician",
  "description": "Custom role for third-party service provider",
  "privileges": ["ASSET_VIEW", "WORKORDER_VIEW", "WORKORDER_UPDATE"]
}
```

### 🟢 Update Role (Privilege: `ROLE_UPDATE`)
- **HTTP Method**: `PUT`
- **Path**: `/api/user/roles/:id`

### 🟢 Delete Role (Privilege: `ROLE_DELETE`)
- **HTTP Method**: `DELETE`
- **Path**: `/api/user/roles/:id`

### 🟢 Get System Privileges List (Public)
Retrieve the complete hardcoded vocabulary of capability keys.
- **HTTP Method**: `GET`
- **Path**: `/api/user/privileges`
- **Success Response**: `200 OK` with JSON array of strings e.g. `["USER_VIEW", "USER_CREATE", ...]`

---

## 🔒 Internal Service-to-Service Endpoints

These endpoints are bypass-routed but blocked from public gateway access. They require `X-Internal-Service: true` header to run.

### 🟢 Get User By Email (Internal)
Retrieve complete user payload, including the hashed password, for validation checks.
- **HTTP Method**: `GET`
- **Path**: `/internal/by-email?email=...`
- **Headers**: `X-Internal-Service: true`

### 🟢 Get User By ID (Internal)
Retrieve complete user payload, including structural privilege list, for validation checks.
- **HTTP Method**: `GET`
- **Path**: `/internal/by-id?id=...`
- **Headers**: `X-Internal-Service: true`
