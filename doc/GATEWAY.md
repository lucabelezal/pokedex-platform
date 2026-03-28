# Gateway

## Purpose

Kong is the platform entry point. It provides a single public access layer in front of the internal applications.

## Current Role

- Receive external traffic.
- Route requests to `mobile-bff`.
- Hide internal service topology from clients.

## Current Configuration

The repository uses a declarative Kong configuration in `core/gateway/kong/kong.yml`.

At the moment:

- Kong exposes `/v1`.
- Requests are forwarded to `http://mobile-bff:8080/api`.

## Architectural Notes

The gateway is intentionally thin right now. That is a good default for an early-stage project because it keeps routing simple while the domain is still evolving.

## Improvement Opportunities

- Document public routes and route ownership more explicitly.
- Add request tracing or correlation headers between gateway and downstream services.
- Decide whether authentication should remain fully in the BFF or partially move to the gateway in the future.
- Add rate limiting and standardized error transformation if the platform becomes externally exposed.

