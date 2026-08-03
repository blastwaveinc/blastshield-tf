// Package planmods provides shared plan modifiers used by the generated
// per-version schema code.
package planmods

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// PreserveStateWhenUnconfigured returns a plan modifier that keeps the prior
// state value — including null — whenever the attribute is absent from the
// configuration.
//
// This exists because Optional+Computed attributes with a null config AND a
// null prior state are planned as "(known after apply)" by the framework,
// producing a perpetual update loop for fields the API legitimately returns
// as null (e.g. endpoint.api_access). UseStateForUnknown cannot fix that
// case: it intentionally skips null prior state.
//
// Semantics: if the config value is null, the planned value is the prior
// state value, whatever it is. Removing a previously-configured value from
// config therefore stops managing it rather than clearing it server-side —
// the same trade-off UseStateForUnknown makes.
func PreserveStateWhenUnconfigured() preserveStateWhenUnconfigured {
	return preserveStateWhenUnconfigured{}
}

type preserveStateWhenUnconfigured struct{}

func (m preserveStateWhenUnconfigured) Description(_ context.Context) string {
	return "Keeps the prior state value (including null) when the attribute is not set in configuration."
}

func (m preserveStateWhenUnconfigured) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m preserveStateWhenUnconfigured) PlanModifyBool(_ context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if req.ConfigValue.IsNull() && !req.State.Raw.IsNull() {
		resp.PlanValue = req.StateValue
	}
}

func (m preserveStateWhenUnconfigured) PlanModifyInt64(_ context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	if req.ConfigValue.IsNull() && !req.State.Raw.IsNull() {
		resp.PlanValue = req.StateValue
	}
}

func (m preserveStateWhenUnconfigured) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.ConfigValue.IsNull() && !req.State.Raw.IsNull() {
		resp.PlanValue = req.StateValue
	}
}

func (m preserveStateWhenUnconfigured) PlanModifyList(_ context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if req.ConfigValue.IsNull() && !req.State.Raw.IsNull() {
		resp.PlanValue = req.StateValue
	}
}

func (m preserveStateWhenUnconfigured) PlanModifySet(_ context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	if req.ConfigValue.IsNull() && !req.State.Raw.IsNull() {
		resp.PlanValue = req.StateValue
	}
}

func (m preserveStateWhenUnconfigured) PlanModifyMap(_ context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {
	if req.ConfigValue.IsNull() && !req.State.Raw.IsNull() {
		resp.PlanValue = req.StateValue
	}
}
