package main

import "tpi-mon/cloud/db"

type siteController struct {
	site      db.Site
	connector remoteSite
	connected bool
}

func (sc siteController) getState() {

}
