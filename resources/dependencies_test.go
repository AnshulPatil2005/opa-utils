package resources

import (
	"context"
	"testing"

	"github.com/open-policy-agent/opa/v1/ast"
	"github.com/open-policy-agent/opa/v1/rego"
)

func TestRegoModulesCompile(t *testing.T) {
	modules := map[string]string{
		"cautils.rego":               RegoCAUtils,
		"designators.rego":           RegoDesignators,
		"kubernetes_api_client.rego": RegoKubernetesApiClient,
	}

	for _, regoVersion := range []ast.RegoVersion{ast.RegoV1, ast.RegoV0} {
		t.Run(regoVersion.String(), func(t *testing.T) {
			compiler := ast.NewCompiler().WithDefaultRegoVersion(regoVersion)

			parsed := make(map[string]*ast.Module, len(modules))
			for filename, src := range modules {
				module, err := ast.ParseModuleWithOpts(filename, src, ast.ParserOptions{RegoVersion: regoVersion})
				if err != nil {
					t.Fatalf("failed to parse %s under %s: %v", filename, regoVersion, err)
				}
				parsed[filename] = module
			}

			compiler.Compile(parsed)
			if compiler.Failed() {
				t.Fatalf("compilation failed under %s: %v", regoVersion, compiler.Errors)
			}
		})
	}
}

func evalSingleResult(t *testing.T, ctx context.Context, query string) interface{} {
	t.Helper()

	r := rego.New(
		rego.Query(query),
		rego.Module("cautils.rego", RegoCAUtils),
		rego.Module("designators.rego", RegoDesignators),
	)

	rs, err := r.Eval(ctx)
	if err != nil {
		t.Fatalf("eval error for query %q: %v", query, err)
	}
	if len(rs) != 1 || len(rs[0].Expressions) != 1 {
		t.Fatalf("unexpected result set for query %q: %#v", query, rs)
	}
	return rs[0].Expressions[0].Value
}

func TestRegoCAUtilsEval(t *testing.T) {
	ctx := context.Background()

	if v := evalSingleResult(t, ctx, `data.cautils.list_contains(["a","b"], "b")`); v != true {
		t.Errorf("list_contains: expected true, got %v", v)
	}

	if v := evalSingleResult(t, ctx, `data.cautils.getPodName({"name": "foo"})`); v != "foo" {
		t.Errorf("getPodName: expected \"foo\", got %v", v)
	}

	permResult := evalSingleResult(t, ctx, `data.cautils.unix_permission(420)`)
	perm, ok := permResult.(map[string]interface{})
	if !ok {
		t.Fatalf("unix_permission: expected map result, got %#v", permResult)
	}
	expectBits := map[string]map[string]bool{
		"user":     {"read": true, "write": true, "exec": false},
		"group":    {"read": true, "write": false, "exec": false},
		"everyone": {"read": true, "write": false, "exec": false},
	}
	for section, bits := range expectBits {
		got, ok := perm[section].(map[string]interface{})
		if !ok {
			t.Fatalf("unix_permission: missing section %q in %#v", section, perm)
		}
		for bit, want := range bits {
			if got[bit] != want {
				t.Errorf("unix_permission.%s.%s: expected %v, got %v", section, bit, want, got[bit])
			}
		}
	}

	if v := evalSingleResult(t, ctx, `data.cautils.unix_permissions_allow(420, 384)`); v != true {
		t.Errorf("unix_permissions_allow(0644, 0600): expected true, got %v", v)
	}
}

func TestRegoDesignatorsEval(t *testing.T) {
	ctx := context.Background()

	if v := evalSingleResult(t, ctx, `data.designators.included_namespaces("default")`); v != true {
		t.Errorf("included_namespaces: expected true, got %v", v)
	}

	// excluded_namespaces("excluded") is undefined (not false): the rule body
	// `not cautils.list_contains(["excluded"], "excluded")` fails since "excluded"
	// is in the list, so the partial-function rule has no matching definition
	// and Rego yields an empty result set rather than a false value.
	r := rego.New(
		rego.Query(`data.designators.excluded_namespaces("excluded")`),
		rego.Module("cautils.rego", RegoCAUtils),
		rego.Module("designators.rego", RegoDesignators),
	)
	rs, err := r.Eval(ctx)
	if err != nil {
		t.Fatalf("eval error for excluded_namespaces: %v", err)
	}
	if len(rs) != 0 {
		t.Errorf("excluded_namespaces: expected undefined (empty result set), got %#v", rs)
	}
}
