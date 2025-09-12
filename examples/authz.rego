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

