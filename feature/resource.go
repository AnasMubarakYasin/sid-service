package feature

import (
	"context"
	"log"
)

type Resource struct {
	usr *User
	pro *Profile
	exp *Experiment
	sub *Subject

	ch chan *ResourceEvent
	sb map[any]rsub
}

type rsub func(*ResourceEvent)

type ResourceEvent struct {
	Type string `form:"type" json:"type" xml:"type"`
	Data any    `form:"data" json:"data" xml:"data"`
}

func NewResource(
	usr *User,
	pro *Profile,
	exp *Experiment,
	sub *Subject,
) *Resource {
	Resource := Resource{usr, pro, exp, sub, make(chan *ResourceEvent), make(map[any]rsub)}
	go Resource.event_loop()
	return &Resource
}

func (f *Resource) Subscribe(c context.Context, s rsub) (e error) {

	f.sb[c] = s

	cexp, e := f.exp.Count(c)
	if e != nil {
		return e
	}
	csub, e := f.sub.Count(c)
	if e != nil {
		return e
	}

	s(&ResourceEvent{Type: "stat", Data: map[string]int{"experiment": cexp, "report": 0, "subject": csub}})

	return nil
}

func (f *Resource) Unsubscribe(c context.Context, s rsub) (e error) {

	delete(f.sb, c)

	return nil
}

func (f *Resource) event_loop() {
	c := context.Background()
	cexp := 0
	csub := 0
	expsb := func(ee *ExperimentEvent) {
		cexp, e := f.exp.Count(c)
		if e != nil {
			return
		}
		log.Printf("Forward Experiment Event: %+v", ee)
		f.ch <- &ResourceEvent{Type: "stat", Data: map[string]int{"experiment": cexp, "report": 0, "subject": csub}}
	}
	subsb := func(ee *SubjectEvent) {
		csub, e := f.sub.Count(c)
		if e != nil {
			return
		}
		f.ch <- &ResourceEvent{Type: "stat", Data: map[string]int{"experiment": cexp, "report": 0, "subject": csub}}
	}
	f.exp.Subscribe(c, expsb)
	f.sub.Subscribe(c, subsb)

	// go func() {
	// 	<-c.Done()
	// 	f.exp.Unsubscribe(c, expsb)
	// 	f.sub.Unsubscribe(c, subsb)
	// }()
	for ch := range f.ch {
		log.Printf("Resource Event: %+v", ch)
		for _, s := range f.sb {
			// log.Printf("Resource Event Broadcast: %+v", ch)
			s(ch)
		}
	}
}
