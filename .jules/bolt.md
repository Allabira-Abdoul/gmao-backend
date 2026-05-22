## 2024-05-18 - Fix N+1 Query in GetAllWorkOrders
**Learning:** Found an N+1 query issue where the application service layer iterated over all retrieved work orders and performed an individual query to fetch `Interventions` for each.
**Action:** Utilized GORM’s `.Preload("Interventions").Preload("Interventions.Measurements")` at the repository layer to fetch all nested entities in a single database round-trip, drastically improving read efficiency and avoiding N+1 queries. Always prefer database-level joins/preloading over application-level loops for data fetching.
