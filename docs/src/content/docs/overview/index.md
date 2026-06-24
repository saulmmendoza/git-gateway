---
title: Overview
description: Overview of Git Gateway
---

## What is Git Gateway?

Git Gateway lets you set up a gateway to your choice of Git provider's API (currently available with GitHub, GitLab, Bitbucket, and Forgejo) that lets tools work with content, branches, and pull requests on your users' behalf.

## How it works

The Git Gateway works with any identity service that can issue JWTs and only allows access when a JSON Web Token with sufficient permissions is present.

## Architecture

Git Gateway acts as a proxy between a client application and a Git provider. It validates JWT tokens and restricts access to specific safe endpoints.

## Comparison

Compared to direct API access, Git Gateway ensures credentials are kept safe and restricts user scopes down to just what is needed.
