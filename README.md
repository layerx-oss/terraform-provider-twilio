# Twilio Terraform Provider (layerx-oss fork)
[![Tests](https://github.com/layerx-oss/terraform-provider-twilio/actions/workflows/test-and-deploy.yml/badge.svg)](https://github.com/layerx-oss/terraform-provider-twilio/actions/workflows/test-and-deploy.yml)

## Project Status

This is an internal fork of [`twilio/terraform-provider-twilio`](https://github.com/twilio/terraform-provider-twilio),
maintained by LayerX for internal use. The upstream project is in PILOT and not actively maintained.

This fork is **not published to the public Terraform Registry**. It is distributed via GitHub Releases and consumed
through a Terraform [filesystem mirror](https://developer.hashicorp.com/terraform/cli/config/config-file#provider-installation)
(see "Installing and Using the Provider" below). The auto-generated resource code under `twilio/resources/**` is frozen
(not regenerated from upstream).

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) v0.15.x or later
- [Go](https://golang.org/doc/install) (see the version in `go.mod`) to build the provider plugin

## Resource Documentation

Documentation on the available resources that can be managed by this provider and their parameters can be found [here](twilio/resources/README.md).

Note that file upload resources are currently not available.

## Building The Provider

Clone repository:

```sh
git clone git@github.com:layerx-oss/terraform-provider-twilio
```

Enter the provider directory and build the provider:

```sh
make build
```

## Installing and Using the Provider

This fork uses the source address `registry.terraform.io/layerx/twilio` and is resolved through a local
filesystem mirror, so Terraform never contacts the public registry for it. Pick one of the two installation paths below.

### Option A: build and install locally (developers)

Run `make install`. This builds the provider and installs it into the filesystem mirror in unpacked layout:

```
~/.terraform.d/plugins/registry.terraform.io/layerx/twilio/<VERSION>/<OS>_<ARCH>/terraform-provider-twilio_v<VERSION>
```

### Option B: install a released build (consumers)

1. Download the zip for your OS/arch from the [GitHub Releases](https://github.com/layerx-oss/terraform-provider-twilio/releases).
2. Place the zip (do **not** unzip it) into the mirror directory in packed layout:

```
~/.terraform.d/plugins/registry.terraform.io/layerx/twilio/terraform-provider-twilio_<VERSION>_<OS>_<ARCH>.zip
```

### Configure the filesystem mirror

Add the following to your `~/.terraformrc` so Terraform resolves this provider from the mirror and never falls back to the
public registry:

```hcl
provider_installation {
  filesystem_mirror {
    path    = "~/.terraform.d/plugins"
    include = ["registry.terraform.io/layerx/*"]
  }
  direct {
    exclude = ["registry.terraform.io/layerx/*"]
  }
}
```

### Use it in your configuration

1. Configure the Twilio provider with your Twilio credentials in your Terraform configuration file (e.g. main.tf). These can also be set via `TWILIO_ACCOUNT_SID` and `TWILIO_AUTH_TOKEN` environment variables.
2. Add your resource configurations to your Terraform configuration file (e.g. main.tf).

```terraform
terraform {
  required_providers {
    twilio = {
      source  = "registry.terraform.io/layerx/twilio"
      version = "0.18.46"
    }
  }
}

# Credentials can be found at www.twilio.com/console.
provider "twilio" {
  //  username defaults to TWILIO_API_KEY with TWILIO_ACCOUNT_SID as the fallback env var
  //  password  defaults to TWILIO_API_SECRET with TWILIO_AUTH_TOKEN as the fallback env var
}

resource "twilio_api_accounts_keys" "key_name" {
  friendly_name = "terraform key"
}

output "messages" {
  value = twilio_api_accounts_keys.key_name
}
```

4. Run `terraform init` and `terraform apply` to initialize and apply changes to your Twilio infrastructure.

### Using environment variables

You can use credentials stored in environment variables for your setup:

#### OPTION 1 (recommended)
* `TWILIO_ACCOUNT_SID` = your Account SID from [your console](https://www.twilio.com/console)
* `TWILIO_API_KEY` = an API Key created in [your console](https://twil.io/get-api-key)
* `TWILIO_API_SECRET` = the secret for the API Key (you would have received this when you created an API key)
* _(optional)_ `TWILIO_REGION` = the Region for the account

#### OPTION 2
* `TWILIO_ACCOUNT_SID` = your Account SID from [your console](https://www.twilio.com/console)
* `TWILIO_AUTH_TOKEN` = your Auth Token from [your console](https://www.twilio.com/console)
* _(optional)_ `TWILIO_REGION` = the Region for the account

## Examples

For usage examples, checkout the [documentation in usage.md](usage.md) and the [examples folder](examples).

## Developing the Provider

The boilerplate includes the following:

- `Makefile` contains helper functions used to build, package and install the Twilio Terraform Provider. The `OS_ARCH` is resolved automatically from `go env`, so it works on both Intel and Apple Silicon machines.

  The `install` target builds the provider and installs it into the `~/.terraform.d/plugins/` filesystem mirror under `registry.terraform.io/layerx/twilio`. Configure `~/.terraformrc` as shown in "Installing and Using the Provider" and use `source = "registry.terraform.io/layerx/twilio"` in your Terraform configuration, then run `terraform init`.

- `examples` contains sample Terraform configuration that can be used to test the Twilio provider
- `twilio` contains the main provider code. This will be where the provider's resources and data source implementations will be defined.

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (the version pinned in `go.mod` is _required_).

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
make build
...
$GOPATH/bin/terraform-provider-twilio
...
```

In order to run the full suite of Acceptance tests, run `make testacc`. Provide your Account SID and Auth Token as environment variables to properly configure the test suite.

_Note:_ Acceptance tests create real resources, and often cost money to run.

```sh
 make testacc TWILIO_ACCOUNT_SID=YOUR_ACCOUNT_SID TWILIO_AUTH_TOKEN=YOUR_AUTH_TOKEN
```

You can also specify a particular suite to run like so:

```shell
 make testacc TEST=./twilio/ TWILIO_ACCOUNT_SID=YOUR_ACCOUNT_SID TWILIO_AUTH_TOKEN=YOUR_AUTH_TOKEN
```

An example test file can be found [here](https://github.com/twilio/terraform-provider-twilio/blob/main/twilio/resources_flex_test.go).

## Debugging

First:

```sh
export TF_LOG=TRACE
```

then refer to the [Terraform Debugging Documentation](https://www.terraform.io/docs/internals/debugging.html).

### Debugging with Delve

You can build and debug the provider locally. When using Goland you can set break point and step through code:

```sh
$ dlv debug main.go -- -debug
Type 'help' for list of commands.
(dlv) c
Provider started, to attach Terraform set the TF_REATTACH_PROVIDERS env var:

	TF_REATTACH_PROVIDERS='{"registry.terraform.io/layerx/twilio":{...}}}'
```

Copy the `TF_REATTACH_PROVIDERS` and run Terraform with this value set:

```sh
$ TF_REATTACH_PROVIDERS='...' terraform init
$ TF_REATTACH_PROVIDERS='...' terraform plan
...
```

Terraform will use the binary running under `dlv` instead of the `layerx/twilio` mirror version. For further details
refer to the [Terraform Debugging Providers](https://www.terraform.io/docs/extend/debugging.html) documentation.

### Debugging with Goland

- Set up GOROOT (initially opening `main.go` should show this option)
- Select `Modify Run Configuration...` on `main.go` and then add `--debug` as `Program arguments`
- Select `Debug "go build main.go"` and then copy the `TF_REATTACH_PROVIDERS` to the shell where `terraform` will be run
- Set breakpoints in Goland as needed and run `terraform`, it will use plugin the process running under the Goland debugger
