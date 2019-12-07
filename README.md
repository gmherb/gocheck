# gocheck
Check HTTP

## gocheck defaults
	gocheck]$ go run gocheck.go
	2019/12/05 01:13:39 Running 'http' mode...
	2019/12/05 01:13:40 status codes `ok`: [200]
	2019/12/05 01:13:40 status codes `warn`: [400]
	2019/12/05 01:13:40 status codes `error`: [500]
	2019/12/05 01:13:40 status code check result: ok
	2019/12/05 01:13:40 status code: 200
	2019/12/05 01:13:40 status: 200 OK
	2019/12/05 01:13:40 response time `ok`: 1s
	2019/12/05 01:13:40 response time `warn`: 2s
	2019/12/05 01:13:40 response time `error`: 3s
	2019/12/05 01:13:40 response time check result: ok
	2019/12/05 01:13:40 response time: 855.545648ms
	2019/12/05 01:13:40 target: api.chucknorris.io
	2019/12/05 01:13:40 header: api.chucknorris.io/jokes/random
	2019/12/05 01:13:40 chuck id: 8plI-Z3IQoCmEcRV1vSYvw
	2019/12/05 01:13:40 chuck url: https://api.chucknorris.io/jokes/8plI-Z3IQoCmEcRV1vSYvw
	2019/12/05 01:13:40 chuck: when you play COD with chuck Norris the escape button doesn't work! No one can escape from Chuck Norris.
	2019/12/05 01:13:40 updated at: 2016-05-01 10:51:41.584544
	2019/12/05 01:13:40 checkHTTP took 855.862645ms

## gocheck w/ Arg's to check a redirect

    gocheck]$ go run gocheck.go -host youtube.com -httpOkStatusCodes=301,302 
    2019/12/02 21:19:36 Running 'http' mode...
    2019/12/02 21:19:37 status codes `ok`: [301 302]
    2019/12/02 21:19:37 status codes `warn`: [400]
    2019/12/02 21:19:37 status codes `error`: [500]
    2019/12/02 21:19:37 status code check result: ok
    2019/12/02 21:19:37 response time `ok`: 1s
    2019/12/02 21:19:37 response time `warn`: 2s
    2019/12/02 21:19:37 response time `error`: 3s
    2019/12/02 21:19:37 response time check result: ok
    2019/12/02 21:19:37 target: youtube.com
    2019/12/02 21:19:37 header: youtube.com
    2019/12/02 21:19:37 status: 301 Moved Permanently
    2019/12/02 21:19:37 status code: 301
    2019/12/02 21:19:37 response time: 185.398518ms
    2019/12/02 21:19:37 protocol: HTTP/2.0
    2019/12/02 21:19:37 content length: 0
    2019/12/02 21:19:37 checkHttp took 185.554812ms

## gocheck with Host Header set and Url

    gocheck]$ go run gocheck.go -host www.google.com -httpOkResponseTime=300ms -httpWarnResponseTime=500ms -header=www.example.com -url /crazycat -httpWarnStatusCodes=400,404,403
    2019/12/02 21:32:13 Running 'http' mode...
    2019/12/02 21:32:13 status codes `ok`: [200]
    2019/12/02 21:32:13 status codes `warn`: [400 404 403]
    2019/12/02 21:32:13 status codes `error`: [500]
    2019/12/02 21:32:13 status code check result: warn
    2019/12/02 21:32:13 response time `ok`: 300ms
    2019/12/02 21:32:13 response time `warn`: 500ms
    2019/12/02 21:32:13 response time `error`: 3s
    2019/12/02 21:32:13 response time check result: ok
    2019/12/02 21:32:13 target: www.google.com
    2019/12/02 21:32:13 header: www.example.com/crazycat
    2019/12/02 21:32:13 status: 404 Not Found
    2019/12/02 21:32:13 status code: 404
    2019/12/02 21:32:13 response time: 175.987496ms
    2019/12/02 21:32:13 protocol: HTTP/2.0
    2019/12/02 21:32:13 content length: 1569
    2019/12/02 21:32:13 checkHttp took 176.212958ms
