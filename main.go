package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/TykTechnologies/tyk/ctx"
	"github.com/TykTechnologies/tyk/headers"
	"github.com/TykTechnologies/tyk/log"
	"github.com/TykTechnologies/tyk/user"
)

// called once plugin is loaded, this is where we put all initialization work for plugin
// i.e. setting exported functions, setting up connection pool to storage and etc.
func init() {
	var logger = log.Get()
	logger.Info("Processing Golang plugin init function!!" )
	//Here you write the code for db connection
}


func ResponseSendCurrentTime(rw http.ResponseWriter, r *http.Request) {

	var logger = log.Get()
	apidef := ctx.GetDefinition(r)
	fmt.Println("Golang plugin - fmt example - API name is ", apidef.Name)

	logger.WithField("api-name", apidef.Name).Info("Processing HTTP request in Golang plugin!!" )

	//Demo injecting header to a request
	logger.WithField("api-name", apidef.Name).Info("Golang plugin - Adding header to a request before it goes upstream.")
	r.Header.Add("Foo", "Bar")

	logger.WithField("api-name", apidef.Name).Info("Golang plugin - ResponseSendCurrentTime")

	now := time.Now().String()


	getTime := r.URL.Query().Get("get_time")
	logger.WithField("api-name", apidef.Name).Info("Golang plugin - get_time is ", getTime)

	// check if we don't need to send reply
	if getTime != "2" {
		// allow request to be processed and sent to upstream
		logger.WithField("api-name", apidef.Name).Info("Golang plugin - Adding current_time as a header in the request. Request to api will continue to the upstream")
		r.Header.Add("current_time", now)
		return
	}

	// send HTTP response from Golang plugin
	logger.WithField("api-name", apidef.Name).Info("Golang plugin - Setting the response header and body. Request will stop in this plugin and OK response will be returned")

	// prepare data to send
	replyData := map[string]string{
		"current_time": now,
	}

	//jsonData, err := json.Marshal(replyData)
	//if err != nil {
	//	rw.WriteHeader(http.StatusInternalServerError)
	//	return
	//}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	writeBody(rw, replyData)
	//rw.Write(jsonData)
}

func writeBody(rw http.ResponseWriter, replyJson map[string]string) error{

	var logger = log.Get()
	jsonData, err := json.Marshal(replyJson)
	if err != nil {
		logger.WithField("map-response", replyJson).Info("Golang auth plugin: Failed to marshal map")
		rw.WriteHeader(http.StatusInternalServerError)
		return err
	}

	rw.Write(jsonData)
	return nil
}
// ------------------------------------------------------------------
// Custom auth plugin code:

func getSessionByKey(key string) *user.SessionState {
	// here goes our logic to check if passed API key is valid and appropriate key session can be retrieved

	// perform auth (only one token "abc" is allowed)
	// Here you add code to query your database
	if key != "abc" {
		return nil
	}

	// return session
	return &user.SessionState{
		OrgID: "default",
		Alias: "abc-session",
	}
}

func MyPluginCustomAuthCheck(rw http.ResponseWriter, r *http.Request) {

	var logger = log.Get()
	apidef := ctx.GetDefinition(r)
	logger.WithField("api-name", apidef.Name).Info("Golang auth plugin - MyPluginCustomAuthCheck")

	// try to get session by API key
	key := r.Header.Get(headers.Authorization)
	session := getSessionByKey(key)
	if session == nil {
		// auth failed, reply with 403
		logger.WithField("api-name", apidef.Name).Info("Golang auth plugin - MyPluginCustomAuthCheck - failed")
		rw.WriteHeader(http.StatusForbidden)

		// prepare data to send
		replyData := map[string]string{
			"reason": "Access forbidden",
		}

		writeBody(rw, replyData)
		//jsonData, err := json.Marshal(replyData)
		//if err != nil {
		//	rw.WriteHeader(http.StatusInternalServerError)
		//
		//	rw.Write(jsonData)
		//	return
		//}

		return
	}

	logger.WithField("api-name", apidef.Name).Info("Golang auth plugin - MyPluginCustomAuthCheck - succeeded")


	// auth was successful, add session and key to request's context so other middle-wares can use it
	ctx.SetSession(r, session, key, true)
}

func main() {}
