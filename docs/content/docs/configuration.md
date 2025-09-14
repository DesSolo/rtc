---
date: '2025-09-12T20:15:36+03:00'
draft: false
title: 'Configuration reference'
weight: 3
---

Actual configuration file in [examples](https://github.com/DesSolo/rtc/blob/master/examples/config.yaml)

## logging

Configures application logging settings.

```yaml
logging:
  level: -4
```

Supported log levels (from slog package):
| Level | Value |
|-------|-------|
| Debug | -4    |
| Info  | 0     |
| Warn  | 4     |
| Error | 8     |

## server

Main server configuration block.

### address

Defines the network address and port the server listens on.

```yaml
address: ":8080"
```

### read_header_timeout

Specifies maximum time to read request headers for security.

```yaml
read_header_timeout: 3s
```

### auth

More info about [auth]({{< ref "auth" >}})

Authentication configuration using JWT and static tokens.

#### jwt

JWT authentication using RSA256 asymmetric cryptography.

##### private_key

{{< callout type="warning" >}}
This parameter is sensitive
{{< /callout >}}

RSA private key for token signing (PEM format).

```yaml
private_key: |
  -----BEGIN PRIVATE KEY-----
  ...
  -----END PRIVATE KEY-----
```

##### public_key

RSA public key for token verification (PEM format).

```yaml
public_key: |
  -----BEGIN PUBLIC KEY-----
  ...
  -----END PUBLIC KEY-----
```

##### ttl

Token time-to-live duration.

```yaml
ttl: 24h
```

#### tokens

{{< callout type="warning" >}}
This parameter is sensitive
{{< /callout >}}

Static token authentication with role-based access.

```yaml
tokens:
  TOKEN_UUID:
    username: token_user
    roles: ["admin"]
```

### authorizer

Authorization system configuration.

#### kind

Authorization implementation type:
- `noop`: Allow all operations
- `rego`: Open Policy Agent Rego-based authorization [example policy](https://github.com/DesSolo/rtc/blob/master/examples/authz.rego)

```yaml
kind: rego
```

#### rego

Open Policy Agent configuration (required when kind=rego).

##### query

Rego policy decision query.

```yaml
query: data.authz.allow
```

##### policy_path

Path to Rego policy file.

```yaml
policy_path: examples/authz.rego
```

## storage

{{< callout type="warning" >}}
This parameter is sensitive
{{< /callout >}}

PostgreSQL database connection configuration.

```yaml
dsn: postgres://user:password@host:port/database?sslmode=disable
```

## values_storage

Etcd key-value storage configuration.

### endpoints

Etcd cluster connection endpoints.

```yaml
endpoints:
  - 127.0.0.1:2379
```

### dial_timeout

Connection establishment timeout.

```yaml
dial_timeout: 3s
```

### path

Root path for all keys in etcd.

```yaml
path: rtc
```