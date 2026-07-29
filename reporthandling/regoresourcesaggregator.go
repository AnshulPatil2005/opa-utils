package reporthandling

import (
	"bytes"
	"encoding/gob"
	"strings"

	"github.com/kubescape/k8s-interface/k8sinterface"
	"github.com/kubescape/k8s-interface/workloadinterface"
	"github.com/kubescape/opa-utils/objectsenvelopes"
)

var aggregatorAttribute = "resourcesAggregator"

func RegoResourcesAggregator(rule *PolicyRule, k8sObjects []workloadinterface.IMetadata) ([]workloadinterface.IMetadata, error) {
	if aggregateBy, ok := rule.Attributes[aggregatorAttribute]; ok {
		switch aggregateBy {
		case "subject-role-rolebinding":
			return AggregateResourcesBySubjects(k8sObjects)
		case "apiserver-pod":
			if obj := AggregateResourcesByAPIServerPod(k8sObjects); obj != nil {
				return []workloadinterface.IMetadata{obj}, nil
			}
			return []workloadinterface.IMetadata{}, nil
		default:
			return k8sObjects, nil
		}
	}
	return k8sObjects, nil
}

// roleRefString returns roleRef.<field> as a string, false if missing or not a string.
func roleRefString(bindingObj map[string]interface{}, field string) (string, bool) {
	v, ok := workloadinterface.InspectMap(bindingObj, "roleRef", field)
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

// namespacesCompatible - an unset namespace means "whatever this is applied into",
// which a file scan cannot know, so it matches anything. Cluster reads always set it.
func namespacesCompatible(roleNamespace, bindingNamespace string) bool {
	return roleNamespace == "" || bindingNamespace == "" || roleNamespace == bindingNamespace
}

// bindingReferencesRole reports whether the roleRef of bindingWorkload resolves to
// roleWorkload. Kind and name alone pair a binding with a same named Role from an
// unrelated namespace, which the subject has no access to at all.
func bindingReferencesRole(bindingWorkload, roleWorkload workloadinterface.IMetadata) bool {
	bindingObj := bindingWorkload.GetObject()

	name, ok := roleRefString(bindingObj, "name")
	if !ok || name != roleWorkload.GetName() {
		return false
	}

	kind, ok := roleRefString(bindingObj, "kind")
	if !ok || kind != roleWorkload.GetKind() {
		return false
	}

	// apiGroup is often omitted in hand written manifests
	apiGroup, _ := roleRefString(bindingObj, "apiGroup")
	roleGroup, _ := k8sinterface.SplitApiVersion(roleWorkload.GetApiVersion())
	if apiGroup != "" && roleGroup != "" && apiGroup != roleGroup {
		return false
	}

	switch kind {
	case "ClusterRole": // cluster scoped, either binding kind may reference it
		return true
	case "Role": // namespaced, only a RoleBinding alongside it resolves
		return bindingWorkload.GetKind() == "RoleBinding" &&
			namespacesCompatible(roleWorkload.GetNamespace(), bindingWorkload.GetNamespace())
	default: // a CRD kind: keep pairing, but not across two namespaces that differ
		return namespacesCompatible(roleWorkload.GetNamespace(), bindingWorkload.GetNamespace())
	}
}

func AggregateResourcesBySubjects(k8sObjects []workloadinterface.IMetadata) ([]workloadinterface.IMetadata, error) {
	aggregatedK8sObjects := []workloadinterface.IMetadata{}
	for _, bindingWorkload := range k8sObjects {
		if !strings.HasSuffix(bindingWorkload.GetKind(), "Binding") {
			continue
		}
		subjects, ok := workloadinterface.InspectMap(bindingWorkload.GetObject(), "subjects")
		if !ok {
			continue
		}
		data, ok := subjects.([]interface{})
		if !ok {
			continue
		}
		for _, roleWorkload := range k8sObjects {
			if !strings.HasSuffix(roleWorkload.GetKind(), "Role") {
				continue
			}
			if !bindingReferencesRole(bindingWorkload, roleWorkload) {
				continue
			}
			for _, subject := range data {
				subjectMap, ok := subject.(map[string]interface{})
				if !ok {
					continue
				}
				// deep copy subject - don't change original subject in rolebinding
				subjectCopy, err := DeepCopyMap(subjectMap)
				if err != nil {
					return aggregatedK8sObjects, err
				}
				subjectCopy[objectsenvelopes.RelatedObjectsKey] = []map[string]interface{}{bindingWorkload.GetObject(), roleWorkload.GetObject()}
				newObj := objectsenvelopes.NewRegoResponseVectorObject(subjectCopy)
				aggregatedK8sObjects = append(aggregatedK8sObjects, newObj)
			}
		}
	}
	return aggregatedK8sObjects, nil
}

// Create custom object of apiserver pod. Has required fields + cmdline
func AggregateResourcesByAPIServerPod(k8sObjects []workloadinterface.IMetadata) workloadinterface.IMetadata {
	for _, obj := range k8sObjects {
		if !k8sinterface.IsTypeWorkload(obj.GetObject()) {
			continue
		}
		workload := workloadinterface.NewWorkloadObj(obj.GetObject())
		if workload.GetKind() == "Pod" && workload.GetNamespace() == "kube-system" {
			if strings.Contains(workload.GetName(), "apiserver") || strings.Contains(workload.GetName(), "api-server") {
				return workload
			}
		}
	}
	return nil
}

// DeepCopyMap performs a deep copy of the given map m.
func DeepCopyMap(m map[string]interface{}) (map[string]interface{}, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)
	err := enc.Encode(m)
	if err != nil {
		return nil, err
	}
	var copy map[string]interface{}
	err = dec.Decode(&copy)
	if err != nil {
		return nil, err
	}
	return copy, nil
}
