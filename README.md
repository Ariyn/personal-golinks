# Personal GoLinks

This is a personal purposed go links implement.

## Usage
1. run the server
    ```bash
    $ go run ./apps/server
    ```

2. add a go link
    ```bash
    $ curl localhost -XPOST -d 'example
    https://example.org'
    ok
    ```

3. check the result
    ```bash
    $ curl localhsot/example -v
    *   Trying 127.0.0.1:80...
    * Connected to localhost (127.0.0.1) port 80 (#0)
    > GET /example333 HTTP/1.1
    > Host: localhost
    > User-Agent: curl/8.1.2
    > Accept: */*
    > 
    < HTTP/1.1 303 See Other
    < Content-Type: text/html; charset=utf-8
    < Location: https://example.org
    < Date: Wed, 23 Oct 2024 12:08:36 GMT
    < Content-Length: 46
    < 
    <a href="https://example.org">See Other</a>.
    
    * Connection #0 to host localhost left intact
    ```