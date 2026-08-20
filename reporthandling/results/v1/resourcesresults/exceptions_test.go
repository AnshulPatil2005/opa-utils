package resourcesresults

import (
	"github.com/armosec/armoapi-go/identifiers"
	"testing"

	"github.com/armosec/armoapi-go/armotypes"
	"github.com/kubescape/k8s-interface/workloadinterface"
	"github.com/kubescape/opa-utils/exceptions"
	"github.com/kubescape/opa-utils/reporthandling"
	"github.com/kubescape/opa-utils/reporthandling/apis"
	helpersv1 "github.com/kubescape/opa-utils/reporthandling/helpers/v1"
	"github.com/stretchr/testify/assert"
)

func mockExceptionDeploymentC0087() *armotypes.PostureExceptionPolicy {
	return &armotypes.PostureExceptionPolicy{
		PortalBase: armotypes.PortalBase{
			Name: "DeploymentC0087",
		},
		Actions: []armotypes.PostureExceptionPolicyActions{armotypes.AlertOnly},
		Resources: []identifiers.PortalDesignator{
			{
				DesignatorType: identifiers.DesignatorAttributes,
				Attributes: map[string]string{
					identifiers.AttributeKind: "Deployment",
				},
			},
		},
		PosturePolicies: []armotypes.PosturePolicy{
			{
				ControlID: "C-0087",
			},
		},
	}
}

func mockExceptionUnitestDeploymentC0087() *armotypes.PostureExceptionPolicy {
	return &armotypes.PostureExceptionPolicy{
		PortalBase: armotypes.PortalBase{
			Name: "unitestDeploymentC0087",
		},
		Actions: []armotypes.PostureExceptionPolicyActions{armotypes.AlertOnly},
		Resources: []identifiers.PortalDesignator{
			{
				DesignatorType: identifiers.DesignatorAttributes,
				Attributes: map[string]string{
					identifiers.AttributeCluster: "unitest",
					identifiers.AttributeKind:    "Deployment",
				},
			},
		},
		PosturePolicies: []armotypes.PosturePolicy{
			{
				ControlID: "C-0087",
			},
		},
	}
}

func mockExceptionUnitestC0088() *armotypes.PostureExceptionPolicy {
	return &armotypes.PostureExceptionPolicy{
		PortalBase: armotypes.PortalBase{
			Name: "unitestC0088",
		},
		Actions: []armotypes.PostureExceptionPolicyActions{armotypes.AlertOnly},
		Resources: []identifiers.PortalDesignator{
			{
				DesignatorType: identifiers.DesignatorAttributes,
				Attributes: map[string]string{
					identifiers.AttributeCluster: "unitest",
				},
			},
		},
		PosturePolicies: []armotypes.PosturePolicy{
			{
				ControlID: "C-0088",
			},
		},
	}
}

func mockExceptionDeploymentC0089() *armotypes.PostureExceptionPolicy {
	return &armotypes.PostureExceptionPolicy{
		PortalBase: armotypes.PortalBase{
			Name: "Deployment0089",
		},
		Actions: []armotypes.PostureExceptionPolicyActions{armotypes.AlertOnly},
		Resources: []identifiers.PortalDesignator{
			{
				DesignatorType: identifiers.DesignatorAttributes,
				Attributes: map[string]string{
					identifiers.AttributeKind: "Deployment",
				},
			},
		},
		PosturePolicies: []armotypes.PosturePolicy{
			{
				ControlID: "C-0089",
			},
		},
	}
}

func mockControlsList() map[string]reporthandling.Control {
	return map[string]reporthandling.Control{
		"C-0087": {},
		"C-0088": {},
		"C-0089": {},
	}
}

func TestSetExceptions(t *testing.T) {
	w := workloadinterface.NewWorkloadMock(nil)
	processor := exceptions.NewProcessor()

	exceptions := []armotypes.PostureExceptionPolicy{}
	exceptions = append(exceptions, *mockExceptionDeploymentC0087())
	exceptions = append(exceptions, *mockExceptionUnitestDeploymentC0087())
	exceptions = append(exceptions, *mockExceptionUnitestC0088())
	exceptions = append(exceptions, *mockExceptionDeploymentC0089())
	c := mockControlsList()
	// simple test
	result1 := mockResultFailed()
	result1.SetExceptions(w, exceptions, "", c, WithExceptionsProcessor(processor))
	assert.Equal(t, 2, result1.ListControlsIDs(nil).Passed())
	assert.Equal(t, 1, result1.ListControlsIDs(nil).Failed())

	// without option to reuse the processor
	result1.SetExceptions(w, exceptions, "", c)
	assert.Equal(t, 2, result1.ListControlsIDs(nil).Passed())
	assert.Equal(t, 1, result1.ListControlsIDs(nil).Failed())

	// test cluster name
	result2 := mockResultFailed()
	result2.SetExceptions(w, exceptions, "unitest", c, WithExceptionsProcessor(processor))
	assert.Equal(t, 3, result2.ListControlsIDs(nil).Passed())
	assert.Equal(t, 0, result2.ListControlsIDs(nil).Failed())

	// test wrong cluster name
	result3 := mockResultFailed()
	result3.SetExceptions(w, exceptions, "unitest2", c, WithExceptionsProcessor(processor))
	assert.Equal(t, 2, result3.ListControlsIDs(nil).Passed())
	assert.Equal(t, 1, result3.ListControlsIDs(nil).Failed())
}

func TestSetExceptionsKeepsRuleStatusForFrameworkScopedEvaluation(t *testing.T) {
	w := workloadinterface.NewWorkloadMock(nil)
	policy := mockExceptionDeploymentC0087()
	policy.PosturePolicies[0].FrameworkName = "NSA"
	policy.PosturePolicies[0].RuleName = "ruleA"

	result := Result{AssociatedControls: []ResourceAssociatedControl{{
		ControlID: "C-0087",
		Status:    apis.StatusInfo{InnerStatus: apis.StatusFailed},
		ResourceAssociatedRules: []ResourceAssociatedRule{
			*mockResourceAssociatedRuleA(),
		},
	}}}

	result.SetExceptions(w, []armotypes.PostureExceptionPolicy{*policy}, "", map[string]reporthandling.Control{"C-0087": {}})

	rule := result.AssociatedControls[0].ResourceAssociatedRules[0]
	assert.Equal(t, apis.StatusFailed, rule.Status, "exceptions must not overwrite the raw rule status")
	assert.Equal(t, apis.StatusPassed, result.AssociatedControls[0].GetStatus(nil).Status())
	assert.Equal(t, apis.SubStatusException, result.AssociatedControls[0].GetStatus(nil).GetSubStatus())
	assert.Equal(t, apis.StatusPassed, result.AssociatedControls[0].GetStatus(&helpersv1.Filters{FrameworkNames: []string{"NSA"}}).Status())
	assert.Equal(t, apis.StatusFailed, result.AssociatedControls[0].GetStatus(&helpersv1.Filters{FrameworkNames: []string{"MITRE"}}).Status())
}
