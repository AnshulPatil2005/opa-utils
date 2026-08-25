package apis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// makeIS will convert any number of parameters to a []interface{}
func makeIS(v ...interface{}) []interface{} {
	return v
}
func TestCompare(t *testing.T) {
	assert.Equal(t, StatusFailed, Compare(StatusFailed, StatusFailed))
	assert.Equal(t, StatusFailed, Compare(StatusFailed, StatusSkipped))
	assert.Equal(t, StatusFailed, Compare(StatusSkipped, StatusFailed))
	assert.Equal(t, StatusFailed, Compare(StatusPassed, StatusFailed))
	assert.Equal(t, StatusSkipped, Compare(StatusSkipped, StatusPassed))
	assert.Equal(t, StatusPassed, Compare(StatusPassed, StatusPassed))
	assert.Equal(t, StatusUnknown, Compare(StatusUnknown, StatusUnknown))
}

func TestCompareStatusAndSubStatus(t *testing.T) {
	assert.Equal(t, makeIS(StatusFailed, SubStatusUnknown), makeIS(CompareStatusAndSubStatus(StatusFailed, StatusPassed, SubStatusUnknown, SubStatusUnknown)))
	assert.Equal(t, makeIS(StatusFailed, SubStatusUnknown), makeIS(CompareStatusAndSubStatus(StatusFailed, StatusSkipped, SubStatusUnknown, SubStatusConfiguration)))
	assert.Equal(t, makeIS(StatusPassed, SubStatusIrrelevant), makeIS(CompareStatusAndSubStatus(StatusPassed, StatusPassed, SubStatusUnknown, SubStatusIrrelevant)))
	assert.Equal(t, makeIS(StatusPassed, SubStatusException), makeIS(CompareStatusAndSubStatus(StatusPassed, StatusPassed, SubStatusException, SubStatusUnknown)))
	assert.Equal(t, makeIS(StatusSkipped, SubStatusConfiguration), makeIS(CompareStatusAndSubStatus(StatusSkipped, StatusPassed, SubStatusConfiguration, SubStatusUnknown)))
	assert.Equal(t, makeIS(StatusSkipped, SubStatusIntegration), makeIS(CompareStatusAndSubStatus(StatusSkipped, StatusPassed, SubStatusIntegration, SubStatusUnknown)))
	assert.Equal(t, makeIS(StatusSkipped, SubStatusManualReview), makeIS(CompareStatusAndSubStatus(StatusPassed, StatusSkipped, SubStatusUnknown, SubStatusManualReview)))
	assert.Equal(t, makeIS(StatusSkipped, SubStatusRequiresReview), makeIS(CompareStatusAndSubStatus(StatusPassed, StatusSkipped, SubStatusUnknown, SubStatusRequiresReview)))

	// notEvaluated should win over other skipped substatuses since it signals a real cluster/RBAC gap
	assert.Equal(t, makeIS(StatusSkipped, SubStatusNotEvaluated), makeIS(CompareStatusAndSubStatus(StatusSkipped, StatusPassed, SubStatusNotEvaluated, SubStatusUnknown)))
	assert.Equal(t, makeIS(StatusSkipped, SubStatusNotEvaluated), makeIS(CompareStatusAndSubStatus(StatusSkipped, StatusSkipped, SubStatusNotEvaluated, SubStatusConfiguration)))
	assert.Equal(t, makeIS(StatusSkipped, SubStatusNotEvaluated), makeIS(CompareStatusAndSubStatus(StatusSkipped, StatusSkipped, SubStatusManualReview, SubStatusNotEvaluated)))
	// failed still beats notEvaluated
	assert.Equal(t, makeIS(StatusFailed, SubStatusUnknown), makeIS(CompareStatusAndSubStatus(StatusFailed, StatusSkipped, SubStatusUnknown, SubStatusNotEvaluated)))
}

func TestConvertStatusToNewStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   ScanningStatus
		expected ScanningStatus
		sub      ScanningSubStatus
	}{
		{
			name:     "StatusExcluded",
			status:   StatusExcluded,
			expected: StatusPassed,
			sub:      SubStatusException,
		},
		{
			name:     "StatusIrrelevant",
			status:   StatusIrrelevant,
			expected: StatusPassed,
			sub:      SubStatusIrrelevant,
		},
		{
			name:     "StatusPassed",
			status:   StatusPassed,
			expected: StatusPassed,
			sub:      SubStatusUnknown,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualStatus, actualSub := ConvertStatusToNewStatus(test.status)
			if actualStatus != test.expected {
				t.Errorf("Expected status %s, but got %s", test.expected, actualStatus)
			}
			if actualSub != test.sub {
				t.Errorf("Expected sub status %s, but got %s", test.sub, actualSub)
			}
		})
	}
}

func TestSubStatusInfo(t *testing.T) {
	tests := []struct {
		name      string
		subStatus ScanningSubStatus
		want      string
	}{
		{
			name:      "configuration",
			subStatus: SubStatusConfiguration,
			want:      string(SubStatusConfigurationInfo),
		},
		{
			name:      "requires review",
			subStatus: SubStatusRequiresReview,
			want:      string(SubStatusRequiresReviewInfo),
		},
		{
			name:      "manual review",
			subStatus: SubStatusManualReview,
			want:      string(SubStatusManualReviewInfo),
		},
		{
			name:      "not evaluated",
			subStatus: SubStatusNotEvaluated,
			want:      string(SubStatusNotEvaluatedInfo),
		},
		{
			name:      "unknown",
			subStatus: SubStatusUnknown,
			want:      "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := SubStatusInfo(test.subStatus); got != test.want {
				t.Errorf("SubStatusInfo() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestCompareStatusAndSubStatus_FailedCarriesException(t *testing.T) {
	// An alertOnly exception keeps the finding failing while annotating it, so the
	// exception sub status has to survive aggregation onto a failed status.
	status, subStatus := CompareStatusAndSubStatus(StatusFailed, StatusPassed, SubStatusException, SubStatusUnknown)
	assert.Equal(t, StatusFailed, status)
	assert.Equal(t, SubStatusException, subStatus)

	// A plain failure alongside an acknowledged one still reports the annotation.
	status, subStatus = CompareStatusAndSubStatus(StatusFailed, StatusFailed, SubStatusUnknown, SubStatusException)
	assert.Equal(t, StatusFailed, status)
	assert.Equal(t, SubStatusException, subStatus)

	// A failure with no exception at all carries no sub status.
	status, subStatus = CompareStatusAndSubStatus(StatusFailed, StatusPassed, SubStatusUnknown, SubStatusUnknown)
	assert.Equal(t, StatusFailed, status)
	assert.Equal(t, SubStatusUnknown, subStatus)

	// Sub statuses that belong to other statuses are not dragged onto a failure.
	status, subStatus = CompareStatusAndSubStatus(StatusFailed, StatusSkipped, SubStatusUnknown, SubStatusNotEvaluated)
	assert.Equal(t, StatusFailed, status)
	assert.Equal(t, SubStatusUnknown, subStatus)

	// Unknown stays unannotated.
	status, subStatus = CompareStatusAndSubStatus(StatusUnknown, StatusUnknown, SubStatusException, SubStatusUnknown)
	assert.Equal(t, StatusUnknown, status)
	assert.Equal(t, SubStatusUnknown, subStatus)
}
