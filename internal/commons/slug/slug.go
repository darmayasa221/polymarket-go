// Package slug builds predictable Polymarket market slugs without requiring an API call.
// Slug format: "{ticker}-updown-5m-{windowStart.Unix()}"
package slug

import (
	"fmt"
	"time"

	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/commons/timeutil"
)

// ForAsset builds a SlugID for the given asset ticker and 5-minute window start.
// windowStart must already be floored to a 5-minute boundary (use timeutil.WindowStart).
func ForAsset(asset string, windowStart time.Time) polyid.SlugID {
	return polyid.SlugID(fmt.Sprintf("%s-updown-5m-%d", asset, windowStart.Unix()))
}

// CurrentWindow returns the SlugID for the active 5-minute window right now.
func CurrentWindow(asset string) polyid.SlugID {
	return ForAsset(asset, timeutil.WindowStart(timeutil.Now()))
}

// NextWindow returns the SlugID for the next 5-minute window.
func NextWindow(asset string) polyid.SlugID {
	next := timeutil.WindowStart(timeutil.Now()).Add(5 * time.Minute)
	return ForAsset(asset, next)
}
