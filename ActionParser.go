package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
)

var requiredHeadersForMode map[uint8][]string = map[uint8][]string{
	uint8(3): []string{"x-github-event"},
}

var eventHeader map[uint8]string = map[uint8]string{
	3: "x-github-event",
}

//Aliases for variablenames
var payloadAlias map[string]string = map[string]string{
	"repo_name_full": "repository.full_name",
	"repo_name":      "repository.name",
	"isprivate":      "repository.private",
	"owner_name":     "repository.owner.name",
	"owner_email":    "repository.full_name",
	"pusher_name":    "pusher.name",
	"pusher_email":   "pusher.email",
}

func validateHeaders(mode uint8, header http.Header) bool {
	var c int
	reqHeaders := requiredHeadersForMode[mode]
	for k := range header {
		for _, header := range reqHeaders {
			if strings.TrimSpace(strings.ToLower(k)) == header {
				c++
			}
		}
	}
	return c >= len(reqHeaders)
}

//Formats the variable names in the actions
func formatAction(subscription *Subscription, webhookData *WebhookData, actionCmd *string, action *ActionFileStruct) (requestValid, hitTrigger bool) {
	if !validateHeaders(subscription.Mode, webhookData.Header) {
		log.Println("Not all required headers were found!")
		return false, false
	}

	//Parse JSON to map
	var data map[string]interface{}
	err := json.Unmarshal([]byte(webhookData.Payload), &data)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		//Get the event from header
		var event string
		for h, v := range webhookData.Header {
			i := strings.ToLower(strings.TrimSpace(h))
			if i == eventHeader[subscription.Mode] {
				event = strings.TrimSpace(strings.ToLower(v[0]))
				break
			}
		}

		//Return if trigger wasn't hit
		if !checkTrigger(event, action.Source.Trigger) {
			return true, false
		}

		//Relpace aliases
		variables := GetVariablesFromCommand(*actionCmd)
		for i, vari := range variables {
			*actionCmd = strings.ReplaceAll(*actionCmd, vari, getReplacedAlias(vari))
			variables[i] = getReplacedAlias(vari)
		}

		//Fill empty map with variable names
		varValMap := make(map[string]string, len(variables))
		for _, v := range variables {
			varValMap[v] = "-"
		}

		//Fill variablenames with values from json
		if len(variables) > 0 {
			findJSONvalues("", data, &varValMap)
		}

		//In case a 'event' variable is used, fill it
		_, k := varValMap["event"]
		if k {
			varValMap["event"] = event
		}

		//Replace all custom variables with json values
		for _, vari := range variables {
			jsonVal := varValMap[vari]

			if action.Shell.SafeMode {
				for _, toReplace := range []string{";", "`", "$("} {
					jsonVal = strings.ReplaceAll(jsonVal, toReplace, "")
				}
			}

			*actionCmd = strings.ReplaceAll(*actionCmd, "%"+vari+"%", jsonVal)
		}
	}

	return true, true
}

//Returns true if trigger was hit
func checkTrigger(event string, triggers []string) bool {
	// Handle trigger
	if len(triggers) > 0 {
		for _, trigger := range triggers {
			if strings.ToLower(trigger) == event {
				return true
			}
		}
		return false
	}

	return true
}

//Fill varlist with json values
func findJSONvalues(name string, payload map[string]interface{}, varlist *map[string]string) {
	for k, v := range payload {
		if v == nil {
			continue
		}
		reft := reflect.TypeOf(v)
		if reft == nil {
			continue
		}
		if reft.Kind() == reflect.Map {
			findJSONvalues(name+k+".", (v).(map[string]interface{}), varlist)
			continue
		}

		currID := name + k
		_, has := (*varlist)[currID]
		if has {
			val := reflectToString(reflect.ValueOf(v))
			if len(val) > 0 {
				(*varlist)[currID] = reflectToString(reflect.ValueOf(v))
			}
		}
	}
}

func replaceAliases(variables []string) []string {
	ret := []string{}
	for _, variable := range variables {
		ret = append(ret, getReplacedAlias(variable))
	}
	return ret
}

func getReplacedAlias(str string) string {
	v, has := payloadAlias[strings.ToLower(str)]
	if has {
		return v
	}
	return str
}
