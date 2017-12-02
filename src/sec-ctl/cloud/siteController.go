package main

import "sec-ctl/cloud/db"

type siteController struct {
	site      db.Site
	connector remoteSite
	connected bool
}

func (sc siteController) getState() {

}
