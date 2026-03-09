package bcrypt

// DefaultCost is the default bcrypt work factor used when none is configured.
const DefaultCost = 12

// Config holds the configuration for the bcrypt encryption adapter.
type Config struct {
	// Cost is the bcrypt work factor. Default: 12. Range: 4–31.
	Cost int
}
