# watchdog

[![Build Status](https://github.com/ailinykh/watchdog/actions/workflows/build.yml/badge.svg)](https://github.com/ailinykh/watchdog/actions/workflows/build.yml)

simple service for monitoring your [Amazon Lightsail](https://lightsail.aws.amazon.com/) instance

### how to run

Create `.env` file from a template:

```bash
cp .env.sample .env
```

paste some credentials into `.env` file.

Run the app
```bash
make run
```