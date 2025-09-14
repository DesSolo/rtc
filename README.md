# 🚀 RTC - Real-time Configuration Management

Welcome to RTC, your solution for dynamic and efficient configuration management! 🎉

RTC provides a robust platform to manage your application configurations in real-time, allowing for seamless updates and deployments without service interruptions.

## ✨ Key Features

- **Dynamic Configuration Updates:** Modify configurations on the fly and have your applications react instantly.
- **Centralized Management:** A single source of truth for all your application settings across different environments.
- **Environment-Specific Configurations:** Easily manage configurations tailored for development, staging, and production environments.
- **User Management & Audit Logs:** Securely manage users and monitor all configuration changes with detailed audit logs.

## ⚠️ Important Note on Configuration Size

Currently, each configuration is limited to **128 keys**. This limitation is due to the underlying `etcd` storage. We are actively exploring solutions to overcome this, but for now, please keep your configurations concise.

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

---

## 📚 Learn More

Dive deeper into RTC with our detailed examples and documentation:

- **Server Configuration:** Learn how to configure `rtcserver` to fit your needs.
  - [Example Configuration File](examples/config.yaml)
- **Client Library:** Integrate RTC into your applications for real-time config updates.
  - [Client Library Guide](examples/client/README.md)
- **Code Generator (`rtcconst`):** Generate type-safe constants for your configuration keys to avoid typos and errors.
  - [rtcconst Generator Guide](examples/const_generator/README.md)
- **Command Line Tool (`rtcctl`):** Manage configurations, users, and permissions directly from your terminal.
  - [Complete rtcctl Documentation](examples/rtcctl/README.md)

We hope you enjoy using RTC! If you have any questions or feedback, feel free to reach out. 🌟

---

## Local run

Dependencies:
- docker
- golang
- npm

```bash
git clone https://github.com/DesSolo/rtc.git
cd rtc
# start storages (postgress and etcd)
docker-compose run -d -f docker-compose.d/docker-compose.yaml
# installing dependencies for go (goose, linter, releaser)
make install-deps
# apply sql migrations
make migrations-up
# start api server
make run
# start frontend
make run-ui
```