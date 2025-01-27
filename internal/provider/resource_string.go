package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceString and resourcePassword both use the same set of CustomizeDiffFunc(s) in order to handle the deprecation
// of the `number` attribute and the simultaneous addition of the `numeric` attribute. planDefaultIfAllNull handles
// ensuring that both `number` and `numeric` default to `true` when they are both absent from config.
// planSyncIfChange handles keeping number and numeric in-sync when either one has been changed.
func resourceString() *schema.Resource {
	customizeDiffFuncs := planDefaultIfAllNull(true, "number", "numeric")
	customizeDiffFuncs = append(customizeDiffFuncs, planSyncIfChange("number", "numeric"))
	customizeDiffFuncs = append(customizeDiffFuncs, planSyncIfChange("numeric", "number"))

	return &schema.Resource{
		Description: "The resource `random_string` generates a random permutation of alphanumeric " +
			"characters and optionally special characters.\n" +
			"\n" +
			"This resource *does* use a cryptographic random number generator.\n" +
			"\n" +
			"Historically this resource's intended usage has been ambiguous as the original example used " +
			"it in a password. For backwards compatibility it will continue to exist. For unique ids please " +
			"use [random_id](id.html), for sensitive random values please use [random_password](password.html).",
		CreateContext: createStringFunc(false),
		ReadContext:   readNil,
		DeleteContext: RemoveResourceFromState,
		// MigrateState is deprecated but the implementation is being left in place as per the
		// [SDK documentation](https://github.com/hashicorp/terraform-plugin-sdk/blob/main/helper/schema/resource.go#L91).
		MigrateState:  resourceRandomStringMigrateState,
		SchemaVersion: 2,
		Schema:        stringSchemaV2(),
		Importer: &schema.ResourceImporter{
			StateContext: importStringFunc,
		},
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 1,
				Type:    resourceStringV1().CoreConfigSchema().ImpliedType(),
				Upgrade: resourcePasswordStringStateUpgradeV1,
			},
		},
		CustomizeDiff: customdiff.All(
			customizeDiffFuncs...,
		),
	}
}

func importStringFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	val := d.Id()

	if err := d.Set("result", val); err != nil {
		return nil, fmt.Errorf("error setting result: %w", err)
	}

	return []*schema.ResourceData{d}, nil
}

func resourceStringV1() *schema.Resource {
	return &schema.Resource{
		Schema: stringSchemaV1(),
	}
}
