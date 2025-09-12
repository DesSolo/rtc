---
date: '2025-09-12T21:21:56+03:00'
draft: false
title: 'Rego'
---

The RTC server supports authorization using [Rego](https://www.openpolicyagent.org/docs/policy-language) policies, the language used by the Open Policy Agent (OPA). This gives you fine-grained control over who can do what within the system.

```rego {base_url="https://github.com/DesSolo/rtc/blob/master/examples/authz.rego",filename="authz.rego"}
package authz

default allow := false

# allow any for admins
allow if {
	"admin" in input.user.roles
}

# allow users endpoints
allow if {
	some rule in users_endpoints
	glob.match(rule.path, ["/"], input.path)
	input.method in rule.methods
}

# allow specific user for change config values
allow if {
    input.method == "PUT"
	glob.match("/api/v1/projects/*/envs/*/releases/*/configs", ["/"], input.path)
	project := project_name(input.path)
	input.user.username in users_configs_changes[project]
}

project_name(path) := project_name if {
	parts := split(path, "/")
	project_name := parts[4]
}

users_endpoints := [
	{"path": "/api/v1/projects", "methods": {"GET"}},
	{"path": "/api/v1/projects/*/envs", "methods": {"GET"}},
	{"path": "/api/v1/projects/*/envs/*/releases", "methods": {"GET"}},
	{"path": "/api/v1/projects/*/envs/*/releases/*/configs", "methods": {"GET"}},
	{"path": "/api/v1/audits", "methods": {"GET"}},
	{"path": "/api/v1/audits/actions", "methods": {"GET"}},
]

users_configs_changes := {
    "example": {"simple_user"},
}
```

### What This Policy Does

This example policy demonstrates a common authorization structure:

*   **Full Access for Admins:** Any user with the `admin` role is allowed to perform any action.
*   **Read Access for Everyone:** All authenticated users can view projects, environments, releases, configurations, and audit logs.
*   **Write Access for Specific Users:** Only the user `simple_user` is permitted to update configuration values (using a `PUT` request) for the project named `example`.

{callout}
You have full flexibility to customize this policy to match your organization's specific security requirements and workflows.
{/callout}

### How to Implement Your Policy

Follow these simple steps to get started with custom Rego policies:

1.  **Create your policy file:** Write your rules and save them in a file (e.g., `authz.rego`).
2.  **Update the configuration:** Point to your new policy file in the `config.yaml`.
3.  **Restart the server:** Apply the changes by restarting the RTC server.

```yaml {linenos=table,hl_lines=[2],filename="config.yaml"}
  authorizer:
    kind: rego
    rego:
      query: data.authz.allow
      policy_path: examples/authz.rego
```