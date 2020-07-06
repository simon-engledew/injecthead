### injecthead

Inject CSRF tokens into static single page webapps.

Useful if you have server side rendered a React application and want to access a CSRF token without making an additional remote request.

Can be run as a `httputil.SingleHostReverseProxy` in front of a webserver like Nginx, or as a `http.Handler` wrapping another `http.Handler`, e.g: `http.FileServer`.

```golang
func csrfMeta(r *http.Request) (string, error) {
    token := csrf.Token(r)
    
    if token == "" {
        return "", errors.New("no csrf middleware installed")
    }
    
    return fmt.Sprintf(`<meta name="csrf" content="%s" />`, template.EscapeHTMLString(token)), nil
}

fs := injecthead.Handle(http.FileServer(http.Dir(".")), meta)

proxy := injecthead.Handle(httputil.NewSingleHostReverseProxy(&url.URL{
    Scheme: "http",
    Host:   "127.0.0.1:4000",
}), meta)

protect := csrf.Protect(
    CSRF_SECRET,
    csrf.Secure(true),
    csrf.FieldName("csrf_token"),
    csrf.RequestHeader("X-CSRF-TOKEN"),
    csrf.Path("/"),
    csrf.CookieName("csrf-token"),
    csrf.MaxAge(0),
)

mux := http.NewServeMux()
mux.Handle("/static", http.StripPrefix("/static", fs))
mux.Handle("/hmr", http.StripPrefix("/hmr", proxy))

http.ListenAndServe("127.0.0.1:3000", protect(mux))
```

The JavaScript application can then access the CSRF token with something like:

```javascript
function csrf() {
  const element = document.querySelector('meta[name="csrf"]');
  return element && element.getAttribute('content');
}
```
