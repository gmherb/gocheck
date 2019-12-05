package main

import (
	"flag"
	"fmt"
	"github.com/thedevsaddam/gojsonq"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
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
func checkHTTP() {
	// Set func time tracker
	defer timeTrack(time.Now(), "checkHTTP")

	// Set default parameters (chuck)
	var defaultHost bool = false
	if *host == "api.chucknorris.io" {
		*url = "/jokes/random"
		*header = "api.chucknorris.io"
		defaultHost = true
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

	// http client request using client.Do and wrapping with Time for metrics
	startTime := time.Now()
	resp, err := client.Do(req)
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
	// Status Code Check
	log.Println("status codes `ok`:", httpOkStatusCodes)
	log.Println("status codes `warn`:", httpWarnStatusCodes)
	log.Println("status codes `error`:", httpErrorStatusCodes)
	log.Println("status code check result:", httpCheckStatusCode)
	log.Println("status code:", resp.StatusCode)
	log.Println("status:", resp.Status)

	// Response Time Check
	log.Println("response time `ok`:", okResponseTime)
	log.Println("response time `warn`:", warnResponseTime)
	log.Println("response time `error`:", errorResponseTime)
	log.Println("response time check result:", httpCheckResponseTime)
	log.Println("response time:", respTime)

	// General checkHTTP
	log.Println("target:", *host)
	log.Println("header:", *header+*url)
	if *verbose {
		log.Println("protocol:", resp.Proto)
		log.Println("content length:", resp.ContentLength)
		log.Println("request:", resp.Request)
		log.Println("header:", resp.Header)
		log.Println("Body:", string(body))
		log.Println("trailer:", resp.Trailer)
	}

	// Parse Json
	if defaultHost {
		chuck := gojsonq.New().FromString(string(body)).Find("value")
		chuckId := gojsonq.New().FromString(string(body)).Find("id")
		chuckUrl := gojsonq.New().FromString(string(body)).Find("url")
		chuckAt := gojsonq.New().FromString(string(body)).Find("updated_at")
		log.Println("chuck id:", chuckId)
		log.Println("chuck url:", chuckUrl)
		log.Println("chuck:", chuck)
		log.Println("updated at:", chuckAt)
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
func checkNET() {
	defer timeTrack(time.Now(), "checkNet")
	log.Println("check net to do includes dns, checks, portscan, connection tests, etc")

	icmpCon, err := icmp.ListenPacket("udp4", "0.0.0.0")
	if err != nil {
		log.Fatal(err)
	}
	defer icmpCon.Close()
	log.Println("ICMP Connection:*", *icmpCon)
	log.Println("ICMP Connection:", icmpCon)

	icmpLocalAddr := icmpCon.LocalAddr()
	log.Println("Local Addr:", icmpLocalAddr)

	icmpMes := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte("PING-PONG-BONG-MONG"),
		},
	}
	log.Println("ICMP Message:", icmpMes)

	icmpMesMar, err := icmpMes.Marshal(nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("ICMP Message Marshal:", icmpMesMar)

	icmpPing, err := icmpCon.WriteTo(icmpMesMar, &net.UDPAddr{IP: net.ParseIP("8.8.8.8"), Port: 0})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("ICMP Ping:", icmpPing)

	fluff := make([]byte, 1500)
	n, target, err := icmpCon.ReadFrom(fluff)
	if err != nil {
		log.Fatal(err)
	}
	icmpResponse, err := icmp.ParseMessage(1, fluff[:n])
	if err != nil {
		log.Fatal(err)
	}
	log.Println("ICMP Response:", icmpResponse)

	switch icmpResponse.Type {
	case ipv4.ICMPTypeEchoReply:
		log.Printf("response received from %v", target)
	default:
		log.Printf("%+v recieved from target", icmpResponse)
	}

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
		log.Println("force selected!")
	}

	// execute check
	switch *subCmd {
	case "http":
		checkHTTP()
	case "net":
		log.Println("Could it be, that `net` mode is not done yet.!?")
		checkNET()
	}
}
