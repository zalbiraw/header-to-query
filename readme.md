# Header To Query Plugin

A Traefik plugin that converts HTTP headers to URL query parameters. Supports mapping, renaming, and optionally keeping headers. Handles multiple headers with the same name.

## Installation & Enabling

Add the plugin to your Traefik static configuration:

```yaml
experimental:
  plugins:
    headertoquery:
      moduleName: github.com/zalbiraw/headertoquery
      version: v0.0.1
```

## Dynamic Configuration Example

```yaml
http:
  routers:
    my-router:
      rule: host(`demo.localhost`)
      service: service-foo
      entryPoints:
        - web
      middlewares:
        - headertoquery

  middlewares:
    headertoquery:
      plugin:
        headers:
          - name: SERVICE_TAG
            key: id
          - name: RANK
          - name: GROUP
            keepHeader: true
```

## How It Works

- Each `headers` entry can specify:
  - `name`: The HTTP header to process
  - `key`: (Optional) The query parameter name to use (defaults to the header name)
  - `keepHeader`: (Optional) If `true`, the header is not removed from the request
- If a header appears multiple times, all values are mapped as repeated query parameters (e.g., `?id=1&id=2`).

### Example

Given this configuration:

```yaml
headers:
  - name: SERVICE_TAG
    key: id
  - name: RANK
  - name: GROUP
    keepHeader: true
```

And a request with these headers:

```
SERVICE_TAG: S117
SERVICE_TAG: SPARTAN-117
SERVICE_TAG: 117
RANK: Masterchief
GROUP: UNSC
```

The resulting query string will be:

```
?id=S117&id=SPARTAN-117&id=117&rank=Masterchief&group=UNSC
```

And the resulting headers will be:

```
GROUP: UNSC
```

The `SERVICE_TAG` and `RANK` headers are removed; `GROUP` remains because `keepHeader: true`.

## Development & Testing

Run tests with:

```sh
go test -v
```

---

For more details, see the source code and test cases.
