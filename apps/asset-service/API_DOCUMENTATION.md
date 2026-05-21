# Asset Service API Specification

The **Asset Service** manages the lifecycle of physical machinery, technical equipment, and operational inventory assets. It registers with Consul under the name `asset-service` and operates on port `8083`.

All external client requests are routed through the API Gateway using the prefix `/api/asset/*`.

---

## 🔒 Security & Privileges

Endpoints require user authentication (Bearer JWT) with the appropriate dynamic privilege:

- **Create Asset**: Requires `ASSET_CREATE`
- **List Assets & Get Asset**: Requires `ASSET_VIEW`
- **Update Asset**: Requires `ASSET_UPDATE`
- **Delete Asset**: Requires `ASSET_DELETE`

---

## 📡 Endpoints Specification

### 🟢 List Assets (Privilege: `ASSET_VIEW`)

Retrieve all registered assets.

- **HTTP Method**: `GET`
- **Path**: `/api/asset/assets` (routed to `GET /assets` downstream)
- **Authentication**: Bearer JWT

#### Success Response
- **Status Code**: `200 OK`
- **Body**:
```json
[
  {
    "id": "7b09db24-f719-482a-a9bd-8311d9d93bb1",
    "name": "Industrial Boiler Model X100",
    "code": "ASSET-B01",
    "status": "OPERATIONAL",
    "category": "HEATING",
    "location": "Sector B / Building 2",
    "purchase_date": "2026-01-15T00:00:00Z",
    "purchase_value": 45000.5,
    "created_at": "2026-05-19T10:15:00Z",
    "updated_at": "2026-05-19T10:15:00Z"
  }
]
```

---

### 🟢 Get Specific Asset (Privilege: `ASSET_VIEW`)

Retrieve an asset by its unique database UUID.

- **HTTP Method**: `GET`
- **Path**: `/api/asset/assets/:id`
- **Authentication**: Bearer JWT

#### Success Response
- **Status Code**: `200 OK`
- **Body**: (Returns a single asset JSON object with the structure shown above)

#### Failure Response (Not Found / Invalid ID)
- **Status Code**: `404 Not Found`
- **Body**:
```json
{
  "error": "Asset not found"
}
```

---

### 🟢 Get Asset By Code (Privilege: `ASSET_VIEW`)

Retrieve an asset by its unique business code (e.g. `ASSET-B01`).

- **HTTP Method**: `GET`
- **Path**: `/api/asset/assets/code/:code`
- **Authentication**: Bearer JWT

#### Success Response
- **Status Code**: `200 OK`
- **Body**: (Returns the single asset JSON object)

---

### 🟢 Create Asset (Privilege: `ASSET_CREATE`)

Register a new physical asset into system inventory.

- **HTTP Method**: `POST`
- **Path**: `/api/asset/assets`
- **Authentication**: Bearer JWT

#### Request Body
| Field | Type | Required | Validation | Description |
| :--- | :--- | :--- | :--- | :--- |
| `name` | `string` | Yes | Min 2, max 255 chars | Full display name of the asset. |
| `code` | `string` | Yes | Min 2, max 50 chars | Unique business/inventory code. |
| `category` | `string` | Yes | Non-empty | Operational category (e.g. HEATING, HVAC, CNC). |
| `location` | `string` | Yes | Non-empty | Physical factory floor location. |
| `purchase_date` | `string (RFC3339)` | Yes | Valid date format | Asset purchase date. |
| `purchase_value` | `float64` | Yes | Positive float | Initial acquisition purchase cost. |

##### Example
```json
{
  "name": "Industrial Boiler Model X100",
  "code": "ASSET-B01",
  "category": "HEATING",
  "location": "Sector B / Building 2",
  "purchase_date": "2026-01-15T00:00:00Z",
  "purchase_value": 45000.50
}
```

#### Success Response
- **Status Code**: `201 Created`
- **Body**: (Returns the newly created asset response object)

---

### 🟢 Update Asset (Privilege: `ASSET_UPDATE`)

Modify an asset's state, category, or location.

- **HTTP Method**: `PUT`
- **Path**: `/api/asset/assets/:id`
- **Authentication**: Bearer JWT

#### Request Body
| Field | Type | Required | Validation / Allowed Values | Description |
| :--- | :--- | :--- | :--- | :--- |
| `name` | `string` | No | Min 2, max 255 chars | Updated display name. |
| `status` | `string` | No | `OPERATIONAL`, `DOWN`, `UNDER_REPAIR`, `SCRAPPED` | Asset state change. |
| `category` | `string` | No | Non-empty | Updated category. |
| `location` | `string` | No | Non-empty | Updated factory floor location. |
| `purchase_value` | `float64` | No | Positive float | Updated acquisition value. |

##### Example
```json
{
  "status": "DOWN",
  "location": "Repair Zone A"
}
```

#### Success Response
- **Status Code**: `200 OK`
- **Body**: (Returns the updated asset response object)

---

### 🟢 Delete Asset (Privilege: `ASSET_DELETE`)

Permanently remove/de-register an asset.

- **HTTP Method**: `DELETE`
- **Path**: `/api/asset/assets/:id`
- **Authentication**: Bearer JWT

#### Success Response
- **Status Code**: `200 OK`
- **Body**:
```json
{
  "message": "Asset deleted successfully"
}
```
