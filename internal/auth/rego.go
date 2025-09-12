package auth

import (
	"context"
	"fmt"

	"github.com/open-policy-agent/opa/v1/rego"
)

// Rego ...
type Rego struct {
	policy *rego.PreparedEvalQuery
}

// NewRego ...
func NewRego(policy *rego.PreparedEvalQuery) *Rego {
	return &Rego{policy: policy}
}

// Authorize ...
func (r *Rego) Authorize(ctx context.Context, input map[string]any) error {
	result, err := r.policy.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return fmt.Errorf("policy.Eval: %w", err)
	}

	if !result.Allowed() {
		return fmt.Errorf("policy.Eval: policy result not allowed")
	}

	return nil
}
