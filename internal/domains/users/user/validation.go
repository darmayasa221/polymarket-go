package user

// Validation limits for User fields.
// Change these constants to adjust validation rules — no other file changes needed.
const (
	UsernameMinLength = 3
	UsernameMaxLength = 50
	FullNameMinLength = 2
	FullNameMaxLength = 100
	EmailMaxLength    = 254 // RFC 5321 limit
)
