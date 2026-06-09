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

# Create an API key under the provider's account.
# The secret is returned by Twilio only at creation time and is stored in
# Terraform state as a sensitive value.  Import does not restore the secret —
# use the key SID + secret from the original apply output (or a secrets manager).
resource "twilio_api_accounts_keys" "runtime" {
  friendly_name = "runtime-key"
}

output "api_key_sid" {
  description = "Use as the Twilio API username."
  value       = twilio_api_accounts_keys.runtime.sid
}

output "api_key_secret" {
  description = "Use as the Twilio API password. Store in Secrets Manager."
  value       = twilio_api_accounts_keys.runtime.secret
  sensitive   = true
}
