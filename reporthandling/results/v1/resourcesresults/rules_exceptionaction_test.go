package resourcesresults

import (
	"testing"

	"github.com/armosec/armoapi-go/armotypes"
	"github.com/kubescape/opa-utils/reporthandling/apis"
	"github.com/stretchr/testify/assert"
)

func policyWithActions(name string, actions ...armotypes.PostureExceptionPolicyActions) armotypes.PostureExceptionPolicy {
	return armotypes.PostureExceptionPolicy{
		PortalBase: armotypes.PortalBase{Name: name},
		Actions:    actions,
	}
}

func failedRuleWithExceptions(exceptions ...armotypes.PostureExceptionPolicy) *ResourceAssociatedRule {
	return &ResourceAssociatedRule{
		Name:      "rule",
		Status:    apis.StatusFailed,
		Exception: exceptions,
	}
}

func TestRuleGetStatusHonoursExceptionAction(t *testing.T) {
	tests := []struct {
		name              string
		exceptions        []armotypes.PostureExceptionPolicy
		expectedStatus    apis.ScanningStatus
		expectedSubStatus apis.ScanningSubStatus
	}{
		{
			name:              "no exception leaves the failure untouched",
			exceptions:        nil,
			expectedStatus:    apis.StatusFailed,
			expectedSubStatus: apis.SubStatusUnknown,
		},
		{
			name:              "disable suppresses the finding",
			exceptions:        []armotypes.PostureExceptionPolicy{policyWithActions("a", armotypes.Disable)},
			expectedStatus:    apis.StatusPassed,
			expectedSubStatus: apis.SubStatusException,
		},
		{
			name:              "alertOnly acknowledges the finding but keeps it failing",
			exceptions:        []armotypes.PostureExceptionPolicy{policyWithActions("a", armotypes.AlertOnly)},
			expectedStatus:    apis.StatusFailed,
			expectedSubStatus: apis.SubStatusException,
		},
		{
			name: "every exception alertOnly keeps the finding failing",
			exceptions: []armotypes.PostureExceptionPolicy{
				policyWithActions("a", armotypes.AlertOnly),
				policyWithActions("b", armotypes.AlertOnly),
			},
			expectedStatus:    apis.StatusFailed,
			expectedSubStatus: apis.SubStatusException,
		},
		{
			name: "a single disable among alertOnly exceptions suppresses the finding",
			exceptions: []armotypes.PostureExceptionPolicy{
				policyWithActions("a", armotypes.AlertOnly),
				policyWithActions("b", armotypes.Disable),
			},
			expectedStatus:    apis.StatusPassed,
			expectedSubStatus: apis.SubStatusException,
		},
		{
			name: "a policy carrying both actions counts as disable",
			exceptions: []armotypes.PostureExceptionPolicy{
				policyWithActions("a", armotypes.AlertOnly, armotypes.Disable),
			},
			expectedStatus:    apis.StatusPassed,
			expectedSubStatus: apis.SubStatusException,
		},
		{
			name:              "a policy with no actions keeps the historical suppressing behaviour",
			exceptions:        []armotypes.PostureExceptionPolicy{policyWithActions("a")},
			expectedStatus:    apis.StatusPassed,
			expectedSubStatus: apis.SubStatusException,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := failedRuleWithExceptions(tt.exceptions...).GetStatus(nil)
			assert.Equal(t, tt.expectedStatus, status.Status())
			assert.Equal(t, tt.expectedSubStatus, status.GetSubStatus())
		})
	}
}

func TestRuleGetStatusLeavesNonFailedRulesAlone(t *testing.T) {
	// A rule that did not fail is never re-stated by an exception, whatever the action.
	rule := &ResourceAssociatedRule{
		Name:      "rule",
		Status:    apis.StatusPassed,
		Exception: []armotypes.PostureExceptionPolicy{policyWithActions("a", armotypes.AlertOnly)},
	}
	status := rule.GetStatus(nil)
	assert.Equal(t, apis.StatusPassed, status.Status())
	assert.Equal(t, apis.SubStatusUnknown, status.GetSubStatus())
}

func TestAllExceptionsAlertOnly(t *testing.T) {
	assert.False(t, allExceptionsAlertOnly(nil), "no exception is not an alertOnly acknowledgement")
	assert.False(t, allExceptionsAlertOnly([]armotypes.PostureExceptionPolicy{}))
	assert.True(t, allExceptionsAlertOnly([]armotypes.PostureExceptionPolicy{
		policyWithActions("a", armotypes.AlertOnly),
	}))
	assert.False(t, allExceptionsAlertOnly([]armotypes.PostureExceptionPolicy{
		policyWithActions("a", armotypes.AlertOnly),
		policyWithActions("b", armotypes.Disable),
	}))
}
