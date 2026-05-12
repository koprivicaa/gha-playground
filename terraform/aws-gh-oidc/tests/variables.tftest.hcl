run "variables_follow_conventions" {
  command = plan

  assert {
    condition     = can(regex("^https://", var.github_oidc_url))
    error_message = "github_oidc_url must start with https://"
  }

  assert {
    condition     = trim(var.github_oidc_audience, " ") != ""
    error_message = "github_oidc_audience must not be empty."
  }

  assert {
    condition     = can(regex("^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$", var.github_repository))
    error_message = "github_repository must be in owner/repository format."
  }

  assert {
    condition     = trim(var.github_branch, " ") != ""
    error_message = "github_branch must not be empty."
  }

  assert {
    condition     = trim(var.github_actions_role_name, " ") != ""
    error_message = "github_actions_role_name must not be empty."
  }

  assert {
    condition     = trim(var.github_actions_policy_name, " ") != ""
    error_message = "github_actions_policy_name must not be empty."
  }

  assert {
    condition = length(var.github_actions_allowed_actions) > 0 && alltrue([
      for action in var.github_actions_allowed_actions : trim(action, " ") != ""
    ])
    error_message = "github_actions_allowed_actions must contain at least one non-empty action."
  }
}
