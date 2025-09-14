---
date: '2025-09-14T22:31:26+03:00'
draft: true
title: 'Workflow'
weight: 3
---

This document outlines the recommended workflow for managing configuration changes using the RTC system.

1.  **Project Registration**
    An administrator (or an automated process) creates and registers a new project within the RTC system.

2.  **Define Configuration**
    A developer defines the application's configuration structure and default values in a `values.yaml` file.

3.  **Generate Constants**
    The developer uses the [rtcconst]({{< ref "ecosystem/rtcconst" >}}) tool to generate type-safe constants from the `values.yaml` file for use in the application's codebase. This ensures all configuration keys are accessed safely without typos.

4.  **Commit to Source Control**
    The developer commits the application code, along with the `values.yaml` file and generated constants, to the Git repository.

5.  **CI/CD Configuration Update**
    The project's CI/CD pipeline uses the [rtcctl]({{< ref "ecosystem/rtcctl/config" >}}) tool to upsert the configuration defined in `values.yaml` to the RTC server *before* deploying the main application. This ensures the new configuration is live and available.

6.  **Application Startup**
    The developer's service application starts up, connects to the RTC server, and immediately receives the new, up-to-date configuration parameters.