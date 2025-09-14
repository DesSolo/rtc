---
date: '2025-09-14T21:55:37+03:00'
draft: true
title: 'config'
---

## The `config upsert` Command

This command is used to create or update (`upsert` = update + insert) configuration values for a specific project, environment, and release.

```shell
rtcctl config upsert --project example --env dev --release latest --values examples/rtccrl/values.yaml
```

### Command Arguments

*   `--project`: The name of the project. **This project must already exist in the system.**
*   `--env`: The target environment (e.g., `dev`, `staging`, `prod`).
*   `--release`: The name of the release. This is often a feature branch name, a ticket number, or a version tag (e.g., `feature-x`, `v1.2.5`).
*   `--values`: The path to the YAML file containing the configuration keys and values.

### Configuration File Example

Here is an example of what the `values.yaml` file should look like:

```yaml {filename="values.yaml"}
example_key:
  usage: This is example_key usage  # A description of what this config key is for.
  group: "tech"                     # A logical group for organization and UI display.
  value: "some string value"        # The current value to be set.
  type: string                      # The data type (e.g., string, int, bool).
  writable: true                    # Whether this value can be edited later via the UI.
```

### Configuration Key Fields Explained

Each key in the YAML file (like `example_key` above) is a configuration parameter. Its properties are defined as follows:

*   **`usage`**: (Required) A human-readable description of the parameter's purpose.
*   **`group`**: (Optional) A category name (e.g., `tech`, `business`, `database`) used to logically group parameters together in the UI.
*   **`value`**: (Required) The actual value you want to set for this parameter.
*   **`type`**: (Required) The data type of the value. Common types include `string`, `int`, `integer`, `float`, `bool`, `boolean`.
*   **`writable`**: (Optional, default: `false`) A boolean (`true`/`false`) defining if this value can be modified after creation through the web UI. Use this to lock down critical configuration values that should only be changed via controlled processes (like CI/CD).