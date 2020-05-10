# gocheck
![Go](https://github.com/gmherb/gocheck/workflows/Go/badge.svg?branch=master)

## gocheck -help
    -cmd string
          http|tcp|icmp (default "http")
    -host string
          target host (default "api.chucknorris.io")
    -http-error-status-codes error
          http | status codes to error on
    -http-follow
          follow http redirects
    -http-header string
          target host header
    -http-ok-status-codes ok
          http | status codes to ok on
    -http-response-time-error string
          error time (default "3s")
    -http-response-time-warn string
          warn time (default "1s")
    -http-url string
          target host url
    -http-warn-status-codes warn
          http | status codes to warn on
    -icmp-count int
          amount of icmp echos (default 1)
    -icmp-response-time-error string
          error time (default "3s")
    -icmp-response-time-warn string
          warn time (default "1s")
    -icmp-timeout string
          timeout for icmp check (default "3s")
    -tcp-port string
          port to connect to in tcp check (default "80")
    -tcp-response-time-error string
          error time (default "3s")
    -tcp-response-time-warn string
          warn time (default "1s")
    -tcp-timeout string
          timeout for tcp check (default "3s")
    -verbose
          enable verbose mode

## gocheck -cmd http
    020/05/10 11:33:42 Running 'http' mode...
    2020/05/10 11:33:42 Host Lookup: [104.28.13.58 104.28.12.58 2606:4700:3035::681c:c3a 2606:4700:3035::681c:d3a]
    2020/05/10 11:33:42 Host: 104.28.13.58
    2020/05/10 11:33:43 status codes `ok`: [200]
    2020/05/10 11:33:43 status codes `warn`: [400]
    2020/05/10 11:33:43 status codes `error`: []
    2020/05/10 11:33:43 status code check result: ok
    2020/05/10 11:33:43 status code: 200
    2020/05/10 11:33:43 status: 200 OK
    2020/05/10 11:33:43 response time `warn`: 1s
    2020/05/10 11:33:43 response time `error`: 3s
    2020/05/10 11:33:43 response time check result: ok
    2020/05/10 11:33:43 response time: 890.200108ms
    2020/05/10 11:33:43 target: api.chucknorris.io
    2020/05/10 11:33:43 header: api.chucknorris.io/jokes/random
    2020/05/10 11:33:43 chuck id: qaetpus9twgdkckd70jisw
    2020/05/10 11:33:43 chuck url: https://api.chucknorris.io/jokes/qaetpus9twgdkckd70jisw
    2020/05/10 11:33:43 chuck: The movie "Delta Force" was extremely hard to make because Chuck had to downplay his abilities. The first few cuts were completely unbelievable.
    2020/05/10 11:33:43 updated at: 2020-01-05 13:42:19.324003
    2020/05/10 11:33:43 checkHTTP took 890.408309ms
    
## gocheck -cmd tcp
    2020/05/10 11:41:15 Running 'tcp' mode...
    2020/05/10 11:41:15 Host Lookup: [104.28.13.58 104.28.12.58 2606:4700:3035::681c:c3a 2606:4700:3035::681c:d3a]
    2020/05/10 11:41:15 Host: 104.28.13.58
    2020/05/10 11:41:15 TCP Port: 80
    2020/05/10 11:41:15 TCP Connection String: 104.28.13.58:80
    2020/05/10 11:41:15 TCP Timeout: 3s
    2020/05/10 11:41:15 response time `warn`: 1s
    2020/05/10 11:41:15 response time `error`: 3s
    2020/05/10 11:41:15 response time check result: ok
    2020/05/10 11:41:15 response time: 14.915503ms
    2020/05/10 11:41:15 checkTCP took 14.955198ms
    
## gocheck -cmd icmp
    [manj gocheck]# go run gocheck.go -cmd icmp
    2020/05/10 11:34:26 Running 'icmp' mode...
    2020/05/10 11:34:26 Host Lookup: [104.28.12.58 104.28.13.58 2606:4700:3035::681c:d3a 2606:4700:3035::681c:c3a]
    2020/05/10 11:34:26 Host: 104.28.12.58
    2020/05/10 11:34:27 ICMP Results: [9.516106ms]
    2020/05/10 11:34:27 checkICMP took 9.605263ms
