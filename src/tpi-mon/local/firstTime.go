package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"tpi-mon/pkg/util"
)

const completeSetupMsg = `
########################################################################
########################################################################

Congratulations!

Your client has been successfully registered. To complete it, please go to:
%s


########################################################################
########################################################################
`

func firstTime(cfg *config) error {
	rsp, err := http.Post(cfg.CloudBaseURL+"/clients", "application/json", nil)
	if err != nil {
		return err
	}

	defer rsp.Body.Close()
	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		panic(err)
	}

	if rsp.StatusCode != 200 {
		return fmt.Errorf("registration request returned unexpected status %v: body:\n%v", rsp.StatusCode, string(data))
	}

	reg := struct {
		SiteID   string
		Token    string
		SetupURL string
	}{}

	if err := json.Unmarshal(data, &reg); err != nil {
		return fmt.Errorf("Unable to parse registration response: %v. Body:\n%v", err, string(data))
	}

	cfg.SiteID = reg.SiteID
	cfg.CloudToken = reg.Token

	if err = saveAuthConfig(cfg); err != nil {
		return err
	}

	url := fmt.Sprintf("%s/%s", cfg.CloudBaseURL, reg.SetupURL)

	fmt.Printf(completeSetupMsg, url)

	return nil
}

func saveAuthConfig(cfg *config) error {

	authCfg := config{SiteID: cfg.SiteID, CloudToken: cfg.CloudToken}

	data, err := json.MarshalIndent(authCfg, "", "  ")
	if err != nil {
		return err
	}

	fname, err := util.GetDefaultConfigFilename(appName)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fname, data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
