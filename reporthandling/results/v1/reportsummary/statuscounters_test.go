package reportsummary

import (
	"testing"

	"github.com/kubescape/opa-utils/reporthandling/apis"
	"github.com/stretchr/testify/assert"
)

var resourcesCounter = StatusCounters{}

// =================================== Counters ============================================

func setResourcesCountersMock() {
	resourcesCounter.PassedResources = 7
	resourcesCounter.FailedResources = 15
	resourcesCounter.SkippedResources = 5
}

func TestSet(t *testing.T) {
	rc := StatusCounters{}
	summaryDetails := MockSummaryDetails()
	rc.Set(summaryDetails.ListResourcesIDs(nil))
	assert.Equal(t, 1, rc.SkippedResources)
	assert.Equal(t, 1, rc.PassedResources)
	assert.Equal(t, 2, rc.FailedResources)
}

// Excluded get the number of skipped resources
func TestExcluded(t *testing.T) {
	setResourcesCountersMock()
	assert.Equal(t, 5, resourcesCounter.Skipped())
}

// Passed get the number of passed resources
func TestPassed(t *testing.T) {
	setResourcesCountersMock()
	assert.Equal(t, 7, resourcesCounter.Passed())
}

// Failed get the number of failed resources
func TestFailed(t *testing.T) {
	setResourcesCountersMock()
	assert.Equal(t, 15, resourcesCounter.Failed())
}

// NumberOfAll get the number of all resources
func TestNumberOfAll(t *testing.T) {
	setResourcesCountersMock()
	assert.Equal(t, 27, resourcesCounter.All())
}

func TestStatusCounters_Increase(t *testing.T) {
	tests := []struct {
		resourceCounters *StatusCounters
		name             string
		status           apis.ScanningStatus
		expectedPassed   int
		expectedFailed   int
		expectedSkipped  int
		expectedExcluded int
	}{
		{
			name:   "Test passed status",
			status: apis.StatusPassed,
			resourceCounters: &StatusCounters{
				FailedResources:  1,
				SkippedResources: 2,
				PassedResources:  3,
			},
			expectedFailed:  1,
			expectedSkipped: 2,
			expectedPassed:  4,
		},
		{
			name:   "Test failed status",
			status: apis.StatusFailed,
			resourceCounters: &StatusCounters{
				FailedResources:  1,
				SkippedResources: 2,
				PassedResources:  3,
			},
			expectedFailed:  2,
			expectedSkipped: 2,
			expectedPassed:  3,
		},
		{
			name:   "Test skipped status",
			status: apis.StatusSkipped,
			resourceCounters: &StatusCounters{
				FailedResources:  1,
				SkippedResources: 2,
				PassedResources:  3,
			},
			expectedFailed:  1,
			expectedSkipped: 3,
			expectedPassed:  3,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			status := &apis.StatusInfo{InnerStatus: test.status}
			test.resourceCounters.Increase(status)
			if test.resourceCounters.PassedResources != test.expectedPassed {
				t.Errorf("Expected PassedResources to be %d, but got %d", test.expectedPassed, test.resourceCounters.PassedResources)
			}
			if test.resourceCounters.FailedResources != test.expectedFailed {
				t.Errorf("Expected FailedResources to be %d, but got %d", test.expectedFailed, test.resourceCounters.FailedResources)
			}
			if test.resourceCounters.SkippedResources != test.expectedSkipped {
				t.Errorf("Expected SkippedResources to be %d, but got %d", test.expectedSkipped, test.resourceCounters.SkippedResources)
			}
		})
	}
}

func TestSubStatusCounters_Increase(t *testing.T) {
	tests := []struct {
		subStatusCounters *SubStatusCounters
		name              string
		status            apis.ScanningStatus
		subStatus         apis.ScanningSubStatus
		expectedIgnored   int
	}{
		{
			name:      "Test ignored and passed status",
			status:    apis.StatusPassed,
			subStatus: apis.SubStatusException,
			subStatusCounters: &SubStatusCounters{
				IgnoredResources: 1,
			},
			expectedIgnored: 2,
		},
		{
			name:      "Test ignored and failed status",
			status:    apis.StatusFailed,
			subStatus: apis.SubStatusException,
			subStatusCounters: &SubStatusCounters{
				IgnoredResources: 1,
			},
			expectedIgnored: 1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			status := &apis.StatusInfo{InnerStatus: test.status, SubStatus: test.subStatus}
			test.subStatusCounters.Increase(status)
			if test.subStatusCounters.IgnoredResources != test.expectedIgnored {
				t.Errorf("Expected IgnoredResources to be %d, but got %d", test.expectedIgnored, test.subStatusCounters.IgnoredResources)
			}
		})
	}
}

func TestSubStatusCounters_IgnoredVsAcknowledged(t *testing.T) {
	counters := &SubStatusCounters{}

	// disable exception: suppressed, counted as ignored
	counters.Increase(&apis.StatusInfo{InnerStatus: apis.StatusPassed, SubStatus: apis.SubStatusException})
	// alertOnly exception: still failing, counted as acknowledged
	counters.Increase(&apis.StatusInfo{InnerStatus: apis.StatusFailed, SubStatus: apis.SubStatusException})
	counters.Increase(&apis.StatusInfo{InnerStatus: apis.StatusFailed, SubStatus: apis.SubStatusException})
	// unrelated statuses do not move either counter
	counters.Increase(&apis.StatusInfo{InnerStatus: apis.StatusFailed})
	counters.Increase(&apis.StatusInfo{InnerStatus: apis.StatusPassed})
	counters.Increase(&apis.StatusInfo{InnerStatus: apis.StatusSkipped, SubStatus: apis.SubStatusNotEvaluated})

	assert.Equal(t, 1, counters.Ignored())
	assert.Equal(t, 2, counters.Acknowledged())
	assert.Equal(t, 3, counters.All())
}

func TestControlSummaryAppend_AlertOnlyStillCountsAsFailed(t *testing.T) {
	// The compliance score is passed/all, so an acknowledged resource has to stay on
	// the failed side of the counters for the score to reflect the real state.
	acknowledged := &ControlSummary{}
	acknowledged.Append(
		&apis.StatusInfo{InnerStatus: apis.StatusFailed, SubStatus: apis.SubStatusException},
		"resource-1",
	)

	statuses, subStatuses := acknowledged.StatusesCounters()
	assert.Equal(t, 1, statuses.Failed())
	assert.Equal(t, 0, statuses.Passed())
	assert.Equal(t, 1, subStatuses.Acknowledged())
	assert.Equal(t, 0, subStatuses.Ignored())

	suppressed := &ControlSummary{}
	suppressed.Append(
		&apis.StatusInfo{InnerStatus: apis.StatusPassed, SubStatus: apis.SubStatusException},
		"resource-1",
	)

	statuses, subStatuses = suppressed.StatusesCounters()
	assert.Equal(t, 0, statuses.Failed())
	assert.Equal(t, 1, statuses.Passed())
	assert.Equal(t, 0, subStatuses.Acknowledged())
	assert.Equal(t, 1, subStatuses.Ignored())
}

// mirrors updateControlsSummaryCounters: append the resource status, then fold in
// the resource-level sub status.
func appendResource(control *ControlSummary, status *apis.StatusInfo, id string) {
	control.Append(status, id)
	control.calculateStatus(status.GetSubStatus())
}

func TestControlSummary_AlertOnlyKeepsControlFailingAndAnnotated(t *testing.T) {
	control := &ControlSummary{ControlID: "C-0034"}
	appendResource(control, &apis.StatusInfo{InnerStatus: apis.StatusFailed, SubStatus: apis.SubStatusException}, "a")

	assert.Equal(t, apis.StatusFailed, control.GetStatus().Status())
	assert.Equal(t, apis.SubStatusException, control.GetStatus().GetSubStatus(),
		"the report must be able to tell an accepted risk apart from an unreviewed failure")
}

func TestControlSummary_PlainFailureStaysUnannotated(t *testing.T) {
	control := &ControlSummary{ControlID: "C-0034"}
	appendResource(control, &apis.StatusInfo{InnerStatus: apis.StatusFailed}, "a")

	assert.Equal(t, apis.StatusFailed, control.GetStatus().Status())
	assert.Equal(t, apis.SubStatusUnknown, control.GetStatus().GetSubStatus())
}

func TestControlSummary_MixedFailuresReportAcknowledgedBreakdown(t *testing.T) {
	control := &ControlSummary{ControlID: "C-0034"}
	appendResource(control, &apis.StatusInfo{InnerStatus: apis.StatusFailed}, "a")
	appendResource(control, &apis.StatusInfo{InnerStatus: apis.StatusFailed, SubStatus: apis.SubStatusException}, "b")

	assert.Equal(t, apis.StatusFailed, control.GetStatus().Status())
	assert.Equal(t, apis.SubStatusException, control.GetStatus().GetSubStatus())

	statuses, subStatuses := control.StatusesCounters()
	assert.Equal(t, 2, statuses.Failed(), "both resources still fail")
	assert.Equal(t, 1, subStatuses.Acknowledged(), "only one of them is an accepted risk")
}
