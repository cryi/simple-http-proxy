# Simple Http Proxy with header injection

We use it to bypass CORS restrictions on some of our setups where it is not possible to inject headers on server side

```
Usage:
  shp [flags]

Flags:
  -f, --forward string              Address to forward requests to.
  -h, --help                        help for shp
  -i, --inject-header stringArray   List of headers to inject in format '<header>:<value>'
  -l, --listen string               Address to listen on.
```