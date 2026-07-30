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

// objectString returns the string at path within obj, false if missing or not a string.
// GetNamespace and GetApiVersion assert v.(string) without comma-ok, so they panic here.
func objectString(obj map[string]interface{}, path ...string) (string, bool) {
	v, ok := workloadinterface.InspectMap(obj, path...)
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

// namespacesCompatible - unset matches anything, on the binding side too: a namespace
// less RoleBinding manifest still pairs with a Role anywhere, which beats losing every
// finding written that way. Cluster reads, where the bug bites, always set it.
func namespacesCompatible(roleNamespace, bindingNamespace string) bool {
	return roleNamespace == "" || bindingNamespace == "" || roleNamespace == bindingNamespace
}

// roleIsNamespaced - by definition for Role and ClusterRole, by its namespace for a CRD.
func roleIsNamespaced(roleKind, roleNamespace string) bool {
	switch roleKind {
	case "Role":
		return true
	case "ClusterRole":
		return false
	default:
		return roleNamespace != ""
	}
}

// bindingResolvesNamespacedRole - a Role is bindable only by a RoleBinding, a CRD role by
// any binding kind except a ClusterRoleBinding, which never reaches into a namespace.
func bindingResolvesNamespacedRole(roleKind, bindingKind string) bool {
	if roleKind == "Role" {
		return bindingKind == "RoleBinding"
	}
	return bindingKind != "ClusterRoleBinding"
}

// bindingReferencesRole reports whether bindingWorkload's roleRef resolves to roleWorkload.
// Kind and name alone pair it with a same named Role from an unrelated namespace.
func bindingReferencesRole(bindingWorkload, roleWorkload workloadinterface.IMetadata) bool {
	bindingObj, roleObj := bindingWorkload.GetObject(), roleWorkload.GetObject()

	name, ok := objectString(bindingObj, "roleRef", "name")
	if !ok || name != roleWorkload.GetName() {
		return false
	}

	kind, ok := objectString(bindingObj, "roleRef", "kind")
	roleKind, roleKindOK := objectString(roleObj, "kind")
	if !ok || !roleKindOK || kind != roleKind {
		return false
	}

	// apiGroup and the namespaces below are qualifiers: missing or corrupt reads as unset,
	// so garbage in one cannot hide a subject from these rules
	apiGroup, _ := objectString(bindingObj, "roleRef", "apiGroup")
	apiVersion, _ := objectString(roleObj, "apiVersion")
	roleGroup, _ := k8sinterface.SplitApiVersion(apiVersion)
	if apiGroup != "" && roleGroup != "" && apiGroup != roleGroup {
		return false
	}

	roleNamespace, _ := objectString(roleObj, "metadata", "namespace")
	if !roleIsNamespaced(roleKind, roleNamespace) {
		return true // cluster scoped, either binding kind may reference it
	}

	// namespaced, only a binding alongside it resolves
	bindingKind, _ := objectString(bindingObj, "kind")
	bindingNamespace, _ := objectString(bindingObj, "metadata", "namespace")
	return bindingResolvesNamespacedRole(roleKind, bindingKind) &&
		namespacesCompatible(roleNamespace, bindingNamespace)
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
