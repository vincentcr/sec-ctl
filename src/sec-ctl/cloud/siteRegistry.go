package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"sec-ctl/cloud/db"
	"sec-ctl/pkg/sites"
	"sec-ctl/pkg/ws"
)

// type siteRegistry struct {
// 	sites map[string]*remoteSite
// }

// func newRegistry() *siteRegistry {
// 	return &siteRegistry{
// 		sites: map[string]*remoteSite{},
// 	}
// }

// func (r *siteRegistry) getSite(id string) (*remoteSite, bool) {
// 	s, ok := r.sites[id]
// 	return s, ok
// }

// func (r *siteRegistry) addSite(c *remoteSite) {
// 	c.registry = r
// 	r.sites[c.id] = c
// }

// func (r *siteRegistry) removeSite(c *remoteSite) {
// 	delete(r.sites, c.id)
// }

const queueNameSiteRemoved = "sites.removed"

type siteRegistry struct {
	db             *db.DB
	queue          *queue
	connectedSites sync.Map
}

func newRegistry(dbConn *db.DB, queue *queue) *siteRegistry {

	sr := &siteRegistry{
		db:             dbConn,
		queue:          queue,
		connectedSites: sync.Map{},
	}

	queue.startConsumeLoop(queueNameSiteRemoved, func(msg qMessage) error {
		siteID := db.UUID(msg.data)
		sr.connectedSites.Delete(siteID)
		return nil
	})

	return sr
}

func (r *siteRegistry) initRemoteSite(site db.Site, conn *ws.Conn) {

	remoteSite := newRemoteSite(site, conn, r.queue)
	r.connectedSites.Store(site.ID, remoteSite)

	r.queue.startConsumeLoop(getSiteQueueName(site.ID, "events"), func(msg qMessage) error {
		var evt sites.Event
		if err := json.Unmarshal(msg.data, &evt); err != nil {
			logger.Panicf("failed to parse event from json %v: %v", msg.data, err)
		}

		return r.db.SaveEvent(string(evt.Level), evt.Time, site.ID, evt)
	})
}

func (r *siteRegistry) getSite(user db.User, id db.UUID) (db.Site, error) {

	s, err := r.db.FetchSiteByID(id)
	if err != nil {
		return db.Site{}, err
	}

	if s.OwnerID != user.ID {
		return db.Site{}, fmt.Errorf("Unauthorized")
	}

	return s, nil
}

func (r *siteRegistry) sendCommand(id db.UUID, cmd sites.UserCommand) error {

	data, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	expires := time.Now().Add(60 * time.Second)
	queueName := getSiteQueueName(id, "commands")
	return r.queue.publishEx(queueName, data, expires)
}

func (r *siteRegistry) getLatestEvents(id db.UUID, max uint) ([]db.Event, error) {
	return r.db.GetLatestEvents(id, max)

}
