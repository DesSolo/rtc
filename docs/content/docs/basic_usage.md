---
date: '2025-09-12T20:05:38+03:00'
draft: false
title: 'Basis usage'
weight: 1
---


## 🚦 Quick Start: Run RTC Locally

Let's get you up and running with a local development instance of RTC in a few simple steps.

### Prerequisites
1. **Docker & Docker Compose:** Needed to run the database.
2. **Goose:** A database migration tool. [Download it here](https://github.com/pressly/goose/releases).
3. **RTC Binaries:** Download the latest `rtcserver` (the server) and `rtcctl` (the CLI control tool) from our [Releases page](https://github.com/DesSolo/rtc/releases).

### Step-by-Step Guide

1. **Start the Database:**
```bash
# Download our docker-compose file and run it
curl -O https://github.com/DesSolo/rtc/blob/master/docker-compose.d/docker-compose.yaml
docker-compose up -d
 ```

This will start a PostgreSQL database in the background.

2. **Setup the Database Schema:**

```bash
# Download the migrations directory from the repository
# Then, apply the migrations using Goose
goose -dir ./migrations postgres "host=localhost port=5432 user=postgres password=postgres dbname=rtc sslmode=disable" up
```

This command creates all the necessary tables in your database.

3. **Configure and Run the Server:**

```bash
# Get the example server configuration file
curl -O https://github.com/DesSolo/rtc/blob/master/examples/config.yaml
# Launch the server
./rtcserver
```

4. **You're all set!** 🎉
Open your browser and go to `http://localhost:8080/ui` to access the RTC management interface.

Default credentials: username `admin` password `rtc`
