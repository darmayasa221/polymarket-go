package constants

const (
	MaxRequestBodyBytes = 10 << 20 // 10 MB

	// Security header values.
	SecurityNoSniff        = "nosniff"
	SecurityFrameDeny      = "DENY"
	SecurityXSSProtection  = "1; mode=block"
	SecurityReferrerPolicy = "strict-origin-when-cross-origin"
)
