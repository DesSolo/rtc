---
date: '2025-09-14T21:54:04+03:00'
draft: false
title: 'rtcconst'
weight: 3
---

`rtcconst` is a CLI utility designed to help developers by automatically generating typed constants (e.g., for the Go programming language) from your configuration YAML files. This ensures type safety and prevents typos when accessing config values in your code.

```text
Usage of const_generator:
  -output string
        path to output file (default "internal/config/config.go")
  -template string
        path to template file
  -yaml_description_key string
        key for description (default "usage")
  -yaml_path string
        path to config block in yaml file (default ".")
```

### Command Flags Explained

*   `--output`: Specifies the file path where the generated constants will be written. Default is `internal/config/config.go`.
*   `--template`: (Optional) The path to a custom Go template file. Use this for advanced control over the generated code's structure and format.
*   `--yaml_description_key`: The key in your YAML file that contains the human-readable description for a config key. The default is `usage`.
*   `--yaml_path`: A dot-separated path (e.g., `foo.bar`) pointing to the specific section *inside* the YAML file that contains your configuration keys. By default (`.`), it uses the entire file content.

### Usage Examples

**Example 1: Flat Structure**

This is for a YAML file that contains *only* your configuration keys at the top level.

```yaml {filename="values_flat.yaml"}
# example 1: flat struct
some_example_feature_key1:
  usage: This is my awesome feature description1
  some_another_sub_key: "generator ignore this"
```

To generate constants from this file, run:
```shell
rtcconst values_flat.yaml
```

**Example 2: Nested Structure**

This is for a YAML file where your configuration keys are nested inside other structures.

```yaml {filename="values_nested.yaml"}
# example 2: nested struct
foo:
  bar:
    baz:
      some_example_feature_key2:
      usage: This is my awesome feature description2
```

To generate constants, you must specify the path to the nested config block:
```shell
rtcconst --yaml_path foo.bar.baz values_nested.yaml
```

Using `rtcconst` helps you avoid hard-coded strings and makes your code more maintainable and less error-prone.