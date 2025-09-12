---
date: '2025-09-12T20:13:11+03:00'
draft: false
title: 'Authorization'
---

The RTC server supports multiple authorization methods to fit different needs:
*   **JWT** for user access.
*   **Token** for services and automated scripts.

## JWT

For seamless integration into your existing infrastructure, we use asymmetric cryptography (a public/private key pair) for signing and verifying JWT tokens. This allows your existing auth system to generate valid tokens that the RTC server can verify using only the public key.

**Example JWT payload:**
```json
{"Username":"admin","Roles":["admin"],"iss":"rtc","exp":1757756630,"iat":1757670230}
```

- *Username*: The user's identifier. This name appears in audit logs and on the user interface (UI).
- *Roles*: A list of the user's roles. These are used for permission checks in Rego policies and to control access to specific UI sections (currently, only the presence of the admin role is checked for UI access).

To use JWT authorization, include the following header in your requests: `Authorization: jwt ${TOKEN}` where `${TOKEN}` is your generated JWT string.

> [!NOTE]
> Remember to configure the public key in your `config.yaml` file for the server to verify the tokens.

## Token
This method is designed for service accounts and automation, such as CI/CD pipelines that need to update configuration files programmatically.

To authorize, use the header: `Authorization: token ${TOKEN}` where `${TOKEN}` is a secret key defined in your config.yaml file.

> [!NOTE]
> Role-based policies apply to tokens as well. Ensure you assign the correct roles to your tokens in the `config.yaml` configuration.

> [!WARNING]
> Tokens are stored in plain text within the configuration file. Please take necessary measures to protect this file and manage the security of your tokens.