package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"turboenigma/model"
	"turboenigma/provider"
)

type Gitlab struct {
	provider provider.Provider
}

func NewGitlab(provider provider.Provider) *Gitlab {
	return &Gitlab{provider: provider}
}

func (g *Gitlab) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	err := func() (err error) {
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			return err
		}

		if string(body) == "" {
			return fmt.Errorf("Body is missing")
		}

		mr, err := jsonDecode(string(body))
		if err != nil {
			return err
		}

		if mr.EventType != "merge_request" {
			fmt.Fprint(writer, "We just care about merge_request events")
			return
		}

		switch mr.ObjectAttributes.Action {
		case "open":
			if err = g.provider.NotifyMergeRequestOpened(mr); err != nil {
				return err
			}
		case "approved":
			if err = g.provider.NotifyMergeRequestApproved(mr); err != nil {
				return err
			}
			fmt.Fprint(writer, "Reacting to approved event")
			return
		case "unapproved":
			if err = g.provider.NotifyMergeRequestUnapproved(mr); err != nil {
				return err
			}
			fmt.Fprint(writer, "Reacting to unapproved event")
			return
		case "close":
			if err = g.provider.NotifyMergeRequestClose(mr); err != nil {
				return err
			}
			fmt.Fprint(writer, "Reacting to close event")
			return
		case "reopen":
			if err = g.provider.NotifyMergeRequestReopen(mr); err != nil {
				return err
			}
			fmt.Fprint(writer, "Reacting to reopen event")
			return
		case "update":
			if err = g.provider.NotifyMergeRequestUpdate(mr); err != nil {
				return err
			}
			fmt.Fprint(writer, "Reacting to update event")
			return
		case "approval":
			if err = g.provider.NotifyMergeRequestApproval(mr); err != nil {
				return err
			}
			fmt.Fprint(writer, "Reacting to approval event")
			return
		case "unapproval":
			if err = g.provider.NotifyMergeRequestUnapproval(mr); err != nil {
				return err
			}
			fmt.Fprint(writer, "Reacting to unapproval event")
			return
		case "merge":
			if err = g.provider.NotifyMergeRequestMerged(mr); err != nil {
				return err
			}
			fmt.Fprint(writer, "Reacting to merge event")
			return
		default:
			fmt.Fprint(writer, fmt.Sprintf("We cannot handle %s event action", mr.ObjectAttributes.Action))
			return
		}

		fmt.Fprint(writer, "OK")

		return
	}()

	if err != nil {
		http.Error(writer, fmt.Sprintf("Error -> %s", err.Error()), http.StatusBadRequest)
	}
}

func jsonDecode(jsonString string) (mergeRequest model.MergeRequestInfo, err error) {
	err = json.Unmarshal([]byte(jsonString), &mergeRequest)

	return
}
