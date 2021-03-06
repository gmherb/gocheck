package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"github.com/thedevsaddam/gojsonq"
)

// main flags
var subCmd = flag.String("cmd", "http", "http|tcp|icmp")
var host = flag.String("host", "api.chucknorris.io", "target host")
var verbose = flag.Bool("verbose", false, "enable verbose mode")

// http mode flags
var header = flag.String("http-header", "", "target host header")
var url = flag.String("http-url", "", "target host url")
var follow = flag.Bool("http-follow", false, "follow http redirects")
var httpWarnResponseTime = flag.String("http-response-time-warn", "1s", "warn time")
var httpErrorResponseTime = flag.String("http-response-time-error", "3s", "error time")

// icmp mode flags
var icmpCount = flag.Int("icmp-count", 1, "amount of icmp echos")
var icmpTimeout = flag.String("icmp-timeout", "3s", "timeout for icmp check")
var icmpWarnResponseTime = flag.String("icmp-response-time-warn", "1s", "warn time")
var icmpErrorResponseTime = flag.String("icmp-response-time-error", "3s", "error time")

// tcp mode flags
var tcpPort = flag.String("tcp-port", "80", "port to connect to in tcp check")
var tcpTimeout = flag.String("tcp-timeout", "3s", "timeout for tcp check")
var tcpWarnResponseTime = flag.String("tcp-response-time-warn", "1s", "warn time")
var tcpErrorResponseTime = flag.String("tcp-response-time-error", "3s", "error time")

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


	//*********** Start Status Code Check ***********//
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
			log.Println("specific http error code was not provided, alerting on all non-ok/warn codes")
		}
	}

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
	} else if len(errorStatusCodeString) == 0 {
		httpCheckStatusCode = "unknown"
	} else {
		httpCheckStatusCode = "unknown"
	}




	//*********** Start Response Time Check ***********//
	// Set responseTime in Duration type to be compared
	warnResponseTime, _ := time.ParseDuration(*httpWarnResponseTime)
	errorResponseTime, _ := time.ParseDuration(*httpErrorResponseTime)

	// Perform response time comparison
	var httpCheckResponseTime string
	if respTime <= warnResponseTime {
		httpCheckResponseTime = "ok"
	} else if respTime <= errorResponseTime {
		httpCheckResponseTime = "warn"
	} else {
		httpCheckResponseTime = "err"
	}




	//*********** Logging Below ***********//
	// Status Code Check
	log.Println("status codes `ok`:", httpOkStatusCodes)
	log.Println("status codes `warn`:", httpWarnStatusCodes)
	log.Println("status codes `error`:", httpErrorStatusCodes)
	log.Println("status code check result:", httpCheckStatusCode)
	log.Println("status code:", resp.StatusCode)
	log.Println("status:", resp.Status)

	// Response Time Check
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

// Check NET function!
func checkICMP(h net.IP) {
	defer timeTrack(time.Now(), "checkICMP")

	// Set icmp Listener
	icmpCon, err := icmp.ListenPacket("udp4", "0.0.0.0")
	if err != nil {
		log.Fatal(err)
	}
	defer icmpCon.Close()
	if *verbose {
		log.Println("ICMP Connection:*", *icmpCon)
		log.Println("ICMP Connection:", icmpCon)
	}

	icmpLocalAddr := icmpCon.LocalAddr()
	if *verbose {
		log.Println("Local Addr:", icmpLocalAddr)
	}

	// Prep ICMP Message
	icmpMes := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte("PING-PONG-BONG-MONG"),
		},
	}
	if *verbose {
		log.Println("ICMP Message:", icmpMes)
	}

	icmpMesMar, err := icmpMes.Marshal(nil)
	if err != nil {
		log.Fatal(err)
	}
	if *verbose {
		log.Println("ICMP Message Marshal:", icmpMesMar)
	}

	// make slice with make for icmp results
	icmpResults := make([]time.Duration, *icmpCount)

	// loop icmp
	i := 0
	for i < *icmpCount {

		icmpPing, err := icmpCon.WriteTo(icmpMesMar, &net.UDPAddr{IP: net.IP(h), Port: 0})
		if err != nil {
			log.Fatal(err)
		}
		if *verbose {
			log.Println("ICMP Ping:", icmpPing)
		}

		icmpTimeout, err := time.ParseDuration(*icmpTimeout)
		if err != nil {
			log.Fatal(err)
		}
		if *verbose {
			log.Println("ICMP Timeout:", icmpTimeout)
		}

		// Set time in future to timeout
		icmpConTimeout := icmpCon.SetDeadline(time.Now().Add(icmpTimeout))
		if *verbose {
			log.Println("ICMP Conn Timeout:", icmpConTimeout)
			log.Println("ICMP Warn Timeout:", *icmpWarnResponseTime)
			log.Println("ICMP Error Timeout:", *icmpErrorResponseTime)
		}

		// make request and read from target, wrapping with Time for metrics
		startTime := time.Now()

		fluff := make([]byte, 1500)
		n, target, err := icmpCon.ReadFrom(fluff)
		if err != nil {
			log.Fatal(err)
		}

		icmpResponse, err := icmp.ParseMessage(1, fluff[:n])
		if err != nil {
			log.Fatal(err)
		}
		if *verbose {
			log.Println("ICMP Response:", icmpResponse)
		}

		// Add time result to slice
		icmpTime := time.Since(startTime)
		icmpResults[i] = icmpTime


		// TODO - ADD response time check here


		// Result of icmp
		switch icmpResponse.Type {
		case ipv4.ICMPTypeEchoReply:
			if *verbose {
				log.Println("ICMP response received from", target, "was", icmpResponse.Type)
			}
		default:
			if *verbose {
				log.Println("ICMP response received from", target, "was", icmpResponse.Type)
			}
		}

		// Add count to loop var
		i = i + 1
	}

	// Log icmp results
	log.Println("ICMP Results:", icmpResults)

	// Close icmp listener <- this is not correct
	icmpCon.Close()

}

// Check TCP
func checkTCP(h string) {
	defer timeTrack(time.Now(), "checkTCP")

	// Connection Port
	connP := *tcpPort
	log.Println("TCP Port:", connP)

	// Connection String
	connS := net.JoinHostPort(h, *tcpPort)
	log.Println("TCP Connection String:", connS)

	// Connection Timeout
	connT, err := time.ParseDuration(*tcpTimeout)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("TCP Timeout:", connT)


	// TCP Connection using net.DialTimeout and wrapping with Time for metrics
	tcpStartTime := time.Now()
	conn, err := net.DialTimeout("tcp", connS, connT)
	if err != nil {
		log.Fatal(err) // handle error
	}

	tcpEndTime := time.Since(tcpStartTime)
	if *verbose {
		log.Println("TCP Connection:", conn)
	}


	//*********** Start Response Time Check ***********//
	// Set responseTime in Duration type to be compared
	warnResponseTime, _ := time.ParseDuration(*tcpWarnResponseTime)
	errorResponseTime, _ := time.ParseDuration(*tcpErrorResponseTime)

	// Perform response time comparison
	var tcpCheckResponseTime string
	if tcpEndTime <= warnResponseTime {
		tcpCheckResponseTime = "ok"
	} else if tcpEndTime <= errorResponseTime {
		tcpCheckResponseTime = "warn"
	} else {
		tcpCheckResponseTime = "err"
	}


	//*********** Logging Below ***********//
	// Response Time Check
	log.Println("response time `warn`:", warnResponseTime)
	log.Println("response time `error`:", errorResponseTime)
	log.Println("response time check result:", tcpCheckResponseTime)
	log.Println("response time:", tcpEndTime)
}

func main() {
	// Parse variables that require new type to be defined
	// Note to self: was unable to define directly into var, thus cannot be outside func
	flag.Var(&httpOkStatusCodes, "http-ok-status-codes", "http | status codes to `ok` on")
	flag.Var(&httpWarnStatusCodes, "http-warn-status-codes", "http | status codes to `warn` on")
	flag.Var(&httpErrorStatusCodes, "http-error-status-codes", "http | status codes to `error` on")

	flag.Parse()

	// showing user selected options
	log.Println("Running '" + *subCmd + "' mode...")

	// DNS Resolve
	hostLookup, err := net.LookupHost(*host)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Host Lookup:", hostLookup)

	// Grab first in DNS resolve results
	firstHost := net.ParseIP(hostLookup[0])
	firstHostStr := hostLookup[0]
	log.Println("Host:", firstHost)

	// execute check
	switch *subCmd {
	case "http":
		checkHTTP()
	case "icmp":
		checkICMP(firstHost)
	case "tcp":
		checkTCP(firstHostStr)
	}
}
