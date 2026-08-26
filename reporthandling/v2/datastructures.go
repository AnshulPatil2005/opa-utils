package v2

import (
	"time"

	armoapi "github.com/armosec/armoapi-go/apis"
	"github.com/kubescape/opa-utils/reporthandling"
	"github.com/kubescape/opa-utils/reporthandling/apis"
	"github.com/kubescape/opa-utils/reporthandling/results/v1/reportsummary"
	"github.com/kubescape/opa-utils/reporthandling/results/v1/resourcesresults"

	"k8s.io/apimachinery/pkg/version"
)

// PostureReport posture scanning report structure
type PostureReport struct {
	ReportGenerationTime  time.Time                         `json:"generationTime"`
	ClusterAPIServerInfo  *version.Info                     `json:"clusterAPIServerInfo"`
	ClusterCloudProvider  string                            `json:"clusterCloudProvider"`
	CustomerGUID          string                            `json:"customerGUID"`
	ClusterName           string                            `json:"clusterName"`
	ReportID              string                            `json:"reportGUID"`
	JobID                 string                            `json:"jobID"`
	SummaryDetails        reportsummary.SummaryDetails      `json:"summaryDetails,omitempty"`
	Resources             []reporthandling.Resource         `json:"resources,omitempty"`
	Attributes            []reportsummary.PostureAttributes `json:"attributes"`
	Results               []resourcesresults.Result         `json:"results,omitempty"`
	Metadata              Metadata                          `json:"metadata,omitempty"`
	PaginationInfo        armoapi.PaginationMarks           `json:"paginationInfo"`
	CustomerGUIDGenerated bool                              `json:"customerGUIDGenerated"`
	TriggeredByCLI        bool                              `json:"triggeredByCLI,omitempty"`
}

type ClusterMetadata struct {
	MapNamespaceToNumberOfResources map[string]int `json:"namespaceToNumberOfResources,omitempty"`
	CloudMetadata                   *CloudMetadata `json:"cloudMetadata,omitempty"`
	CloudProvider                   string         `json:"cloudProvider,omitempty"` // Deprecated - info should be in cloudMetadata
	ContextName                     string         `json:"contextName,omitempty"`
	NumberOfWorkerNodes             int            `json:"numberOfWorkerNodes,omitempty"`
}

// CloudMetadata metadata of the cloud the cluster is running on. Compatible with the reporthandling.ICloudMetadata interface
type CloudMetadata struct {
	CloudProvider apis.CloudProviderName `json:"cloudProvider,omitempty"`
	ShortName     string                 `json:"shortName,omitempty"`
	FullName      string                 `json:"fullName,omitempty"`
	PrefixName    string                 `json:"prefixName,omitempty"`
}

type RepoContextMetadata struct {
	Provider      string                    `json:"provider,omitempty"` // repo provider name. e.g. github, gitlab
	Repo          string                    `json:"repo,omitempty"`
	Owner         string                    `json:"owner,omitempty"`
	Branch        string                    `json:"branch,omitempty"`
	DefaultBranch string                    `json:"defaultBranch,omitempty"`
	RemoteURL     string                    `json:"remoteURL,omitempty"`
	LastCommit    reporthandling.LastCommit `json:"lastCommit,omitempty"`
	LocalRootPath string                    `json:"localRootPath,omitempty"` // repo root path (local)
}

type FileContextMetadata struct {
	FilePath string `json:"filePath,omitempty"` // like "hostpath"
	HostName string `json:"hostName,omitempty"` // like "hostpath"
}
type DirectoryContextMetadata struct {
	BasePath string `json:"basePath,omitempty"` // the scanning request base path
	HostName string `json:"hostName,omitempty"` // like "hostpath"
}

type HelmContextMetadata struct {
	ChartName string `json:"chartName,omitempty"`
}
type ContextMetadata struct {
	ClusterContextMetadata   *ClusterMetadata          `json:"clusterContextMetadata,omitempty"`
	RepoContextMetadata      *RepoContextMetadata      `json:"gitRepoContextMetadata,omitempty"`
	FileContextMetadata      *FileContextMetadata      `json:"fileContextMetadata,omitempty"`
	HelmContextMetadata      *HelmContextMetadata      `json:"helmContextMetadata,omitempty"`
	DirectoryContextMetadata *DirectoryContextMetadata `json:"directoryContextMetadata,omitempty"`
}

type Metadata struct {
	ContextMetadata    ContextMetadata     `json:"targetMetadata,omitempty"`
	ClusterMetadata    ClusterMetadata     `json:"clusterMetadata,omitempty"`
	ScanMetadata       ScanMetadata        `json:"scanMetadata,omitempty"`
	EncryptionMetadata *EncryptionMetadata `json:"encryptionMetadata,omitempty"`
}

type EncryptionMetadata struct {
	Version      string `json:"version,omitempty"`
	DEKAlgorithm string `json:"dekAlgorithm,omitempty"`
	KEKAlgorithm string `json:"kekAlgorithm,omitempty"`
	EncryptedDEK string `json:"encryptedDEK,omitempty"`
}

type ScanningTarget uint16

const (
	Cluster   ScanningTarget = 0
	File      ScanningTarget = 1
	Repo      ScanningTarget = 2
	GitLocal  ScanningTarget = 3
	Directory ScanningTarget = 4
)

type ScanMetadata struct {
	TargetType       string `json:"targetType,omitempty"`
	KubescapeVersion string `json:"kubescapeVersion,omitempty"`
	FormatVersion    string `json:"formatVersion,omitempty"`
	ControlsInputs   string `json:"controlsInputs,omitempty"`
	// Format that has been requested for the output results.
	//
	// Since Kubescape added support for multiple outputs, might be not a
	// single format, but a comma-separated string of the multiple
	// requested formats.
	//
	// Deprecated: Since Kubescape added support for multiple outputs,
	// `Format` exists only for backward compatibility. Please use the
	// `Formats` field instead.
	Format string `json:"format,omitempty"`
	// Formats that have been requested for the output results.
	Formats             []string       `json:"formats,omitempty"`
	UseExceptions       string         `json:"useExceptions,omitempty"`
	Logger              string         `json:"logger,omitempty"`
	ExcludedNamespaces  []string       `json:"excludedNamespaces,omitempty"`
	IncludeNamespaces   []string       `json:"includeNamespaces,omitempty"`
	TargetNames         []string       `json:"targetNames,omitempty"`
	FailThreshold       float32        `json:"failThreshold,omitempty"`
	ComplianceThreshold float32        `json:"complianceThreshold,omitempty"`
	ScanningTarget      ScanningTarget `json:"scanningTarget,omitempty"`
	HostScanner         bool           `json:"hostScanner,omitempty"`
	Submit              bool           `json:"submit,omitempty"`
	VerboseMode         bool           `json:"verboseMode,omitempty"`
	// ScanContract captures the selected repository scan contract and the
	// resolved scan inputs it influenced. It is absent for scans that did not
	// use a repository scan contract.
	ScanContract *ScanContractMetadata `json:"scanContract,omitempty"`
}

// ScanContractMetadata is the provenance envelope for a selected repository
// scan contract. It deliberately records only identifiers, normalized scan
// inputs, and digests. It must never contain runner credentials, file contents,
// or absolute host paths.
type ScanContractMetadata struct {
	APIVersion              string                         `json:"apiVersion,omitempty"`
	Name                    string                         `json:"name,omitempty"`
	Contract                string                         `json:"contract,omitempty"`
	MinimumKubescapeVersion string                         `json:"minimumKubescapeVersion,omitempty"`
	DigestSchema            string                         `json:"digestSchema,omitempty"`
	ContractDigest          string                         `json:"contractDigest,omitempty"`
	EffectiveRunDigest      string                         `json:"effectiveRunDigest,omitempty"`
	Source                  string                         `json:"source,omitempty"`
	AllowedSections         []string                       `json:"allowedSections,omitempty"`
	DeniedSections          []string                       `json:"deniedSections,omitempty"`
	Effective               *ScanContractEffectiveSettings `json:"effective,omitempty"`
	RunnerInputs            []ScanContractRunnerInput      `json:"runnerInputs,omitempty"`
	GateResolution          *ScanContractGateResolution    `json:"gateResolution,omitempty"`
	OrdinaryCLIOverrides    *ScanContractCLIOverrides      `json:"ordinaryCliOverrides,omitempty"`
}

// ScanContractRunnerInput identifies a runner-owned file that affected a scan.
// Source is repository-relative when the input belongs to the repository, or
// the literal "external" when revealing a host path would be unsafe.
type ScanContractRunnerInput struct {
	Role   string `json:"role,omitempty"`
	Source string `json:"source,omitempty"`
	Digest string `json:"digest,omitempty"`
}

// ScanContractEffectiveSettings contains the post-resolution contract values
// that were provided to the scan. A nil section means that it did not
// participate in the resolved contract.
type ScanContractEffectiveSettings struct {
	Policy     *ScanContractPolicy     `json:"policy,omitempty"`
	Scope      *ScanContractScope      `json:"scope,omitempty"`
	Evaluation *ScanContractEvaluation `json:"evaluation,omitempty"`
	Failure    *ScanContractFailure    `json:"failure,omitempty"`
	Output     *ScanContractOutput     `json:"output,omitempty"`
}

type ScanContractPolicy struct {
	Frameworks      []string `json:"frameworks,omitempty"`
	Controls        []string `json:"controls,omitempty"`
	ControlsVersion string   `json:"controlsVersion,omitempty"`
}

type ScanContractScope struct {
	IncludeNamespaces []string `json:"includeNamespaces,omitempty"`
	ExcludeNamespaces []string `json:"excludeNamespaces,omitempty"`
}

type ScanContractEvaluation struct {
	ScanTimeout    string `json:"scanTimeout,omitempty"`
	ControlTimeout string `json:"controlTimeout,omitempty"`
}

// Pointer fields preserve the distinction between omitted gates and explicit
// zero or false values.
type ScanContractFailure struct {
	SeverityAtLeast     *string  `json:"severityAtLeast,omitempty"`
	ComplianceBelow     *float64 `json:"complianceBelow,omitempty"`
	CoverageBelow       *float64 `json:"coverageBelow,omitempty"`
	DegradedPolicyInput *bool    `json:"degradedPolicyInput,omitempty"`
}

type ScanContractOutput struct {
	Formats          []string `json:"formats,omitempty"`
	OmitRawResources *bool    `json:"omitRawResources,omitempty"`
}

// ScanContractGateResolution records the contract value, the trusted runner
// floor, and the resulting effective value for each monotonic gate.
type ScanContractGateResolution struct {
	SeverityAtLeast     ScanContractStringGateResolution `json:"severityAtLeast,omitempty"`
	ComplianceBelow     ScanContractNumberGateResolution `json:"complianceBelow,omitempty"`
	CoverageBelow       ScanContractNumberGateResolution `json:"coverageBelow,omitempty"`
	DegradedPolicyInput ScanContractBoolGateResolution   `json:"degradedPolicyInput,omitempty"`
}

type ScanContractStringGateResolution struct {
	Contract    *string `json:"contract,omitempty"`
	RunnerFloor *string `json:"runnerFloor,omitempty"`
	Effective   *string `json:"effective,omitempty"`
}

type ScanContractNumberGateResolution struct {
	Contract    *float64 `json:"contract,omitempty"`
	RunnerFloor *float64 `json:"runnerFloor,omitempty"`
	Effective   *float64 `json:"effective,omitempty"`
}

type ScanContractBoolGateResolution struct {
	Contract    *bool `json:"contract,omitempty"`
	RunnerFloor *bool `json:"runnerFloor,omitempty"`
	Effective   *bool `json:"effective,omitempty"`
}

// ScanContractCLIOverrides preserves explicit runner-owned ordinary settings.
// Pointer fields retain an explicit empty list or false value separately from
// an omitted flag.
type ScanContractCLIOverrides struct {
	Policy     *ScanContractPolicyOverrides     `json:"policy,omitempty"`
	Scope      *ScanContractScopeOverrides      `json:"scope,omitempty"`
	Evaluation *ScanContractEvaluationOverrides `json:"evaluation,omitempty"`
	Output     *ScanContractOutputOverrides     `json:"output,omitempty"`
}

type ScanContractPolicyOverrides struct {
	Frameworks      *[]string `json:"frameworks,omitempty"`
	Controls        *[]string `json:"controls,omitempty"`
	ControlsVersion *string   `json:"controlsVersion,omitempty"`
}

type ScanContractScopeOverrides struct {
	IncludeNamespaces *[]string `json:"includeNamespaces,omitempty"`
	ExcludeNamespaces *[]string `json:"excludeNamespaces,omitempty"`
}

type ScanContractEvaluationOverrides struct {
	ScanTimeout    *string `json:"scanTimeout,omitempty"`
	ControlTimeout *string `json:"controlTimeout,omitempty"`
}

type ScanContractOutputOverrides struct {
	Formats          *[]string `json:"formats,omitempty"`
	OmitRawResources *bool     `json:"omitRawResources,omitempty"`
}

// Moved to apis/cloudmetadata.go
// const (
// 	GKE = "GKE"
// 	GCP = "GCP"
// 	EKS = "EKS"
// )

func (st *ScanningTarget) String() string {
	switch *st {
	case 0:
		return "Cluster"
	case 1:
		return "File"
	case 2:
		return "Repo"
	case 3:
		return "GitLocal"
	case 4:
		return "Directory"
	default:
		return ""
	}
}
