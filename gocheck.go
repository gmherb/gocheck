package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	//"net"
	"net/http"
	"strconv"
	"strings"
	"time"
	//"encoding/json"
)

/*
*  `-http` is selected by default
*  $ gocheck [-http|-net] [options]
*
*  checking redirect works,
*  $ gocheck --host follow youtube.com --httpOkStatusCode 301 --httpOkStatusCode 302
 */

// main flags
var subCmd = flag.String("subCmd", "http", "the subcommand")
var host = flag.String("host", "api.chucknorris.io", "target host")

// optional flags
var force = flag.Bool("force", false, "enable to overide prompts")
var verbose = flag.Bool("verbose", false, "enable verbose mode")

// http mode flags
var header = flag.String("header", "", "target host header")
var url = flag.String("url", "", "target host url")
var follow = flag.Bool("follow", false, "follow redirects in http check")
var httpOkResponseTime = flag.String("httpOkResponseTime", "1s", "ok time")
var httpWarnResponseTime = flag.String("httpWarnResponseTime", "2s", "warn time")
var httpErrorResponseTime = flag.String("httpErrorResponseTime", "3s", "error time")

// new type for httpStatusCodes, a slice of strings
// to be used by OK, WARN, ERROR statements
type httpStatusCodes []string

// two flag.Value interface methods required:
// (1) flag String() string
func (s *httpStatusCodes) String() string {
	return fmt.Sprint(*s)
}

// (2) flag Set(value string) error
func (s *httpStatusCodes) Set(value string) error {
	// if httpStatusCodes given contains a comma, split value
	ss := strings.Split(value, ",")
	*s = append(*s, ss...) // "..." as suffice in order to append slice into another
	return nil
}

// variables from new type
var httpOkStatusCodes httpStatusCodes
var httpWarnStatusCodes httpStatusCodes
var httpErrorStatusCodes httpStatusCodes

// Wrapper for time tracking
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////

/////    /////   ///////////////////   ///////////////////   /////////////
/////    /////   ///////////////////   ///////////////////   /////////////
/////    /////          /////                 /////          /////    ////
/////    /////          /////                 /////          /////    ////
//////////////          /////                 /////          /////    ////
//////////////          /////                 /////          /////////////
/////    /////          /////                 /////          /////////////
/////    /////          /////                 /////          /////
/////    /////          /////                 /////          /////
/////    /////          /////                 /////          /////

/////  ///  ///  /////  /////  //  //
//     ///  ///  //     //     // //
//     ////////  ///    //     ////
//     ///  ///  //     //     // //
/////  ///  ///  /////  /////  //  //   //   //  //   //     //

//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////

// Check HTTP
func checkHttp() {
	// Set func time tracker
	defer timeTrack(time.Now(), "checkHttp")

	// Set default parameters (chuck)
	if *host == "api.chucknorris.io" {
		*url = "/jokes/random"
		*header = "api.chucknorris.io"
	}

	// Setting client options
	tran := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{
		Transport: tran,
		Timeout:   10 * time.Second,
		// Uncomment to not follow redirects - TODO -make bool option
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// If no header was passed, set to host
	if *header == "" {
		*header = *host
	}

	// Setting host header to a request
	req, err := http.NewRequest("GET", "https://"+*host+*url, nil)
	req.Header.Set("Host", *header)

	// Start Timer for Http Request
	startTime := time.Now()

	// http client request using golang client.Do
	resp, err := client.Do(req)

	// End Timer for Http Request
	respTime := time.Since(startTime)

	// http error handling needs work
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	// loading body in ioreader
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// Start Status Code Check //

	// Check if httpStatusCodes were provided, otherwise set default.
	if httpOkStatusCodes == nil {
		if *verbose {
			log.Println("no httpOkStatusCode value was provided; set default to `200`")
		}
		httpOkStatusCodes = append(httpOkStatusCodes, "200")
	}
	if httpWarnStatusCodes == nil {
		if *verbose {
			log.Println("no httpWarnStatusCode value was provided; set default to `400`")
		}
		httpWarnStatusCodes = append(httpWarnStatusCodes, "400")
	}
	if httpErrorStatusCodes == nil {
		if *verbose {
			log.Println("no httpErrorStatusCode value was provided; set default to `500`")
		}
		httpErrorStatusCodes = append(httpErrorStatusCodes, "500")
	}

	// End Status Code Check //

	// Start Response Time Check //

	// Convert response statusCode from int to string
	statusCodeString := strconv.Itoa(resp.StatusCode)

	// Convert type httpStatusCode values to String's
	okStatusCodeString := fmt.Sprint(httpOkStatusCodes)
	warnStatusCodeString := fmt.Sprint(httpWarnStatusCodes)
	errorStatusCodeString := fmt.Sprint(httpErrorStatusCodes)

	// Perform status code comparison
	var httpCheckStatusCode string
	if strings.Contains(okStatusCodeString, statusCodeString) {
		httpCheckStatusCode = "ok"
	} else if strings.Contains(warnStatusCodeString, statusCodeString) {
		httpCheckStatusCode = "warn"
	} else if strings.Contains(errorStatusCodeString, statusCodeString) {
		httpCheckStatusCode = "err"
	} else {
		httpCheckStatusCode = "unknown"
	}

	// Set responseTime in Duration type to be compared
	okResponseTime, _ := time.ParseDuration(*httpOkResponseTime)
	warnResponseTime, _ := time.ParseDuration(*httpWarnResponseTime)
	errorResponseTime, _ := time.ParseDuration(*httpErrorResponseTime)

	// Perform response time comparison
	var httpCheckResponseTime string
	if respTime <= okResponseTime {
		httpCheckResponseTime = "ok"
	} else if respTime <= warnResponseTime {
		httpCheckResponseTime = "warn"
	} else if respTime <= errorResponseTime {
		httpCheckResponseTime = "err"
	} else {
		httpCheckResponseTime = "unknown"
	}

	// End Response Time Check //

	// Logging Below //

	// Status Code Check Out
	log.Println("status codes `ok`:", httpOkStatusCodes)
	log.Println("status codes `warn`:", httpWarnStatusCodes)
	log.Println("status codes `error`:", httpErrorStatusCodes)
	log.Println("status code check result:", httpCheckStatusCode)
	log.Println("status code:", resp.StatusCode)
	log.Println("status:", resp.Status)

	// Response Time Check Out
	log.Println("response time `ok`:", okResponseTime)
	log.Println("response time `warn`:", warnResponseTime)
	log.Println("response time `error`:", errorResponseTime)
	log.Println("response time check result:", httpCheckResponseTime)
	log.Println("response time:", respTime)

	// Print checkHttp ouput
	log.Println("target:", *host)
	log.Println("header:", *header+*url)
	log.Println("protocol:", resp.Proto)
	log.Println("content length:", resp.ContentLength)
	if *verbose {
		log.Println("request:", resp.Request)
		log.Println("header:", resp.Header)
		log.Println("Body:", string(body))
		log.Println("trailer:", resp.Trailer)
	}
}

//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////

/////     /////   ////////////////  //////////////////
//////    /////   ///////////////   ////////////////
///////   /////    /////                 /////
////////  /////     /////                /////
///////// /////    //////////            /////
////  /////////     ///////              /////
////   ////////    /////                 /////
////    ///////    ////                  /////
////     //////    //////////////        /////
////      /////    //////////////        /////

/////  ///  ///  /////  /////  //  //
//     ///  ///  //     //     // //
//     ////////  ///    //     ////
//     ///  ///  //     //     // //
/////  ///  ///  /////  /////  //  //   //   //  //   //     //

//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////

// Check Network
func checkNet() {
	defer timeTrack(time.Now(), "checkNet")
	log.Println("check net to do:\n- dns\n- icmp\n- nmap scan\n")

}

func main() {
	// Parse variables that require new type to be defined
	// Note to self: was unable to define directly into var, thus cannot be outside func
	flag.Var(&httpOkStatusCodes, "httpOkStatusCodes", "http | status codes to `ok` on")
	flag.Var(&httpWarnStatusCodes, "httpWarnStatusCodes", "http | status codes to `warn` on")
	flag.Var(&httpErrorStatusCodes, "httpErrorStatusCodes", "http | status codes to `error` on")

	flag.Parse()

	// showing user selected options
	log.Println("Running '" + *subCmd + "' mode...")

	// showing force was selected
	if *force {
		log.Println("\nforce selected!\n")
	}

	// execute check
	switch *subCmd {
	case "http":
		checkHttp()
	case "net":
		log.Println("Could it be, that `net` mode is not done yet.!?")
		checkNet()
	}
}
