# JWT Redirect
JWT Redirect is a simple middleware plugin for [Traefik](https://github.com/traefik/traefik) which based on some a JWT 
field extracted from the `Authorization` header, that is configurable, and a list of values also configurable a redirect 307 will be done.

## Configuration

### Static

```toml
[pilot]
  token = "xxxx"

[experimental.plugins.jwtredirect]
  modulename = "github.com/dragosnutu/traefik-plugin"
  version = "v0.0.5"
```

### Dynamic

To configure the `JWT Redirect` plugin you should create a [middleware](https://docs.traefik.io/middlewares/overview/) in 
your dynamic configuration as explained [here](https://docs.traefik.io/middlewares/overview/). The following example creates
and uses the `jwtredirect` middleware plugin to redirect the user to a specific url based on values and a field value extracted from jwt.

```toml
[http.routers]
  [http.routers.my-router]
    rule = "Host(`localhost`)"
    middlewares = ["jwtredirect-sub"]
    service = "my-service"

[http.middlewares]
  [http.middlewares.jwtredirect-sub.plugin.jwtredirect]
    jwtField = "sub"
    jwtVaues = ["val1", "val2"]
    redirect = "https://google.com"

[http.services]
  [http.services.my-service]
    [http.services.my-service.loadBalancer]
      [[http.services.my-service.loadBalancer.servers]]
        url = "http://127.0.0.1"
```

This configuration will basically look at the `sub` field form JWT `Authorization` of `Bearer` type, and if the value is `val1` or `val2` it will 
return a `307` status and a `Location` header to the `https://google.com` 


## Release notes

### 0.0.6 - DEV
### 0.0.5 (28.09.2020)
* documentation added
* renamed the plugin
* added logs
