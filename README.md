# Traefik Plugin: Basic Auth powered by Kuzzle
Traefik Plugin using Kuzzle as authentication provider for Basic Auth middleware

## Why?

One authentication system to rule them all. Kuzzle offer a complex and fine-grained RBAC authentication system, why do not use it everywhere? 

## How it works?

This plugin rely on [Basic Auth Traefik middleware](https://doc.traefik.io/traefik/middlewares/basicauth/).
The principle is rather simple, configure the middleware with a single user/password pair by following its [documentation](https://doc.traefik.io/traefik/middlewares/basicauth/), add this information to the plugin configuration as well as the connection info to the Kuzzle server (`host`, `port`...) and.... that's it. Enjoy going to your applications hidden by Basic Auth using your Kuzzle user! :tada:

## What is Kuzzle?

Kuzzle is a [generic backend](https://docs.kuzzle.io/core/2/guides/introduction/general-purpose-backend/) offering **the basic building blocks common to every application**.

Rather than developing the same standard features over and over again each time you create a new application, Kuzzle proposes them off the shelf, allowing you to focus on building **high-level, high-value business functionalities**.

Kuzzle enables you to build modern web applications and complex IoT networks in no time.

* **API First**: use a standardised multi-protocol API.
* **Persisted Data**: store your data and perform advanced searches on it.
* **Realtime Notifications**: use the pub/sub system or subscribe to database notifications.
* **User Management**: login, logout and security rules are no more a burden.
* **Extensible**: develop advanced business feature directly with the integrated framework.
* **Client SDKs**: use our SDKs to accelerate the frontend development.


