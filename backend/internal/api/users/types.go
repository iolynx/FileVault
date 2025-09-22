package users

// User represents a user in the system for API responses.
// ID: unique identifier of the user.
// Email: user's email address.
// Name: user's display name.
type User struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// MeResponse represents the detailed user information returned by the /auth/me endpoint.
// ID: unique identifier of the user.
// Email: user's email address.
// Name: user's display name.
// Role: user's role in the system.
// StorageUsedBytes: total raw storage used by the user (before deduplication).
// DeduplicatedUsageBytes: total storage used after deduplication.
// StorageQuotaBytes: the user's assigned storage quota.
// SavingsBytes: total storage saved due to deduplication.
// SavingsPercentage: percentage of storage saved compared to raw usage.
type MeResponse struct {
	ID                     int64   `json:"id"`
	Email                  string  `json:"email"`
	Name                   string  `json:"name"`
	Role                   string  `json:"role"`
	StorageUsedBytes       int64   `json:"storage_used_bytes"`       // "Original storage usage"
	DeduplicatedUsageBytes int64   `json:"deduplicated_usage_bytes"` // "Total storage used (deduplicated)"
	StorageQuotaBytes      int64   `json:"storage_quota_bytes"`
	SavingsBytes           int64   `json:"savings_bytes"`
	SavingsPercentage      float64 `json:"savings_percentage"`
}

// signupRequest represents the expected JSON payload for user signup.
// Contains email, name, and plaintext password (hashed before storage)
type signupRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// loginRequest represents the expected JSON payload for user login.
// Contains email and plaintext password
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
