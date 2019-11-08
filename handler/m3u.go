package handler

import (
	"net/http"

	"fmt"

	"io/ioutil"
	"log"

	"github.com/boscolai/m3uproxy/filter"
	"github.com/spf13/viper"
)

type account struct {
	ID     string
	M3U    string
	EPG    string
	Groups []string
}

var accountMap map[string]account

func HandleM3U(w http.ResponseWriter, req *http.Request) {
	acct := getUserAccount(w, req)
	if acct == nil {
		return
	}

	if len(acct.Groups) > 0 {
		tracks, err := filter.FilterByGroupNames(acct.M3U, acct.Groups)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "group filtering failed: %s\n", err)
		}
		filter.WriteTracks(w, tracks)
		return
	}

	// no filtering, just fetch and return
	fetchURL(w, acct.M3U)
}

func HandleEPG(w http.ResponseWriter, req *http.Request) {
	acct := getUserAccount(w, req)
	if acct == nil {
		return
	}
	fetchURL(w, acct.EPG)
}

func HandleGroups(w http.ResponseWriter, req *http.Request) {
	acct := getUserAccount(w, req)
	if acct == nil {
		return
	}

	groups, err := filter.GetGroups(acct.M3U)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Get groups failed: %s", err)
	}
	for k, v := range groups {
		fmt.Fprintf(w, "%s: %d\n", k, v)
	}
}

func InitHandlers(v *viper.Viper) {
	if err := v.UnmarshalKey("accounts", &accountMap); err != nil {
		log.Fatalf("failed to fetch accounts from config: %s", err)
	}
	log.Printf("user accounts: %v\n", accountMap)
}

func getUserAccount(w http.ResponseWriter, req *http.Request) *account {
	var user string
	if user = req.URL.Query().Get("u"); user == "" {
		w.WriteHeader(http.StatusForbidden)
		return nil
	}

	var acct account
	var exists bool
	if acct, exists = accountMap[user]; !exists {
		w.WriteHeader(http.StatusForbidden)
		return nil
	}
	return &acct
}

func fetchURL(w http.ResponseWriter, url string) {
	resp, err := http.Get(url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to fetch data for %s: %s", url, err)
		return
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to read data from url for %s: %s", url, err)
		return
	}

	w.Write(data)
}
