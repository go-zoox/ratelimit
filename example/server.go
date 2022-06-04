package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-zoox/ratelimit"
)

func main() {
	limiter := ratelimit.NewMemory("example", 10*time.Second, 2)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ip := strings.Split(r.RemoteAddr, ":")[0]
		limiter.Inc(ip)

		if limiter.IsExceeded(ip) {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		w.Header().Set("X-RateLimit-Remaing", fmt.Sprintf("%d", limiter.Remaining(ip)))
		w.Header().Set("X-RateLimit-Reset-After", fmt.Sprintf("%d", limiter.ResetAfter(ip)))
		w.Header().Set("X-RateLimit-Total", fmt.Sprintf("%d", limiter.Total(ip)))

		w.Write([]byte("Hello World!"))
	})

	fmt.Println("server start at http://127.0.0.1:8080")
	http.ListenAndServe(":8080", nil)
}

// 1. curl -I http://127.0.0.1:8080
// HTTP/1.1 200 OK
// X-Ratelimit-Remaing: 1
// X-Ratelimit-Reset-After: 10000
// X-Ratelimit-Total: 2
// Date: Sat, 04 Jun 2022 05:04:25 GMT
// Content-Length: 12
// Content-Type: text/plain; charset=utf-8

// 1. curl -I http://127.0.0.1:8080
// HTTP/1.1 200 OK
// X-Ratelimit-Remaing: 0
// X-Ratelimit-Reset-After: 8867
// X-Ratelimit-Total: 2
// Date: Sat, 04 Jun 2022 05:08:07 GMT
// Content-Length: 12
// Content-Type: text/plain; charset=utf-8

// 3. curl -I http://127.0.0.1:8080
// HTTP/1.1 429 Too Many Requests
// Content-Type: text/plain; charset=utf-8
// X-Content-Type-Options: nosniff
// Date: Sat, 04 Jun 2022 05:03:19 GMT
// Content-Length: 20
