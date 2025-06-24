## Reverse Caching Proxy Workflow (Go + Redis)

A beginner level project to get started on Reverse Proxy and Caching. This proxy server fetches data from an origin API and caches it using Redis. It also respects a custom `Cache-Required` HTTP header for flexible cache control.

---

### Installation

1. Clone the project
2. Run ```go mod tidy```
3. Run redis-server using ```docker-compose up -d``` to run redis-db on detach mode
4. Either run **main.go** using ```go run main.go -origin <origin-url> -port <portnumber>``` or compile into binary using ```go build ```
5. Use the server on the port assigned (default port is 1234 🤓)

### Flow Overview

1. **Start the Server**
   - Accepts two flags:
     - `--port`: Port to listen on (default `1234`)
     - `--origin`: The origin base URL (e.g., `https://api.adviceslip.com`)

2. **Incoming HTTP Request**
   - The request path is appended to the origin URL.
   - Checks for the custom header: `Cache-Required`.

3. **Cache Behavior**

Some APIS like [random advice API](https://api.adviceslip.com/advice) return random responses on same endpoint. To resolve this we can append a custom flag to the request to get cached response or fresh response.
``` bash
curl -H "Cache-Required: (no^yes)" <request url>
```

#### If `Cache-Required` is not set to `no`:
   - Try to get the response from Redis using the full URL as the key.
   - **If found:**
     - Return cached response.
     - Add header: `X-Cache: Hit`
   - **If not found:**
     - Forward the request to the origin.
     - Cache the new response.
     - Add header: `X-Cache: Miss`

#### If `Cache-Required: no`:
   - Always send the request to the origin.
   - Cache the response (overwriting any existing value).
   - Add header: `X-Cache: Refreshed`

---

### Cache Logic Summary

| `Cache-Required` Header | Behavior                      | `X-Cache` Response |
|--------------------------|-------------------------------|--------------------|
| not `no`                   | Try cache, then fallback to origin | `Hit` or `Miss`    |
| `no` | Always fetch new response, update cache | `Refreshed`        |

---

### Redis Caching

- **Key**: Full request URL (origin + path)
- **Value**: Response body (as string)
- **TTL**: *(Optional)* Currently not set
- **Cache Clearing**: To clear cache we can send a get req to the proxy server on **/redigo/clear**, it will clear the cache
``` bash
curl <proxy-serverurl>/redigo/clear
```

### [Project URL](https://roadmap.sh/projects/caching-server)
