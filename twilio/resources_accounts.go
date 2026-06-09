package twilio

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"

	twilioclient "github.com/twilio/terraform-provider-twilio/client"
	"github.com/twilio/terraform-provider-twilio/core"
	"github.com/twilio/terraform-provider-twilio/twilio/resources"
)

// buildResourcesMap returns the generated resource map with the hand-written
// additions and overrides applied on top.
func buildResourcesMap() map[string]*schema.Resource {
	m := resources.NewTwilioResources().Map
	for name, r := range handwrittenResources() {
		m[name] = r
	}
	return m
}

// handwrittenResources holds resources that are added to, or override, the
// OpenAPI-generated resource map. They live outside the frozen generated code
// and are merged in by buildResourcesMap.
//
// Why these exist:
//   - twilio_api_accounts: the generator does not emit a resource for managing
//     (sub)accounts, but multi-tenant provisioning needs to create one Twilio
//     subaccount per tenant. (upstream issue twilio/terraform-provider-twilio#99)
//   - twilio_api_accounts_keys: the generated resource drops the API key secret,
//     which is only returned on creation, so the created key is unusable. This
//     override captures it. (upstream issue #82)
func handwrittenResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"twilio_api_accounts":      resourceAccount(),
		"twilio_api_accounts_keys": resourceAccountsKeysWithSecret(),
	}
}

// resourceAccount manages a Twilio (sub)account. Twilio does not allow deleting
// an account through the API, so Delete closes it instead.
func resourceAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: createAccount,
		ReadContext:   readAccount,
		UpdateContext: updateAccount,
		DeleteContext: deleteAccount,
		Schema: map[string]*schema.Schema{
			"friendly_name":     core.AsString(core.SchemaComputedOptional),
			"status":            core.AsString(core.SchemaComputedOptional),
			"sid":               core.AsString(core.SchemaComputed),
			"owner_account_sid": core.AsString(core.SchemaComputed),
			// Only returned at creation/fetch; treat as sensitive.
			"auth_token": core.AsString(core.SchemaComputedSensitive),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func createAccount(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	params := openapi.CreateAccountParams{}
	if err := core.UnmarshalSchema(&params, d); err != nil {
		return diag.FromErr(err)
	}

	r, err := m.(*twilioclient.Config).Client.Api.CreateAccount(&params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*r.Sid)
	if err := core.MarshalSchema(d, r); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func readAccount(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	r, err := m.(*twilioclient.Config).Client.Api.FetchAccount(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if err := core.MarshalSchema(d, r); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func updateAccount(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	params := openapi.UpdateAccountParams{}
	if err := core.UnmarshalSchema(&params, d); err != nil {
		return diag.FromErr(err)
	}

	r, err := m.(*twilioclient.Config).Client.Api.UpdateAccount(d.Id(), &params)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := core.MarshalSchema(d, r); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

// deleteAccount closes the (sub)account; Twilio has no delete endpoint.
func deleteAccount(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	closed := "closed"
	params := openapi.UpdateAccountParams{Status: &closed}
	if _, err := m.(*twilioclient.Config).Client.Api.UpdateAccount(d.Id(), &params); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

// resourceAccountsKeysWithSecret overrides the generated twilio_api_accounts_keys
// resource so the secret returned only at creation is captured into state.
func resourceAccountsKeysWithSecret() *schema.Resource {
	return &schema.Resource{
		CreateContext: createAccountsKeysWithSecret,
		ReadContext:   readAccountsKeysWithSecret,
		UpdateContext: updateAccountsKeysWithSecret,
		DeleteContext: deleteAccountsKeysWithSecret,
		Schema: map[string]*schema.Schema{
			"path_account_sid": core.AsString(core.SchemaComputedOptional),
			"friendly_name":    core.AsString(core.SchemaComputedOptional),
			"sid":              core.AsString(core.SchemaComputed),
			// Returned only when the key is first created.
			"secret": core.AsString(core.SchemaComputedSensitive),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func createAccountsKeysWithSecret(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	params := openapi.CreateNewKeyParams{}
	if err := core.UnmarshalSchema(&params, d); err != nil {
		return diag.FromErr(err)
	}

	r, err := m.(*twilioclient.Config).Client.Api.CreateNewKey(&params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*r.Sid)
	// r is *ApiV2010NewKey and includes Secret, so MarshalSchema captures it.
	if err := core.MarshalSchema(d, r); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func readAccountsKeysWithSecret(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	params := openapi.FetchKeyParams{}
	if err := core.UnmarshalSchema(&params, d); err != nil {
		return diag.FromErr(err)
	}

	// The fetch response (ApiV2010Key) has no Secret field, so MarshalSchema
	// leaves the secret already stored in state untouched.
	r, err := m.(*twilioclient.Config).Client.Api.FetchKey(d.Id(), &params)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := core.MarshalSchema(d, r); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func updateAccountsKeysWithSecret(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	params := openapi.UpdateKeyParams{}
	if err := core.UnmarshalSchema(&params, d); err != nil {
		return diag.FromErr(err)
	}

	r, err := m.(*twilioclient.Config).Client.Api.UpdateKey(d.Id(), &params)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := core.MarshalSchema(d, r); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func deleteAccountsKeysWithSecret(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	params := openapi.DeleteKeyParams{}
	if err := core.UnmarshalSchema(&params, d); err != nil {
		return diag.FromErr(err)
	}

	if err := m.(*twilioclient.Config).Client.Api.DeleteKey(d.Id(), &params); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
