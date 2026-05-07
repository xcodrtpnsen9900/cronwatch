// Package fingerprint produces a stable string key that uniquely
// identifies an alert event so that deduplication and rate-limiting
// layers can correlate repeated occurrences of the same condition.
package fingerprint

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

// AlertType classifies the kind of alert being fingerprinted.
type AlertType string

const (
	Missed    AlertType = "missed"
	Failed    AlertType = "failed"
	Recovered AlertType = "recovered"
)

// Of returns a short, stable hex fingerprint for the combination of
// job name and alert type. The result is safe to use as a map key,
// cache key, or dedup token.
//
//	key := fingerprint.Of("backup-db", fingerprint.Failed)
func Of(jobName string, alertType AlertType) string {
	raw := strings.ToLower(strings.TrimSpace(jobName)) + ":" + string(alertType)
	sum := sha256.Sum256([]byte(raw))
	return fmt.Sprintf("%x", sum[:8]) // 16 hex chars – collision-resistant enough
}

// Parts reconstructs the canonical input string that was hashed.
// Useful for logging and diagnostics; it is NOT the reverse of Of.
func Parts(jobName string, alertType AlertType) string {
	return strings.ToLower(strings.TrimSpace(jobName)) + ":" + string(alertType)
}
