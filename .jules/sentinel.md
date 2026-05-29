## 2024-05-29 - Hardcoded Default JWT Secret Fallback
**Vulnerability:** JWT secrets had insecure hardcoded fallback values in environment variable reads across multiple services.
**Learning:** Default arguments in `getEnv` pattern bypassed explicit configuration requirements, potentially deploying to production with known compromised development keys.
**Prevention:** Remove hardcoded defaults for critical security keys. Use `os.Getenv` and explicitly trigger `log.Fatal` if the environment variable is missing to force secure configuration at runtime.
