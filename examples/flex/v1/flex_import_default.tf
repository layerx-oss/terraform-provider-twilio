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

resource "twilio_chat_services_v2" "flex_chat_service" {
  friendly_name = "Flex Chat Service"
}

resource "twilio_taskrouter_workspaces_v1" "flex_task_assignment" {
  friendly_name = "Flex Task Assignment"
}

resource "twilio_taskrouter_workspaces_task_queues_v1" "everyone" {
  workspace_sid = twilio_taskrouter_workspaces_v1.flex_task_assignment.sid
  friendly_name = "Everyone"
}

resource "twilio_taskrouter_workspaces_workflows_v1" "assign_to_anyone" {
  workspace_sid = twilio_taskrouter_workspaces_v1.flex_task_assignment.sid
  friendly_name = "Assign To Anyone"
  configuration = jsonencode({
    task_routing = {
      default_filter = {
        queue = twilio_taskrouter_workspaces_task_queues_v1.everyone.sid
      }
    }
  })
}

resource "twilio_taskrouter_workspaces_activities_v1" "offline" {
  workspace_sid = twilio_taskrouter_workspaces_v1.flex_task_assignment.sid
  friendly_name = "Offline"
}

resource "twilio_taskrouter_workspaces_activities_v1" "available" {
  workspace_sid = twilio_taskrouter_workspaces_v1.flex_task_assignment.sid
  friendly_name = "Available"
}

resource "twilio_taskrouter_workspaces_activities_v1" "unavailable" {
  workspace_sid = twilio_taskrouter_workspaces_v1.flex_task_assignment.sid
  friendly_name = "Unavailable"
}

resource "twilio_taskrouter_workspaces_activities_v1" "break" {
  workspace_sid = twilio_taskrouter_workspaces_v1.flex_task_assignment.sid
  friendly_name = "Break"
}
