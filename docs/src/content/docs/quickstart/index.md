---
title: Quickstart
description: Get started with Git Gateway.
---

## Basic Usage

Git Gateway provides secure role-based access to the APIs of common Git Hosting providers.

## Docker Compose

```yaml
version: '3'
services:
  git-gateway:
    image: netlify/git-gateway
    ports:
      - "8081:8081"
    environment:
      - GITGATEWAY_GITHUB_ACCESS_TOKEN=your_token
      - GITGATEWAY_JWT_SECRET=your_jwt_secret
```

## Examples

To run Git Gateway locally:

```bash
git-gateway -c config.yml
```
