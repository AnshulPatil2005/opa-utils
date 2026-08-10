package resources

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/kubescape/opa-utils/reporthandling"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetFilteredPostureControlInputs(t *testing.T) {
	regoInputData := RegoDependenciesData{}
	regoInputData.PostureControlInputs = map[string][]string{"sensitiveKeyNames": {"keyA", "keyB"}}
	s := []string{"settings.postureControlInputs.sensitiveKeyNames", "settings.postureControlInputs.blabla"}
	postureControlInputs := regoInputData.GetFilteredPostureControlInputs(s)
	splitted0 := strings.Split(s[0], ".")
	_, ok := postureControlInputs[splitted0[2]]
	assert.True(t, ok)

	splitted1 := strings.Split(s[1], ".")
	_, ok = postureControlInputs[splitted1[2]]
	assert.False(t, ok)
}

func TestGetFilteredPostureControlConfigInputs(t *testing.T) {
	regoInputData := RegoDependenciesData{}
	regoInputData.PostureControlInputs = map[string][]string{"sensitiveKeyNames": {"keyA", "keyB"}}

	inputs := []reporthandling.ControlConfigInputs{
		{
			Path: "settings.postureControlInputs.sensitiveKeyNames",
			Name: "Sensitive Key Names",
		},
		{
			Path: "settings.postureControlInputs.blabla",
			Name: "Blabla",
		},
	}

	postureControlInputs := regoInputData.GetFilteredPostureControlConfigInputs(inputs)

	splitted0 := strings.Split(inputs[0].Path, ".")
	_, ok := postureControlInputs[splitted0[2]]
	assert.True(t, ok)

	splitted1 := strings.Split(inputs[1].Path, ".")
	_, ok = postureControlInputs[splitted1[2]]
	assert.False(t, ok)
}

func TestLoadRegoFiles(t *testing.T) {
	dir := t.TempDir()

	// File names chosen so that trimming ".rego" as a character cutset rather
	// than as a suffix produces a visibly wrong key: every one of these shares
	// at least one leading or trailing character with the set {., r, e, g, o}.
	files := map[string]string{
		"raw.rego":     "package raw",
		"rule.rego":    "package rule",
		"role.rego":    "package role",
		"error.rego":   "package error",
		"storage.rego": "package storage",
		"generic.rego": "package generic",
		"cautils.rego": "package cautils",
		"notes.txt":    "not a rego file",
	}
	for name, content := range files {
		require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600))
	}

	modules := LoadRegoFiles(dir)

	for _, name := range []string{"raw", "rule", "role", "error", "storage", "generic", "cautils"} {
		content, ok := modules[name]
		assert.Truef(t, ok, "module %q missing; got keys %v", name, keysOf(modules))
		assert.Equal(t, "package "+name, content)
	}

	assert.NotContains(t, modules, "", "a file must never produce an empty module key")
	assert.Len(t, modules, len(files)-1, "non-.rego files must be skipped")
}

func TestLoadRegoFilesNonExistentDir(t *testing.T) {
	assert.Empty(t, LoadRegoFiles(filepath.Join(t.TempDir(), "does-not-exist")))
}

func keysOf(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
