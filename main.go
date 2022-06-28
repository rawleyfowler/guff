package main

import(
	"fmt"
	"os"
	"strings"
	"flag"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

var HTTP_METHODS [6]string = [6]string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"}

type HttpMethod string
const (
	Get HttpMethod = "GET"
	Post           = "POST"
	Patch          = "PATCH"
	Put            = "PUT"
	Delete         = "DELETE"
	Options        = "OPTIONS"
)

type Query struct {
	Uri string
	Headers map[string]string
	Payload string
	Method HttpMethod
}

func main() {
	// Read flags
	p := flag.String("p", "", "Specify a payload")
	P := flag.String("P", "", "Specify a payload located inside of a file")
	x := flag.String("x", "", "Specify additional headers")
	M := flag.String("M", "GET", "Specify the HTTP method to use")
	l := flag.String("l", "", "Log results to a given file")
	c := flag.Int("c", 1, "The number of times to make the request")
	h := flag.Bool("h", false, "Display this message")

	flag.Parse()

	if *h {
		fmt.Println("USAGE: guff URI [-pPxMlch]")
		flag.PrintDefaults()
		fmt.Println("EXAMPLE: guff www.gnu.org -c 10")
		fmt.Println("NOTE: Additional headers should be provided in the format 'KEY:VALUE;KEY:VALUE...'")
		os.Exit(1)
	}
	
	e := []error{}
	
	args := flag.Args()
	var uri string
	switch len(args) {
	case 0:
		e = append(e, errors.New("You must specify a target URI as the first argument"))
		break
	case 1:
		if !strings.HasPrefix(args[0], "https://") || !strings.HasPrefix(args[0], "http://") {
			uri = "http://" + args[0]
		} else {
			uri = args[0]
		}
		break
	default:
		e = append(e, errors.New("Too many arguments..."))
		break
	}

	if len(*p) > 0 && len(*P) > 0 {
		e = append(e, errors.New("Cannot specify two payloads, you can either use -p or -P"))
	}

	if !IsValidHTTPMethod(HttpMethod(*M)) {
		e = append(e, errors.New("Invalid HTTP method"))
	}

	if *c < 1 {
		e = append(e, errors.New("Cannot make less than 1 request"))
	}

	// Load the payload
	if len(*P) > 0 {
		contents, err := ioutil.ReadFile(*P)
		if err != nil {
			e = append(e, errors.New("Could not open file " + *P))
		} else {
			strContents := string(contents)
			p = &strContents
		}
	}

	client := &http.Client{}
	
	q := Query{
		Uri: uri,
		Method: HttpMethod(*M),
		Payload: *p,
		Headers: make(map[string]string),
	}
	
	// Get custom headers, add them to the client
	if len(*x) > 0 {
		headerPairs := strings.Split(*l, ";")
		for _, pair := range headerPairs {
			v := strings.Split(pair, ":")
			if len(v) < 2 {
				e = append(e, errors.New("Headers are not in the correct format, 'HEADER:VALUE;HEADER:VALUE...'"))
				break
			}
			q.Headers[v[0]] = v[1]
		}
	}
	
	// Log any errors and exit 1
	if len(e) > 0 {
		fmt.Println("Failed with the following errors:")
		for _, v := range e {
			fmt.Println("ERROR: " + v.Error())
		}
		fmt.Println("Try using -h for help")
		os.Exit(1)
	}
	
	res, queryErr := MakeRequest(q, client)

	if queryErr != nil {
		fmt.Println("ERROR: failed to make connection to requested URI")
		os.Exit(1)
	}

	fmt.Println(res)
}

func IsValidHTTPMethod(m HttpMethod) bool {
	if m == Get ||
		m == Post ||
		m == Put ||
		m == Patch ||
		m == Delete ||
		m == Options {
		return true
	}
	return false
}

func MakeRequest(q Query, client *http.Client) (string, error) {
	var err error
	var res *http.Response
	var req *http.Request

	body := strings.NewReader(q.Payload)
	
	switch q.Method {
	case Get:
		req, err = http.NewRequest(string(Get), q.Uri, body)
		break
	case Post:
		req, err = http.NewRequest(string(Post), q.Uri, body)
		break
	case Put:
		req, err = http.NewRequest(string(Put), q.Uri, body)
		break
	case Patch:
		req, err = http.NewRequest(string(Patch), q.Uri, body)
		break
	}

	for key, val := range q.Headers {
		req.Header.Add(key, val)
	}
	
	res, err = client.Do(req)
	if err != nil {
		return "", err
	}
	
	dump, err := httputil.DumpResponse(res, true)

	if err != nil {
		return "", err
	}
	return string(dump), nil
}
