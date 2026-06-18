package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"sid/service/feature"
	"sid/service/model"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Experiment struct {
	f *feature.Experiment
}

func NewExperiment(feature *feature.Experiment) *Experiment {
	e := &Experiment{f: feature}
	return e
}

func (h *Experiment) CreateLab(c *gin.Context) {
	p := &feature.ParamCreateLabExperiment{}

	if e := c.ShouldBindJSON(p); e != nil {
		_ = c.Error(ErrBadRequest)
		return
	}

	usr, e := GetUser(c)
	if e != nil {
		_ = c.Error(ErrUnauthorized)
		return
	}

	r, e := h.f.CreateLab(c, p, usr)
	if e != nil {
		_ = c.Error(e)
		return
	}

	labs[r.ID] = newLab(h.f.Collaboration(c, usr, r))

	OK(c, r)
}

func (h *Experiment) RemoveLab(c *gin.Context) {
	p := &ParamID{}

	if err := c.ShouldBindUri(p); err != nil {
		_ = c.Error(ErrBadRequest)
		return
	}

	usr, e := GetUser(c)
	if e != nil {
		_ = c.Error(ErrUnauthorized)
		return
	}

	r, e := h.f.DeleteLab(c, &p.ID, &usr.ID)
	if e != nil {
		_ = c.Error(e)
		return
	}

	labs[r.ID].destroy()

	delete(labs, r.ID)

	OK(c, r)
}

func (h *Experiment) Collaboration(c *gin.Context) {
	ws, e := GetWebsocket(c)
	if e != nil {
		log.Printf("Failed get websocket")
		return
	}

	p := &ParamID{}

	if e := c.ShouldBindUri(p); e != nil {
		log.Printf("Bad uri: %+v", e)
		return
	}

	usr, e := GetUser(c)
	if e != nil {
		log.Printf("Failed get user: %v\n", p.ID)
		return
	}

	lab, ok := labs[p.ID]
	if !ok {
		log.Printf("Global Lab not used: %+v", p.ID)
		exp, e := h.f.Get(c, &p.ID)
		if e != nil {
			log.Printf("Lab not found: %+v", p.ID)
			return
		}
		if exp == nil {
			log.Printf("Experiment not found: %+v", p.ID)
			return
		}
		lab = newLab(h.f.Collaboration(c, usr, exp))
		labs[exp.ID] = lab
	} else {
		log.Printf("Global Lab used: %+v", p.ID)
	}

	log.Printf("Collaboration: %+v", usr.ID)

	lab.mu.Lock()
	lab.join(ws, usr)
	// defer lab.leave(ws, usr)
	lab.osync(ws, "all")
	lab.mu.Unlock()

	log.Printf("Synchronized: %+v", usr.ID)

	// lab.heartbeat(ws)
	lab.repl(ws, usr)
}

func (h *Experiment) Event(c *gin.Context) {
	cr := c.Request.Context()
	sub := func(ev *feature.ExperimentEvent) {
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

func (h *Experiment) Create(c *gin.Context) {
	p := &feature.ParamCreateExperiment{}

	if e := c.ShouldBindJSON(p); e != nil {
		_ = c.Error(ErrBadRequest)
		return
	}

	r, e := h.f.Create(c, p)

	if e != nil {
		Fail(c, 500, e)
		return
	}

	OK(c, r)
}

func (h *Experiment) Update(c *gin.Context) {
	p := &ParamID{}

	if err := c.ShouldBindUri(p); err != nil {
		_ = c.Error(ErrBadRequest)
		return
	}

	pp := &feature.ParamUpdateExperiment{}

	if e := c.ShouldBindJSON(pp); e != nil {
		_ = c.Error(ErrBadRequest)
		return
	}

	pp.ID = p.ID

	r, e := h.f.Update(c, pp)

	if e != nil {
		Fail(c, 500, e)
		return
	}

	OK(c, r)
}

func (h *Experiment) Delete(c *gin.Context) {
	p := &ParamID{}

	if err := c.ShouldBindUri(p); err != nil {
		_ = c.Error(ErrBadRequest)
		return
	}

	r, e := h.f.Delete(c, &p.ID)

	if e != nil {
		Fail(c, 500, e)
		return
	}

	OK(c, r)
}

func (h *Experiment) List(c *gin.Context) {
	r, e := h.f.List(c)

	if e != nil {
		Fail(c, 500, e)
		return
	}

	OKM(c, r, &Meta{})
}

func (h *Experiment) Get(c *gin.Context) {
	p := &ParamID{}

	if err := c.ShouldBindUri(p); err != nil {
		_ = c.Error(ErrBadRequest)
		return
	}

	r, e := h.f.Get(c, &p.ID)

	if e != nil {
		Fail(c, 500, e)
		return
	}

	OK(c, r)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var collabs = make(map[int64]*feature.ExperimentCollaboration)
var labs = make(map[int64]*ExperimentLab)

var (
	messageTypeLeave = 0
	messageTypeJoin  = 1
	messageTypeOP    = 2
)

type ExperimentLab struct {
	mu     sync.RWMutex
	cl     map[*websocket.Conn]*model.User
	collab *feature.ExperimentCollaboration
	// usr map[int64]*model.User
	// exp *model.Experiment
}

func newLab(collab *feature.ExperimentCollaboration) *ExperimentLab {
	return &ExperimentLab{
		mu:     sync.RWMutex{},
		cl:     make(map[*websocket.Conn]*model.User),
		collab: collab,
	}
}

func (h *ExperimentLab) join(ws *websocket.Conn, usr *model.User) {
	log.Printf("join: %+v", usr.ID)
	h.cl[ws] = usr
	h.broadcast(h.collab.Join(usr))
}

func (h *ExperimentLab) leave(ws *websocket.Conn, usr *model.User) {
	log.Printf("leave: %+v", usr.ID)
	h.broadcast(h.collab.Leave(usr))
	delete(h.cl, ws)
}

func (h *ExperimentLab) sync(a string) {
	h.broadcast(h.collab.DownSync(a))
}

func (h *ExperimentLab) osync(ws *websocket.Conn, a string) {
	h.send(ws, h.collab.DownSync(a))
}

func (h *ExperimentLab) destroy() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.broadcast(h.collab.Destroy())
	for ws := range h.cl {
		delete(h.cl, ws)
	}
}

func (h *ExperimentLab) kick(usr *model.User, msg any) {
	for cl, u := range h.cl {
		if u.ID != usr.ID {
			continue
		}
		h.send(cl, msg)
		return
	}
	log.Printf("failed kick: %+v", usr)
}

func (h *ExperimentLab) send(ws *websocket.Conn, msg any) {
	if err := ws.WriteJSON(msg); err != nil {
		log.Printf("Send error: %+v", err)
	}
}

func (h *ExperimentLab) publish(ws *websocket.Conn, msg any) {
	u := h.cl[ws]
	delete(h.cl, ws)
	for cl, u := range h.cl {
		log.Printf("broadcast: %+v", u.ID)
		if err := cl.WriteJSON(msg); err != nil {
			log.Printf("Broadcast error: %+v", err)
		}
	}
	h.cl[ws] = u
}

func (h *ExperimentLab) broadcast(msg any) {
	log.Printf("subscribers: %+v", len(h.cl))
	for cl, u := range h.cl {
		log.Printf("broadcast: %+v", u.ID)
		if err := cl.WriteJSON(msg); err != nil {
			log.Printf("Broadcast error: %+v", err)
		}
	}
}

func (h *ExperimentLab) repl(ws *websocket.Conn, usr *model.User) {
	for {
		// ws.SetReadDeadline(time.Now().Add(5 * time.Second))
		ws.SetReadDeadline(time.Now().Add(5 * time.Minute))
		msg := &feature.ExperimentCollaborationMessage{}
		e := ws.ReadJSON(msg)
		if e != nil {
			log.Printf("Read Loop Error: %+v", e)
			h.mu.Lock()
			h.leave(ws, usr)
			h.mu.Unlock()
			break
		}
		log.Printf("Read: %+v", msg)
		h.mu.Lock()
		switch msg.Type {
		case "sync":
			s := &feature.ExperimentSync{}
			if msg.Attr == "groups" {
				s.Groups = msg.Groups
			}
			if msg.Attr == "subject" {
				s.Subject = msg.Subject
			}
			if msg.Attr == "status" {
				s.Status = msg.Status
			}
			if msg.Attr == "questions" {
				s.Questions = msg.Questions
			}
			if msg.Attr == "answers" {
				s.Answers = msg.Answers
			}
			if msg.Attr == "simulation" {
				s.Simulation = msg.Simulation
			}
			if msg.Attr == "collections" {
				s.Collections = msg.Collections
			}
			if msg.Attr == "conclusion" {
				s.Conclusion = msg.Conclusion
			}
			rpl := h.collab.UpSync(msg.Attr, s, usr)
			if rpl != nil {
				h.send(ws, rpl)
				break
			}
			h.broadcast(msg)
			break
		case "kick":
			rpl := h.collab.Kick(usr, msg.User)
			if rpl.Type == "error" {
				h.send(ws, rpl)
				break
			}
			h.kick(rpl.User, rpl)
			h.leave(ws, rpl.User)
			break
		case "request":
			path := strings.Split(msg.RequestPath, ":")
			switch path[0] {
			case "sync":
				res := h.collab.DownSync("all")
				res.Type = "response"
				res.ResponsePath = msg.RequestPath
				res.ResponseBody = nil
				h.send(ws, res)
			}
			break
		case "publish":
			h.publish(ws, h.collab.Publish(msg))
		default:
			log.Printf("Read Unhandled: %+v", msg.Type)
		}
		h.mu.Unlock()
		// log.Printf("Read Experiment: %+v", msg.Experiment)
		// log.Printf("Read Experiment Groups: %v", msg.Experiment.Groups)
	}
}

func (h *ExperimentLab) heartbeat(conn *websocket.Conn) {
	const (
		pongWait   = 5 * time.Minute
		pingPeriod = (pongWait * 9) / 10 // must be less than pongWait
	)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	// Start a goroutine to send pings.
	go func() {
		ticker := time.NewTicker(pingPeriod)
		defer ticker.Stop()
		for range ticker.C {
			conn.SetWriteDeadline(time.Now().Add(pongWait))
			log.Printf("Pinging\n")
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}()
}
