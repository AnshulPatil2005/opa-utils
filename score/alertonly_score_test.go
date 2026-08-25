package score

import (
	"testing"

	"github.com/kubescape/opa-utils/reporthandling/apis"
	"github.com/kubescape/opa-utils/reporthandling/results/v1/reportsummary"
	"github.com/stretchr/testify/assert"
)

// controlWithResources builds a control summary the way updateControlsSummaryCounters
// does, one resource status at a time.
func controlWithResources(t *testing.T, statuses ...*apis.StatusInfo) *reportsummary.ControlSummary {
	t.Helper()
	control := &reportsummary.ControlSummary{ControlID: "C-0034"}
	for i, status := range statuses {
		control.Append(status, resourceID(i))
	}
	control.CalculateStatus()
	return control
}

func resourceID(i int) string {
	return string(rune('a' + i))
}

// The compliance score is passed/all. An exception with a disable action moves a
// resource to the passed side and lifts the score; an alertOnly exception must not,
// because it only acknowledges the finding.
func TestComplianceScore_AlertOnlyDoesNotLiftTheScore(t *testing.T) {
	su := NewScore(nil)

	failed := &apis.StatusInfo{InnerStatus: apis.StatusFailed}
	suppressed := &apis.StatusInfo{InnerStatus: apis.StatusPassed, SubStatus: apis.SubStatusException}
	acknowledged := &apis.StatusInfo{InnerStatus: apis.StatusFailed, SubStatus: apis.SubStatusException}

	baseline := controlWithResources(t, failed, failed)
	assert.Equal(t, float32(0), su.GetControlComplianceScore(baseline, ""))

	withDisable := controlWithResources(t, suppressed, failed)
	assert.Equal(t, float32(50), su.GetControlComplianceScore(withDisable, ""),
		"a disable exception removes the finding, so the score rises")

	withAlertOnly := controlWithResources(t, acknowledged, failed)
	assert.Equal(t, float32(0), su.GetControlComplianceScore(withAlertOnly, ""),
		"an alertOnly exception acknowledges the finding, so the score must not move")
}
