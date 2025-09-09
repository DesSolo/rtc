# const generator examples
this is examples for usage const generator in your project

```shell
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

basic usage [values_flat.yaml](values_flat.yaml)
```shell
const_generator values_flat.yaml
```

nested yaml struct [values_nested.yaml](values_nested.yaml)
```shell
const_generator --yaml_path foo.bar values_nested.yaml
```