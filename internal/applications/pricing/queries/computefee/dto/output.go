package dto

import "github.com/darmayasa221/polymarket-go/internal/applications/shared/feecalc"

// Output wraps the FeeResult from ComputeFee.
type Output struct {
	Fee feecalc.FeeResult
}
