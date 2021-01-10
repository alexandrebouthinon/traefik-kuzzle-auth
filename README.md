<p align="center">
  <img src="https://user-images.githubusercontent.com/7868838/103894440-3a45f080-50ef-11eb-86a0-336682af6147.png"/>
</p>
<p align="center">
    <img alt="GitHub release (latest by date)" src="https://img.shields.io/github/v/release/alexandrebouthinon/traefik-kuzzle-auth">
    <a href="https://github.com/alexandrebouthinon/traefik-kuzzle-auth/actions?query=branch%3Amaster">
        <img alt="GitHub branch checks state" src="https://img.shields.io/github/checks-status/alexandrebouthinon/traefik-kuzzle-auth/master">
    </a>
    <a href="https://codecov.io/gh/alexandrebouthinon/traefik-kuzzle-auth">
        <img src="https://codecov.io/gh/alexandrebouthinon/traefik-kuzzle-auth/branch/master/graph/badge.svg?token=YZ8WADNYRH"/>
    </a>
    <a href="https://goreportcard.com/report/github.com/alexandrebouthinon/traefik-kuzzle-auth">
        <img src="https://goreportcard.com/badge/github.com/alexandrebouthinon/traefik-kuzzle-auth"/>
    </a>
    <a href="https://github.com/alexandrebouthinon/traefik-kuzzle-auth/blob/master/LICENSE">
        <img alt="GitHub" src="https://img.shields.io/github/license/alexandrebouthinon/traefik-kuzzle-auth">
    </a>
</p>

<!-- TOC -->

- [What?](#what)
- [Why?](#why)
- [How?](#how)
  - [Prerequisites](#prerequisites)
  - [Demo](#demo)
  - [Installation](#installation)
  - [Configuration](#configuration)
- [Development](#development)
- [Roadmap](#roadmap)
- [What is Kuzzle?](#what-is-kuzzle)

<!-- /TOC -->

# What?
This is a Traefik Basic Auth Plugin using Kuzzle as authentication provider.

# Why?

*One authentication system to rule them all* :sunglasses:

Kuzzle offer a complex and fine-grained RBAC authentication system, why do not use it everywhere? 


# How?
> :warning: At this time, Traefik Plugin system is still an experimental feature use it with caution. You can freeze your Traefik version to increase stability if you want to use this plugin on a real world use case

## Prerequisites

* A valid [Traefik Pilot](https://pilot.traefik.io) token for your Traefik instance.
* A running Kuzzle server in which one or more users are configured.


## Demo
You can found a demonstration Docker Compose file (`docker-compose.demo.yml`) in the repository root. 

```shell
TRAEFIK_PILOT_TOKEN="xxxx" docker-compose -f docker-compose.demo.yml up -d
```
This will launch:
* A complete [Kuzzle stack](http://localhost:7512) (Kuzzle, Elasticsearch and Redis containers).
* A Traefik instance with [dashboard](http://traefik.localhost) and latest released plugin version enabled and only available using `admin` Kuzzle user
* A [`whoami` instance](http://whoami.localhost) available using both `admin` and `developer` Kuzzle users

Once all containers are started and healthy, you can use the [Kuzzle Admin Console](https://next-console.kuzzle.io) to create your users (`admin` and `developer`).

## Installation
Declare it in the Traefik configuration:

**YAML**
```yaml
pilot:
  token: "xxxx"
experimental:
  plugins:
    traefik-kuzzle-auth:
        moduleName: github.com/alexandrebouthinon/traefik-kuzzle-auth
        version: v0.1.0
```

**TOML**
```toml
[pilot]
  token = "xxxx"
[experimental.plugins.fail2ban]
    moduleName = "github.com/alexandrebouthinon/traefik-kuzzle-auth"
    version = "v0.1.0"
```

**CLI**
```shell
--pilot.token=${TRAEFIK_PILOT_TOKEN}
--experimental.plugins.traefik-kuzzle-auth.moduleName=github.com/alexandrebouthinon/traefik-kuzzle-auth
--experimental.plugins.traefik-kuzzle-auth.version=v0.1.0
```

## Configuration

**YAML**
```yaml
middlewares:
  your-well-named-middleware:
    plugin:
      traefik-kuzzle-auth:
        customRealm: "Use a valid Kuzzle user to authenticate" # optional
        kuzzle:
          url: "http://localhost:7512" # required
          routes: # optional
            ping: /_publicApi
            login: /_login/local
            getCurrentUser: /_me # With Kuzzle v1 you must use '/users/_me'
          allowedUsers: # optional
            - admin
            - developer
```

**TOML**
```toml
[middlewares]
  [middlewares.your-well-named-middleware]
    [middlewares.your-well-named-middleware.plugin]
      [middlewares.your-well-named-middleware.plugin.traefik-kuzzle-auth]
        customRealm = "Use a valid Kuzzle user to authenticate" # optional
        
        [middlewares.your-well-named-middleware.plugin.traefik-kuzzle-auth.kuzzle]
          url = "http://localhost:7512" # required
          allowedUsers = ["admin", "developer"] # optional

          [middlewares.your-well-named-middleware.plugin.traefik-kuzzle-auth.kuzzle.routes] # optional
            ping = "/_publicApi"
            login = "/_login/local"
            getCurrentUser = "/_me" # With Kuzzle v1 you must use '/users/_me'

```

**Docker Compose Labels**
```yaml
labels:
  - "traefik.http.middlewares.your-well-named-middleware.plugin.traefik-kuzzle-auth.customRealm=Use a valid Kuzzle user to authenticate" # optional
  - "traefik.http.middlewares.your-well-named-middleware.plugin.traefik-kuzzle-auth.kuzzle.url=http://kuzzle:7512" # required
  - "traefik.http.middlewares.your-well-named-middleware.plugin.traefik-kuzzle-auth.kuzzle.routes.ping=/_publicApi" # optional
  - "traefik.http.middlewares.your-well-named-middleware.plugin.traefik-kuzzle-auth.kuzzle.routes.login=/_login/local" # optional
  - "traefik.http.middlewares.your-well-named-middleware.plugin.traefik-kuzzle-auth.kuzzle.routes.getCurrentUser=/_me" # With Kuzzle v1 you must use '/users/_me' (optional)
  - "traefik.http.middlewares.kuzzle-auth.plugin.traefik-kuzzle-auth.kuzzle.allowedUsers=admin,developer" # optional
```

# Development
You can found a development Docker Compose file (`docker-compose.dev.yml`) in the repository root. 

```shell
TRAEFIK_PILOT_TOKEN="xxxx" docker-compose -f docker-compose.dev.yml up -d
```
This will launch:
* A complete [Kuzzle stack](http://localhost:7512) (Kuzzle, Elasticsearch and Redis containers).
* A Traefik instance with [dashboard](http://traefik.localhost) and latest released plugin version enabled and only available using `admin` Kuzzle user
* A [`whoami` instance](http://whoami.localhost) available using both `admin` and `developer` Kuzzle users

Once all containers are started and healthy, you can use the [Kuzzle Admin Console](https://next-console.kuzzle.io) to create your users (`admin` and `developer`).

# Roadmap

- [x] [Users](https://docs.kuzzle.io/core/2/guides/main-concepts/permissions/#users) greenlisting
- [ ] [Profiles](https://docs.kuzzle.io/core/2/guides/main-concepts/permissions/#profiles) greenlisting
- [ ] [Kuzzle API Key](https://docs.kuzzle.io/core/2/guides/advanced/api-keys/) authentication

New ideas are welcome, feel free to fill out an issue and let's discuss it :wink:

# What is Kuzzle?

Kuzzle is a [generic backend](https://docs.kuzzle.io/core/2/guides/introduction/general-purpose-backend/) offering **the basic building blocks common to every application**.

Rather than developing the same standard features over and over again each time you create a new application, Kuzzle proposes them off the shelf, allowing you to focus on building **high-level, high-value business functionalities**.

Kuzzle enables you to build modern web applications and complex IoT networks in no time.

* **API First**: use a standardised multi-protocol API.
* **Persisted Data**: store your data and perform advanced searches on it.
* **Realtime Notifications**: use the pub/sub system or subscribe to database notifications.
* **User Management**: login, logout and security rules are no more a burden.
* **Extensible**: develop advanced business feature directly with the integrated framework.
* **Client SDKs**: use our SDKs to accelerate the frontend development.


