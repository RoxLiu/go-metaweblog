package xmlrpc

import (
	"testing"
	"../../src/xmlrpc"
	"bufio"
	"os"
	"time"
)

type BlogInfo struct {
	BlogId 			string 			`xml:"blogid"`
	URL				string			`xml:"url"`
	BlogName		string			`xml:"blogName"`
}

func TestEncodeRequest(t *testing.T)  {
	v := BlogInfo{BlogId: "blog-001", URL: "http://localhost:8080", BlogName: "Test Blog"}

	s := xmlrpc.EncodeRequest("GetBlogList", v)
	t.Log(s)
}

func TestDecodeResponse(t *testing.T)  {
	fi, err := os.Open("../xml/TestDecodeResponse.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer fi.Close()
	r := bufio.NewReader(fi)

	v, err := xmlrpc.DecodeResponse(r)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(v)
}

func TestDecodeRequest(t *testing.T)  {
	fi, err := os.Open("../xml/TestDecodeRequest.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer fi.Close()
	r := bufio.NewReader(fi)

	v, err := xmlrpc.DecodeRequest(r)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(v)
}


func TestNewPost(t *testing.T)  {
	fi, err := os.Open("../xml/newPost.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer fi.Close()
	r := bufio.NewReader(fi)

	rsp, err := xmlrpc.DecodeRequest(r)
	if err != nil {
		t.Fatal(err)
	}

	blogId := rsp.GetParameter(0).(string)
	t.Log(blogId)
	//username := r.GetParameter(1)
	//password := r.GetParameter(2)
	post := rsp.GetParameter(3).(xmlrpc.Struct)

	title := post["title"].(string)
	t.Log(title)
	content := post["description"].(string)
	t.Log(content)

	publish := rsp.GetParameter(4).(bool)
	t.Log(publish)

	dateCreated := post["dateCreated"].(time.Time)
	t.Log(dateCreated)

	t.Log(time.Now().UTC().Format("2006-01-02T15:04:05Z07:00"))
}

func TestTimeParse(t *testing.T) {
	s := "2019-01-04T02:39:49Z"
	_, e := time.Parse("2006-01-02T15:04:05Z07:00", s)
	if e != nil {
		t.Fail()
	} else {
		t.Log("Parse Successfully.")
	}

	s = "20181214T06:16:34Z"
	_, e = time.Parse("20060102T15:04:05Z07:00", s)
	if e != nil {
		t.Fail()
	} else {
		t.Log("Parse Successfully.")
	}
}
/*
func makeRequest(name string, args ...interface{}) *bytes.Buffer {
	buf := new(bytes.Buffer)
	buf.WriteString(`<?xml version="1.0"?><methodCall>`)
	buf.WriteString("<methodName>" + xmlEscape(name) + "</methodName>")
	buf.WriteString("<params>")
	for _, arg := range args {
		buf.WriteString("<param><value>")
		buf.WriteString(toXml(arg, true))
		buf.WriteString("</value></param>")
	}
	buf.WriteString("</params></methodCall>")
	return buf
}

//发送请求
func call(client *http.Client, url, name string, args ...interface{}) (v interface{}, e error) {
	r, e := client.Post(url, "text/xml", makeRequest(name, args...))
	if e != nil {
		return nil, e
	}

	// Since we do not always read the entire body, discard the rest, which
	// allows the http transport to reuse the connection.
	defer io.Copy(ioutil.Discard, r.Body)
	defer r.Body.Close()

	if r.StatusCode/100 != 2 {
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	p := xml.NewDecoder(r.Body)
	se, e := nextStart(p) // methodResponse
	if se.Name.Local != "methodResponse" {
		return nil, errors.New("invalid response: missing methodResponse")
	}
	se, e = nextStart(p) // params
	if se.Name.Local != "params" {
		return nil, errors.New("invalid response: missing params")
	}
	se, e = nextStart(p) // param
	if se.Name.Local != "param" {
		return nil, errors.New("invalid response: missing param")
	}
	se, e = nextStart(p) // value
	if se.Name.Local != "value" {
		return nil, errors.New("invalid response: missing value")
	}
	_, v, e = next(p)
	return v, e
}


// Client is client of XMLRPC
type Client struct {
	HttpClient *http.Client
	url        string
}

// NewClient create new Client
func NewClient(url string) *Client {
	return &Client{
		HttpClient: &http.Client{Transport: http.DefaultTransport, Timeout: 10 * time.Second},
		url:        url,
	}
}

// Call call remote procedures function name with args
func (c *Client) Call(name string, args ...interface{}) (v interface{}, e error) {
	return call(c.HttpClient, c.url, name, args...)
}

// Global httpClient allows us to pool/reuse connections and not wastefully
// re-create transports for each request.
var httpClient = &http.Client{Transport: http.DefaultTransport, Timeout: 10 * time.Second}

// Call call remote procedures function name with args
func Call(url, name string, args ...interface{}) (v interface{}, e error) {
	return call(httpClient, url, name, args...)
}

*/



/*

func createServer(path, name string, f func(args ...interface{}) (interface{}, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != path {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		p := xml.NewDecoder(r.Body)
		se, _ := nextStart(p) // methodResponse
		if se.Name.Local != "methodCall" {
			http.Error(w, "missing methodCall", http.StatusBadRequest)
			return
		}
		se, _ = nextStart(p) // params
		if se.Name.Local != "methodName" {
			http.Error(w, "missing methodName", http.StatusBadRequest)
			return
		}
		var s string
		if err := p.DecodeElement(&s, &se); err != nil {
			http.Error(w, "wrong function name", http.StatusBadRequest)
			return
		}
		if s != name {
			http.Error(w, fmt.Sprintf("want function name %q but got %q", name, s), http.StatusBadRequest)
			return
		}
		se, _ = nextStart(p) // params
		if se.Name.Local != "params" {
			http.Error(w, "missing params", http.StatusBadRequest)
			return
		}
		var args []interface{}
		for {
			se, _ = nextStart(p) // param
			if se.Name.Local == "" {
				break
			}
			if se.Name.Local != "param" {
				http.Error(w, "missing param", http.StatusBadRequest)
				return
			}
			se, _ = nextStart(p) // value
			if se.Name.Local != "value" {
				http.Error(w, "missing value", http.StatusBadRequest)
				return
			}
			_, v, err := next(p)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			args = append(args, v)
		}

		ret, err := f(args...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write([]byte(`
		<?xml version="1.0"?>
		<methodResponse>
		<params>
			<param>
				<value>` + toXml(ret, true) + `</value>
			</param>
		</params>
		</methodResponse>
		`))
	}
}

func TestAddInt(t *testing.T) {
	ts := httptest.NewServer(createServer("/api", "AddInt", func(args ...interface{}) (interface{}, error) {
		if len(args) != 2 {
			return nil, errors.New("bad number of arguments")
		}
		switch args[0].(type) {
		case int:
		default:
			return nil, errors.New("args[0] should be int")
		}
		switch args[1].(type) {
		case int:
		default:
			return nil, errors.New("args[1] should be int")
		}
		return args[0].(int) + args[1].(int), nil
	}))
	defer ts.Close()

	client := NewClient(ts.URL + "/api")
	v, err := client.Call("AddInt", 1, 2)
	if err != nil {
		t.Fatal(err)
	}
	i, ok := v.(int)
	if !ok {
		t.Fatalf("want int but got %T: %v", v, v)
	}
	if i != 3 {
		t.Fatalf("want %v but got %v", 3, v)
	}
}

func TestAddString(t *testing.T) {
	ts := httptest.NewServer(createServer("/api", "AddString", func(args ...interface{}) (interface{}, error) {
		if len(args) != 2 {
			return nil, errors.New("bad number of arguments")
		}
		switch args[0].(type) {
		case string:
		default:
			return nil, errors.New("args[0] should be string")
		}
		switch args[1].(type) {
		case string:
		default:
			return nil, errors.New("args[1] should be string")
		}
		return args[0].(string) + args[1].(string), nil
	}))
	defer ts.Close()

	client := NewClient(ts.URL + "/api")
	v, err := client.Call("AddString", "hello", "world")
	if err != nil {
		t.Fatal(err)
	}
	s, ok := v.(string)
	if !ok {
		t.Fatalf("want string but got %T: %v", v, v)
	}
	if s != "helloworld" {
		t.Fatalf("want %q but got %q", "helloworld", v)
	}
}
*/