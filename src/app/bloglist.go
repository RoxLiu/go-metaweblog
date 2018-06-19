package app

import (
	"net/http"
	"strings"
	"log"
	"io/ioutil"
	"fmt"
	"time"
)

type BlogListHandler struct {
}

/*
/blog => Blog首页地址

*/
//应答请求
func (p *BlogListHandler)Process(ctx *Context, w http.ResponseWriter, r *http.Request) (error) {
	//打印收到的消息
	s := fmt.Sprintf("Time: %s, Request: %s, %s", time.Now().String(), r.Method, r.URL)
	log.Println(s)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed to read the request body: " +err.Error())
		return err
	}
	log.Println("Request Body:\n" + string(body))

	split := strings.Split(r.URL.Path[1:], "/")

	//url: /blog
	if len(split) == 1 {
		return handleIndex(ctx, w, r)
	}
	if split[1] == "" {
		// Allow a trailing / on requests
		//return handleClusterList(ctx, w, r)
	}

/*	switch split[3] {
	case "consumer":
		switch {
		case r.Method == "DELETE":
			switch {
			case (len(split) == 5) || (split[5] == ""):
				return handleConsumerDrop(ctx, w, r, split[2], split[4])
			default:
				return makeErrorResponse(http.StatusMethodNotAllowed, "request method not supported", w, r)
			}
		case r.Method == "GET":
			switch {
			case (len(split) == 4) || (split[4] == ""):
				return handleConsumerList(ctx, w, r, split[2])
			case (len(split) == 5) || (split[5] == ""):
				// Consumer detail - list of consumer streams/hosts? Can be config info later
				return makeErrorResponse(http.StatusNotFound, "unknown API call", w, r)
			case split[5] == "topic":
				switch {
				case (len(split) == 6) || (split[6] == ""):
					return handleConsumerTopicList(ctx, w, r, split[2], split[4])
				case (len(split) == 7) || (split[7] == ""):
					return handleConsumerTopicDetail(ctx, w, r, split[2], split[4], split[6])
				}
			case split[5] == "status":
				return handleConsumerStatus(ctx, w, r, split[2], split[4], false)
			case split[5] == "lag":
				return handleConsumerStatus(ctx, w, r, split[2], split[4], true)
			}
		default:
			return makeErrorResponse(http.StatusMethodNotAllowed, "request method not supported", w, r)
		}
	case "topic":
		switch {
		case r.Method != "GET":
			return makeErrorResponse(http.StatusMethodNotAllowed, "request method not supported", w, r)
		case (len(split) == 4) || (split[4] == ""):
			return handleBrokerTopicList(ctx, w, r, split[2])
		case (len(split) == 5) || (split[5] == ""):
			return handleBrokerTopicDetail(ctx, w, r, split[2], split[4])
		}
	case "offsets":
		// Reserving this endpoint to implement later
		return makeErrorResponse(http.StatusNotFound, "unknown API call", w, r)
	}
*/

	return nil
}

func handleIndex(ctx *Context, w http.ResponseWriter, r *http.Request) (error) {
	s := `<html>
	<body>
		<center>GitPage Index File</center>
	</body>
</html>`
	log.Println("Sent The Response:\n" + s)

	w.Write([]byte(s))
	return nil
}