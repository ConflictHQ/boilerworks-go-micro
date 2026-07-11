# Security Policy

## Reporting a Vulnerability

If you discover a security vulnerability in Boilerworks, please report it responsibly.

**Do not open a public issue.**

Instead, email **security@weareconflict.com** with:

- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

We will acknowledge your report within 48 hours and aim to release a fix within 7 days for critical issues.

## Supported Versions

| Version | Supported |
| ------- | --------- |
| latest  | Yes       |

## Security Best Practices

When deploying Boilerworks:

- Rotate the seed API key: set `API_KEY_SEED` to a strong random value (never ship `bw_seed_key_change_me_in_production`), then mint scoped keys via `/api-keys` and revoke the seed key
- Change the default Postgres credentials in `DATABASE_URL` (and `docker-compose.yml` if you deploy with it)
- Use HTTPS in production -- API keys travel in the `X-API-Key` header
- Grant keys the narrowest scopes that work (`events.read`, `events.write`, `keys.manage`); avoid the `*` wildcard
- Keep the built-in rate limiter enabled
- Review the security hardening in `bootstrap.md`
