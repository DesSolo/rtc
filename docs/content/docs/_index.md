---
date: '2025-09-12T20:04:44+03:00'
draft: false
title: ''
---

# Welcome to RTC documentation

> [!WARNING]
> Currently, each configuration is limited to **128 keys**. This limitation is due to the underlying `etcd` storage. We are actively exploring solutions to overcome this, but for now, please keep your configurations concise.

## Principal diagram

```mermaid
sequenceDiagram
    actor User
    participant RTC
    participant DB
    participant etcd
    actor Client_Code

    Note over User, etcd: Configuration recording process
    User->>RTC: Sending key and value
    RTC->>DB: Saving a key with metadata
    RTC->>etcd: Saving a key with value
    etcd-->>RTC: Confirmation of entry
    DB-->>RTC: Confirmation of entry
    RTC-->>User: Successful update

    Note over Client_Code, etcd: Configuration reading process
    Client_Code->>etcd: Query value by key
    etcd-->>Client_Code: Return value

```

Explore the following sections to learn how to use RTC:

{{< cards >}}
  {{< card link="basic_usage" title="Basic usage" icon="document-text" >}}
  {{< card link="configuration" title="Configuration reference" icon="adjustments" >}}
{{< /cards >}}
