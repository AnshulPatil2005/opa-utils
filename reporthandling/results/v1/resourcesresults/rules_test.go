package resourcesresults

import (
	"testing"

	"github.com/armosec/armoapi-go/armotypes"
	apisv1 "github.com/kubescape/opa-utils/reporthandling/apis"
	helpersv1 "github.com/kubescape/opa-utils/reporthandling/helpers/v1"
	"github.com/stretchr/testify/assert"
)

func TestSetGetRuleName(t *testing.T) {
	r := ResourceAssociatedRule{}
	id := "my-rule"
	r.SetName(id)
	assert.Equal(t, id, r.GetName())
}

func TestRuleStatusIsFrameworkScoped(t *testing.T) {
	exception := armotypes.PostureExceptionPolicy{
		PosturePolicies: []armotypes.PosturePolicy{
			{
				FrameworkName: "NSA",
				ControlID:     "C-0034",
				RuleName:      "R1",
			},
		},
	}

	newRule := func() ResourceAssociatedRule {
		return ResourceAssociatedRule{
			Name:      "R1",
			Status:    apisv1.StatusFailed,
			Exception: []armotypes.PostureExceptionPolicy{exception},
		}
	}

	t.Run("NSA exception does not suppress MITRE failure", func(t *testing.T) {
		rule := newRule()

		rule.SetStatus(apisv1.StatusFailed, &helpersv1.Filters{
			FrameworkNames: []string{"NSA"},
		})

		assert.Equal(t, apisv1.StatusPassed, rule.GetStatus(&helpersv1.Filters{
			FrameworkNames: []string{"NSA"},
		}).Status())
		assert.Equal(t, apisv1.StatusFailed, rule.GetStatus(nil).Status())
		assert.Empty(t, rule.GetSubStatus())

		assert.Equal(t, apisv1.StatusFailed, rule.GetStatus(&helpersv1.Filters{
			FrameworkNames: []string{"MITRE"},
		}).Status())
	})

	t.Run("framework evaluation order does not change result", func(t *testing.T) {
		rule := newRule()

		mitreStatus := rule.GetStatus(&helpersv1.Filters{
			FrameworkNames: []string{"MITRE"},
		}).Status()

		rule.SetStatus(apisv1.StatusFailed, &helpersv1.Filters{
			FrameworkNames: []string{"NSA"},
		})

		nsaStatus := rule.GetStatus(&helpersv1.Filters{
			FrameworkNames: []string{"NSA"},
		}).Status()

		mitreStatusAfterNSA := rule.GetStatus(&helpersv1.Filters{
			FrameworkNames: []string{"MITRE"},
		}).Status()

		assert.Equal(t, apisv1.StatusFailed, mitreStatus)
		assert.Equal(t, apisv1.StatusPassed, nsaStatus)
		assert.Equal(t, apisv1.StatusFailed, mitreStatusAfterNSA)
	})
}
