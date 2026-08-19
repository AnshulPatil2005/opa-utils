package resourcesresults

import (
	"github.com/kubescape/opa-utils/reporthandling/apis"
	helpersv1 "github.com/kubescape/opa-utils/reporthandling/helpers/v1"
)

// GetName get rule name
func (rule *ResourceAssociatedRule) GetName() string {
	return rule.Name
}

// SetName set rule name
func (rule *ResourceAssociatedRule) SetName(n string) {
	rule.Name = n
}

// =============================== Status ====================================

// GetStatus get rule status
func (rule *ResourceAssociatedRule) GetStatus(f *helpersv1.Filters) apis.IStatus {
	if rule.Status != apis.StatusFailed {
		return rule.statusInfo()
	}

	exceptions := rule.Exception
	if f != nil {
		exceptions = f.FilterExceptions(exceptions)
	}
	if len(exceptions) > 0 {
		return &apis.StatusInfo{
			InnerStatus: apis.StatusPassed,
			SubStatus:   apis.SubStatusException,
		}
	}

	return rule.statusInfo()
}

func (rule *ResourceAssociatedRule) statusInfo() apis.IStatus {
	return &apis.StatusInfo{
		InnerStatus: rule.Status,
		SubStatus:   rule.SubStatus,
		InnerInfo:   apis.SubStatusInfo(rule.SubStatus),
	}
}

// GetSubStatus get rule sub status
func (rule *ResourceAssociatedRule) GetSubStatus() apis.ScanningSubStatus {
	return rule.SubStatus
}

// SetStatus set rule status and sub status
func (rule *ResourceAssociatedRule) SetStatus(s apis.ScanningStatus, _ *helpersv1.Filters) {
	rule.Status = s
}
