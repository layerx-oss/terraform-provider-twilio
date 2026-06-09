terraform {
  required_providers {
    twilio = {
      source  = "registry.terraform.io/layerx/twilio"
      version = ">=0.4.0"
    }
  }
}

provider "twilio" {
  //  username defaults to TWILIO_API_KEY with TWILIO_ACCOUNT_SID as the fallback env var
  //  password  defaults to TWILIO_API_SECRET with TWILIO_AUTH_TOKEN as the fallback env var
}

# Create a subaccount for a tenant.
# Twilio bills all subaccount usage to the parent account.
# There is no delete API: destroying this resource closes the subaccount (status=closed).
resource "twilio_api_accounts" "tenant_a" {
  friendly_name = "Tenant A"
}

# Create an API key scoped to the subaccount.
# The secret is returned only once (at creation) and is stored in Terraform state.
# Store it in a secrets manager (e.g. AWS Secrets Manager) rather than using it
# directly in other Terraform resources.
resource "twilio_api_accounts_keys" "tenant_a_runtime" {
  path_account_sid = twilio_api_accounts.tenant_a.sid
  friendly_name    = "tenant-a-runtime"
}

# Provision a phone number under the subaccount.
resource "twilio_api_accounts_incoming_phone_numbers" "tenant_a_number" {
  path_account_sid = twilio_api_accounts.tenant_a.sid
  area_code        = "415"
  friendly_name    = "Tenant A main number"
  voice_url        = "https://example.com/twiml/voice"
  sms_url          = "https://example.com/twiml/sms"

  # Prevent accidental release (and loss) of the provisioned phone number.
  lifecycle {
    prevent_destroy = true
  }
}

output "tenant_a_account_sid" {
  value = twilio_api_accounts.tenant_a.sid
}

output "tenant_a_api_key_sid" {
  description = "Use this as the Twilio API username in the runtime."
  value       = twilio_api_accounts_keys.tenant_a_runtime.sid
}

output "tenant_a_api_key_secret" {
  description = "Use this as the Twilio API password in the runtime. Store in Secrets Manager."
  value       = twilio_api_accounts_keys.tenant_a_runtime.secret
  sensitive   = true
}

output "tenant_a_phone_number" {
  value = twilio_api_accounts_incoming_phone_numbers.tenant_a_number.phone_number
}
