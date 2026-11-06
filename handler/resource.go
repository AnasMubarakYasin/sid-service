package handler

import (
	"encoding/json"
	"log"
	"sid/service/feature"
	"time"

	"github.com/gin-gonic/gin"
)

type Resource struct {
	f *feature.Resource
}

func NewResource(feature *feature.Resource) *Resource {
	return &Resource{f: feature}
}

func (h *Resource) Event(c *gin.Context) {
	cr := c.Request.Context()
	sub := func(ev *feature.ResourceEvent) {
		log.Printf("Resource Event Stream: %+v", ev.Data)
		p, e := json.Marshal(ev.Data)
		if e != nil {
			c.SSEvent("error", gin.H{"name": "parse", "message": "failed parsing"})
		} else {
			c.SSEvent(ev.Type, p)
		}
		c.Writer.Flush()
	}

	h.f.Subscribe(c, sub)

	for {
		select {
		case <-cr.Done():
			h.f.Unsubscribe(c, sub)
			return
		case <-time.After(60 * time.Second):
			c.SSEvent("ping", gin.H{})
			c.Writer.Flush()
		}
	}
}
