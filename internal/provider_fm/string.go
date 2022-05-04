package provider_fm

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"math/big"
	"sort"
)

func getStringSchema(sensitive bool, description string) tfsdk.Schema {
	idDesc := "The generated random string."
	if sensitive {
		idDesc = "A static value used internally by Terraform, this should not be referenced in configurations."
	}

	return tfsdk.Schema{
		Description: description,
		Attributes: map[string]tfsdk.Attribute{
			"keepers": {
				Description: "Arbitrary map of values that, when changed, will trigger recreation of " +
					"resource. See [the main provider documentation](../index.html) for more information.",
				Type: types.MapType{
					ElemType: types.StringType,
				},
				Optional:      true,
				PlanModifiers: []tfsdk.AttributePlanModifier{tfsdk.RequiresReplace()},
			},

			"length": {
				Description: "The length of the string desired. The minimum value for length is 1 and, length " +
					"must also be >= (`min_upper` + `min_lower` + `min_numeric` + `min_special`).",
				Type:          types.Int64Type,
				Required:      true,
				PlanModifiers: []tfsdk.AttributePlanModifier{tfsdk.RequiresReplace()},
				Validators:    []tfsdk.AttributeValidator{lengthValidator{}},
			},

			"special": {
				Description: "Include special characters in the result. These are `!@#$%&*()-_=+[]{}<>:?`. Default value is `true`.",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					tfsdk.RequiresReplace(),
					defaultBool{},
				},
			},

			"upper": {
				Description: "Include uppercase alphabet characters in the result. Default value is `true`.",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					tfsdk.RequiresReplace(),
					defaultBool{},
				},
			},

			"lower": {
				Description: "Include lowercase alphabet characters in the result. Default value is `true`.",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					tfsdk.RequiresReplace(),
					defaultBool{},
				},
			},

			"number": {
				Description: "Include numeric characters in the result. Default value is `true`.",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					tfsdk.RequiresReplace(),
					defaultBool{},
				},
			},

			"min_numeric": {
				Description: "Minimum number of numeric characters in the result. Default value is `0`.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					tfsdk.RequiresReplace(),
					defaultInt{},
				},
			},

			"min_upper": {
				Description: "Minimum number of uppercase alphabet characters in the result. Default value is `0`.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					tfsdk.RequiresReplace(),
					defaultInt{},
				},
			},

			"min_lower": {
				Description: "Minimum number of lowercase alphabet characters in the result. Default value is `0`.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					tfsdk.RequiresReplace(),
					defaultInt{},
				},
			},

			"min_special": {
				Description: "Minimum number of special characters in the result. Default value is `0`.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					tfsdk.RequiresReplace(),
					defaultInt{},
				},
			},

			"override_special": {
				Description: "Supply your own list of special characters to use for string generation.  This " +
					"overrides the default character list in the special argument.  The `special` argument must " +
					"still be set to true for any overwritten characters to be used in generation.",
				Type:     types.StringType,
				Optional: true,
				Computed: true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					tfsdk.RequiresReplace(),
					defaultOverrideSpecial{},
				},
			},

			"result": {
				Description: "The generated random string.",
				Type:        types.StringType,
				Computed:    true,
				Sensitive:   sensitive,
			},

			"id": {
				Description: idDesc,
				Computed:    true,
				Type:        types.StringType,
			},
		},
	}
}

type lengthValidator struct{}

func (l lengthValidator) Description(context.Context) string {
	return "Length validator ensures that length is at least 1"
}

func (l lengthValidator) MarkdownDescription(context.Context) string {
	return "Length validator ensures that `length` is at least 1"
}

func (l lengthValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	t := req.AttributeConfig.(types.Int64)

	if t.Value < 1 {
		resp.Diagnostics.AddError(
			fmt.Sprintf("expected length to be at least 1, got %d", t.Value),
			fmt.Sprintf("expected length to be at least 1, got %d", t.Value),
		)
	}
}

type defaultBool struct{}

func (d defaultBool) Description(ctx context.Context) string {
	return "If the plan does not contain a value, a default will be set."
}

func (d defaultBool) MarkdownDescription(ctx context.Context) string {
	return "If the plan does not contain a value, a default will be set."
}

func (d defaultBool) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
	t := req.AttributeConfig.(types.Bool)

	if t.Null {
		resp.AttributePlan = types.Bool{
			Value: true,
		}
	}
}

type defaultInt struct{}

func (d defaultInt) Description(ctx context.Context) string {
	return "If the plan does not contain a value, a default will be set."
}

func (d defaultInt) MarkdownDescription(ctx context.Context) string {
	return "If the plan does not contain a value, a default will be set."
}

func (d defaultInt) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
	t := req.AttributeConfig.(types.Int64)

	if t.Null {
		resp.AttributePlan = types.Int64{
			Null:  false,
			Value: 0,
		}
	}
}

type defaultOverrideSpecial struct{}

func (d defaultOverrideSpecial) Description(ctx context.Context) string {
	return "If the plan does not contain a value, a default will be set."
}

func (d defaultOverrideSpecial) MarkdownDescription(ctx context.Context) string {
	return "If the plan does not contain a value, a default will be set."
}

func (d defaultOverrideSpecial) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
	t := req.AttributeConfig.(types.String)

	if t.Null {
		resp.AttributePlan = types.String{
			Null:  false,
			Value: "",
		}
	}
}

func createString(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse, sensitive bool) {
	const numChars = "0123456789"
	const lowerChars = "abcdefghijklmnopqrstuvwxyz"
	const upperChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var specialChars = "!@#$%&*()-_=+[]{}<>:?"
	var plan String

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	length := plan.Length.Value
	upper := plan.Upper.Value
	minUpper := plan.MinUpper.Value
	lower := plan.Lower.Value
	minLower := plan.MinLower.Value
	number := plan.Number.Value
	minNumeric := plan.MinNumeric.Value
	special := plan.Special.Value
	minSpecial := plan.MinSpecial.Value
	overrideSpecial := plan.OverrideSpecial.Value

	if overrideSpecial != "" {
		specialChars = overrideSpecial
	}

	var chars = string("")
	if upper {
		chars += upperChars
	}
	if lower {
		chars += lowerChars
	}
	if number {
		chars += numChars
	}
	if special {
		chars += specialChars
	}

	minMapping := map[string]int64{
		numChars:     minNumeric,
		lowerChars:   minLower,
		upperChars:   minUpper,
		specialChars: minSpecial,
	}

	var result = make([]byte, 0, length)

	for k, v := range minMapping {
		s, err := generateRandomBytes(&k, v)
		if err != nil {
			resp.Diagnostics.AddError(
				"error generating random bytes",
				fmt.Sprintf("error generating random bytes: %s", err),
			)
			return
		}
		result = append(result, s...)
	}

	s, err := generateRandomBytes(&chars, length-int64(len(result)))
	if err != nil {
		resp.Diagnostics.AddError(
			"error generating random bytes",
			fmt.Sprintf("error generating random bytes: %s", err),
		)
		return
	}

	result = append(result, s...)

	order := make([]byte, len(result))
	if _, err := rand.Read(order); err != nil {
		resp.Diagnostics.AddError(
			"error generating random bytes",
			fmt.Sprintf("error generating random bytes: %s", err),
		)
		return
	}

	sort.Slice(result, func(i, j int) bool {
		return order[i] < order[j]
	})

	str := String{
		ID:              types.String{Value: string(result)},
		Keepers:         plan.Keepers,
		Length:          types.Int64{Value: length},
		Special:         types.Bool{Value: special},
		Upper:           types.Bool{Value: upper},
		Lower:           types.Bool{Value: lower},
		Number:          types.Bool{Value: number},
		MinNumeric:      types.Int64{Value: minNumeric},
		MinUpper:        types.Int64{Value: minUpper},
		MinLower:        types.Int64{Value: minLower},
		MinSpecial:      types.Int64{Value: minSpecial},
		OverrideSpecial: types.String{Value: overrideSpecial},
		Result:          types.String{Value: string(result)},
	}

	if sensitive {
		str.ID.Value = "none"
	}

	diags = resp.State.Set(ctx, str)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func generateRandomBytes(charSet *string, length int64) ([]byte, error) {
	bytes := make([]byte, length)
	setLen := big.NewInt(int64(len(*charSet)))
	for i := range bytes {
		idx, err := rand.Int(rand.Reader, setLen)
		if err != nil {
			return nil, err
		}
		bytes[i] = (*charSet)[idx.Int64()]
	}
	return bytes, nil
}

func importString(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse, sensitive bool) {
	id := req.ID

	state := String{
		ID:     types.String{Value: id},
		Result: types.String{Value: id},
	}

	state.Keepers.ElemType = types.StringType

	if sensitive {
		state.ID.Value = "none"
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func validateLength(ctx context.Context, req tfsdk.ValidateResourceConfigRequest, resp *tfsdk.ValidateResourceConfigResponse) {
	var config String
	req.Config.Get(ctx, &config)

	length := config.Length.Value
	minUpper := config.MinUpper.Value
	minLower := config.MinLower.Value
	minNumeric := config.MinNumeric.Value
	minSpecial := config.MinSpecial.Value

	if length < minUpper+minLower+minNumeric+minSpecial {
		resp.Diagnostics.AddError(
			fmt.Sprintf("length (%d) must be >= min_upper + min_lower + min_numeric + min_special (%d)", length, minUpper+minLower+minNumeric+minSpecial),
			fmt.Sprintf("length (%d) must be >= min_upper + min_lower + min_numeric + min_special (%d)", length, minUpper+minLower+minNumeric+minSpecial),
		)
	}
}