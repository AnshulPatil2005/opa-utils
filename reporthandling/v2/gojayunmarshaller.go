package v2

import (
	"time"

	"github.com/francoispqt/gojay"
)

/*
  responsible on fast unmarshaling of various COMMON posture report v2 structure for basic validation

*/
// UnmarshalJSONObject - File inside a pkg
func (r *PostureReport) UnmarshalJSONObject(dec *gojay.Decoder, key string) (err error) {

	switch key {
	case "customerGUID":
		err = dec.String(&(r.CustomerGUID))

	case "clusterName":
		err = dec.String(&(r.ClusterName))

	case "reportGUID":
		err = dec.String(&(r.ReportID))
	case "jobID":
		err = dec.String(&(r.JobID))
	case "generationTime":
		err = dec.Time(&(r.ReportGenerationTime), time.RFC3339)
		r.ReportGenerationTime = r.ReportGenerationTime.Local()
	case "metadata":
		err = dec.Object(&(r.Metadata))
	case "customerGUIDGenerated":
		err = dec.Bool(&(r.CustomerGUIDGenerated))
	case "triggeredByCLI":
		err = dec.Bool(&(r.TriggeredByCLI))
	case "paginationInfo":
		err = dec.Object(&(r.PaginationInfo))
	}
	return err
}

// func (files *PkgFiles) UnmarshalJSONArray(dec *gojay.Decoder) error {
// 	lae := PackageFile{}
// 	if err := dec.Object(&lae); err != nil {
// 		return err
// 	}

// 	*files = append(*files, lae)
// 	return nil
// }

func (file *PostureReport) NKeys() int {
	return 0
}

// UnmarshalJSONObject unmarshals incoming JSON data into a Metadata object
func (m *Metadata) UnmarshalJSONObject(dec *gojay.Decoder, key string) (err error) {

	switch key {
	case "scanMetadata":
		err = dec.Object(&(m.ScanMetadata))

	case "clusterMetadata":
		err = dec.Object(&(m.ClusterMetadata))

	case "targetMetadata":
		err = dec.Object(&(m.ContextMetadata))

	case "encryptionMetadata":
		encryptionMetadata := &EncryptionMetadata{}

		if err = dec.Object(encryptionMetadata); err == nil {
			m.EncryptionMetadata = encryptionMetadata
		}

	}

	return err
}

func (file *Metadata) NKeys() int {
	return 0
}

func (c *ContextMetadata) UnmarshalJSONObject(dec *gojay.Decoder, key string) (err error) {
	switch key {
	case "clusterContextMetadata":
		clusterMetadata := &ClusterMetadata{}
		if err = dec.Object(clusterMetadata); err == nil {
			c.ClusterContextMetadata = clusterMetadata
		}
	case "gitRepoContextMetadata":
		repoContextMetadata := &RepoContextMetadata{}
		if err = dec.Object(repoContextMetadata); err == nil {
			c.RepoContextMetadata = repoContextMetadata
		}
	}
	return err
}

func (c *ContextMetadata) NKeys() int {
	return 0
}

func (c *RepoContextMetadata) UnmarshalJSONObject(dec *gojay.Decoder, key string) (err error) {
	switch key {
	case "provider":
		err = dec.String(&(c.Provider))
	case "repo":
		err = dec.String(&(c.Repo))
	case "owner":
		err = dec.String(&(c.Owner))
	case "branch":
		err = dec.String(&(c.Branch))
	case "remoteURL":
		err = dec.String(&(c.RemoteURL))
	}
	return err
}

func (c *RepoContextMetadata) NKeys() int {
	return 0
}

// UnmarshalJSONObject unmarshals incoming JSON data into a ScanMetadata object
func (m *ScanMetadata) UnmarshalJSONObject(dec *gojay.Decoder, key string) (err error) {

	switch key {
	case "format": // string
		err = dec.String(&(m.Format))
	case "formats":
		err = dec.SliceString(&(m.Formats))
	case "excludedNamespaces": // []string
		err = dec.SliceString(&(m.ExcludedNamespaces))
	case "includeNamespaces": // []string
		err = dec.SliceString(&(m.IncludeNamespaces))
	case "failThreshold": // float32
		err = dec.Float32(&(m.FailThreshold))
	case "submit": // bool
		err = dec.Bool(&(m.Submit))
	case "hostScanner": // bool
		err = dec.Bool(&(m.HostScanner))
	case "logger": // string
		err = dec.String(&(m.Logger))
	case "targetType": // string
		err = dec.String(&(m.TargetType))
	case "targetNames": // []string
		err = dec.SliceString(&(m.TargetNames))
	case "useExceptions": // string
		err = dec.String(&(m.UseExceptions))
	case "controlsInputs": // string
		err = dec.String(&(m.ControlsInputs))
	case "verboseMode": // bool
		err = dec.Bool(&(m.VerboseMode))
	case "scanContract":
		contract := &ScanContractMetadata{}
		if err = dec.Object(contract); err == nil {
			m.ScanContract = contract
		}
	case "scanningTarget":
		var value uint16
		err = dec.Uint16(&value)
		m.ScanningTarget = ScanningTarget(value)
	}
	return err

}

func (file *ScanMetadata) NKeys() int {
	return 0
}

func (m *ScanContractMetadata) UnmarshalJSONObject(dec *gojay.Decoder, key string) (err error) {
	switch key {
	case "apiVersion":
		err = dec.String(&m.APIVersion)
	case "name":
		err = dec.String(&m.Name)
	case "contract":
		err = dec.String(&m.Contract)
	case "minimumKubescapeVersion":
		err = dec.String(&m.MinimumKubescapeVersion)
	case "digestSchema":
		err = dec.String(&m.DigestSchema)
	case "contractDigest":
		err = dec.String(&m.ContractDigest)
	case "effectiveRunDigest":
		err = dec.String(&m.EffectiveRunDigest)
	case "source":
		err = dec.String(&m.Source)
	case "allowedSections":
		err = dec.SliceString(&m.AllowedSections)
	case "deniedSections":
		err = dec.SliceString(&m.DeniedSections)
	case "effective":
		value := &ScanContractEffectiveSettings{}
		if err = dec.Object(value); err == nil {
			m.Effective = value
		}
	case "runnerInputs":
		var inputs scanContractRunnerInputs
		if err = dec.Array(&inputs); err == nil {
			m.RunnerInputs = inputs
		}
	case "gateResolution":
		value := &ScanContractGateResolution{}
		if err = dec.Object(value); err == nil {
			m.GateResolution = value
		}
	case "ordinaryCliOverrides":
		value := &ScanContractCLIOverrides{}
		if err = dec.Object(value); err == nil {
			m.OrdinaryCLIOverrides = value
		}
	}
	return err
}

func (*ScanContractMetadata) NKeys() int { return 0 }

type scanContractRunnerInputs []ScanContractRunnerInput

func (inputs *scanContractRunnerInputs) UnmarshalJSONArray(dec *gojay.Decoder) error {
	input := ScanContractRunnerInput{}
	if err := dec.Object(&input); err != nil {
		return err
	}
	*inputs = append(*inputs, input)
	return nil
}

func (input *ScanContractRunnerInput) UnmarshalJSONObject(dec *gojay.Decoder, key string) (err error) {
	switch key {
	case "role":
		err = dec.String(&input.Role)
	case "source":
		err = dec.String(&input.Source)
	case "digest":
		err = dec.String(&input.Digest)
	}
	return err
}

func (*ScanContractRunnerInput) NKeys() int { return 0 }

func (settings *ScanContractEffectiveSettings) UnmarshalJSONObject(dec *gojay.Decoder, key string) (err error) {
	switch key {
	case "policy":
		value := &ScanContractPolicy{}
		if err = dec.Object(value); err == nil {
			settings.Policy = value
		}
	case "scope":
		value := &ScanContractScope{}
		if err = dec.Object(value); err == nil {
			settings.Scope = value
		}
	case "evaluation":
		value := &ScanContractEvaluation{}
		if err = dec.Object(value); err == nil {
			settings.Evaluation = value
		}
	case "failure":
		value := &ScanContractFailure{}
		if err = dec.Object(value); err == nil {
			settings.Failure = value
		}
	case "output":
		value := &ScanContractOutput{}
		if err = dec.Object(value); err == nil {
			settings.Output = value
		}
	}
	return err
}

func (*ScanContractEffectiveSettings) NKeys() int { return 0 }

func (policy *ScanContractPolicy) UnmarshalJSONObject(dec *gojay.Decoder, key string) (err error) {
	switch key {
	case "frameworks":
		err = dec.SliceString(&policy.Frameworks)
	case "controls":
		err = dec.SliceString(&policy.Controls)
	case "controlsVersion":
		err = dec.String(&policy.ControlsVersion)
	}
	return err
}

func (*ScanContractPolicy) NKeys() int { return 0 }

func (scope *ScanContractScope) UnmarshalJSONObject(dec *gojay.Decoder, key string) (err error) {
	switch key {
	case "includeNamespaces":
		err = dec.SliceString(&scope.IncludeNamespaces)
	case "excludeNamespaces":
		err = dec.SliceString(&scope.ExcludeNamespaces)
	}
	return err
}

func (*ScanContractScope) NKeys() int { return 0 }

func (evaluation *ScanContractEvaluation) UnmarshalJSONObject(dec *gojay.Decoder, key string) (err error) {
	switch key {
	case "scanTimeout":
		err = dec.String(&evaluation.ScanTimeout)
	case "controlTimeout":
		err = dec.String(&evaluation.ControlTimeout)
	}
	return err
}

func (*ScanContractEvaluation) NKeys() int { return 0 }

func (failure *ScanContractFailure) UnmarshalJSONObject(dec *gojay.Decoder, key string) (err error) {
	switch key {
	case "severityAtLeast":
		failure.SeverityAtLeast, err = decodeContractString(dec)
	case "complianceBelow":
		failure.ComplianceBelow, err = decodeContractFloat64(dec)
	case "coverageBelow":
		failure.CoverageBelow, err = decodeContractFloat64(dec)
	case "degradedPolicyInput":
		failure.DegradedPolicyInput, err = decodeContractBool(dec)
	}
	return err
}

func (*ScanContractFailure) NKeys() int { return 0 }

func (output *ScanContractOutput) UnmarshalJSONObject(dec *gojay.Decoder, key string) (err error) {
	switch key {
	case "formats":
		err = dec.SliceString(&output.Formats)
	case "omitRawResources":
		output.OmitRawResources, err = decodeContractBool(dec)
	}
	return err
}

func (*ScanContractOutput) NKeys() int { return 0 }

func (resolution *ScanContractGateResolution) UnmarshalJSONObject(dec *gojay.Decoder, key string) (err error) {
	switch key {
	case "severityAtLeast":
		err = dec.Object(&resolution.SeverityAtLeast)
	case "complianceBelow":
		err = dec.Object(&resolution.ComplianceBelow)
	case "coverageBelow":
		err = dec.Object(&resolution.CoverageBelow)
	case "degradedPolicyInput":
		err = dec.Object(&resolution.DegradedPolicyInput)
	}
	return err
}

func (*ScanContractGateResolution) NKeys() int { return 0 }

func (resolution *ScanContractStringGateResolution) UnmarshalJSONObject(dec *gojay.Decoder, key string) (err error) {
	switch key {
	case "contract":
		resolution.Contract, err = decodeContractString(dec)
	case "runnerFloor":
		resolution.RunnerFloor, err = decodeContractString(dec)
	case "effective":
		resolution.Effective, err = decodeContractString(dec)
	}
	return err
}

func (*ScanContractStringGateResolution) NKeys() int { return 0 }

func (resolution *ScanContractNumberGateResolution) UnmarshalJSONObject(dec *gojay.Decoder, key string) (err error) {
	switch key {
	case "contract":
		resolution.Contract, err = decodeContractFloat64(dec)
	case "runnerFloor":
		resolution.RunnerFloor, err = decodeContractFloat64(dec)
	case "effective":
		resolution.Effective, err = decodeContractFloat64(dec)
	}
	return err
}

func (*ScanContractNumberGateResolution) NKeys() int { return 0 }

func (resolution *ScanContractBoolGateResolution) UnmarshalJSONObject(dec *gojay.Decoder, key string) (err error) {
	switch key {
	case "contract":
		resolution.Contract, err = decodeContractBool(dec)
	case "runnerFloor":
		resolution.RunnerFloor, err = decodeContractBool(dec)
	case "effective":
		resolution.Effective, err = decodeContractBool(dec)
	}
	return err
}

func (*ScanContractBoolGateResolution) NKeys() int { return 0 }

func (overrides *ScanContractCLIOverrides) UnmarshalJSONObject(dec *gojay.Decoder, key string) (err error) {
	switch key {
	case "policy":
		value := &ScanContractPolicyOverrides{}
		if err = dec.Object(value); err == nil {
			overrides.Policy = value
		}
	case "scope":
		value := &ScanContractScopeOverrides{}
		if err = dec.Object(value); err == nil {
			overrides.Scope = value
		}
	case "evaluation":
		value := &ScanContractEvaluationOverrides{}
		if err = dec.Object(value); err == nil {
			overrides.Evaluation = value
		}
	case "output":
		value := &ScanContractOutputOverrides{}
		if err = dec.Object(value); err == nil {
			overrides.Output = value
		}
	}
	return err
}

func (*ScanContractCLIOverrides) NKeys() int { return 0 }

func (overrides *ScanContractPolicyOverrides) UnmarshalJSONObject(dec *gojay.Decoder, key string) (err error) {
	switch key {
	case "frameworks":
		overrides.Frameworks, err = decodeContractStringSlice(dec)
	case "controls":
		overrides.Controls, err = decodeContractStringSlice(dec)
	case "controlsVersion":
		overrides.ControlsVersion, err = decodeContractString(dec)
	}
	return err
}

func (*ScanContractPolicyOverrides) NKeys() int { return 0 }

func (overrides *ScanContractScopeOverrides) UnmarshalJSONObject(dec *gojay.Decoder, key string) (err error) {
	switch key {
	case "includeNamespaces":
		overrides.IncludeNamespaces, err = decodeContractStringSlice(dec)
	case "excludeNamespaces":
		overrides.ExcludeNamespaces, err = decodeContractStringSlice(dec)
	}
	return err
}

func (*ScanContractScopeOverrides) NKeys() int { return 0 }

func (overrides *ScanContractEvaluationOverrides) UnmarshalJSONObject(dec *gojay.Decoder, key string) (err error) {
	switch key {
	case "scanTimeout":
		overrides.ScanTimeout, err = decodeContractString(dec)
	case "controlTimeout":
		overrides.ControlTimeout, err = decodeContractString(dec)
	}
	return err
}

func (*ScanContractEvaluationOverrides) NKeys() int { return 0 }

func (overrides *ScanContractOutputOverrides) UnmarshalJSONObject(dec *gojay.Decoder, key string) (err error) {
	switch key {
	case "formats":
		overrides.Formats, err = decodeContractStringSlice(dec)
	case "omitRawResources":
		overrides.OmitRawResources, err = decodeContractBool(dec)
	}
	return err
}

func (*ScanContractOutputOverrides) NKeys() int { return 0 }

func decodeContractString(dec *gojay.Decoder) (*string, error) {
	value := ""
	if err := dec.String(&value); err != nil {
		return nil, err
	}
	return &value, nil
}

func decodeContractFloat64(dec *gojay.Decoder) (*float64, error) {
	var value float64
	if err := dec.Float64(&value); err != nil {
		return nil, err
	}
	return &value, nil
}

func decodeContractBool(dec *gojay.Decoder) (*bool, error) {
	var value bool
	if err := dec.Bool(&value); err != nil {
		return nil, err
	}
	return &value, nil
}

func decodeContractStringSlice(dec *gojay.Decoder) (*[]string, error) {
	var values []string
	if err := dec.SliceString(&values); err != nil {
		return nil, err
	}
	return &values, nil
}

// UnmarshalJSONObject unmarshals incoming JSON data into a ClusterMetadata object
func (m *ClusterMetadata) UnmarshalJSONObject(dec *gojay.Decoder, key string) (err error) {

	switch key {

	case "numberOfWorkerNodes": //int
		err = dec.Int(&(m.NumberOfWorkerNodes))

	case "cloudProvider": //string
		err = dec.String(&(m.CloudProvider))

	case "contextName": //string
		err = dec.String(&(m.ContextName))

	}
	return err

}

func (file *ClusterMetadata) NKeys() int {
	return 0
}

// UnmarshalJSONObject unmarshals incoming JSON data into an EncryptionMetadata object
func (e *EncryptionMetadata) UnmarshalJSONObject(
	dec *gojay.Decoder,
	key string,
) (err error) {

	switch key {
	case "version":
		err = dec.String(&(e.Version))

	case "dekAlgorithm":
		err = dec.String(&(e.DEKAlgorithm))

	case "kekAlgorithm":
		err = dec.String(&(e.KEKAlgorithm))

	case "encryptedDEK":
		err = dec.String(&(e.EncryptedDEK))
	}

	return err
}

func (e *EncryptionMetadata) NKeys() int {
	return 0
}
