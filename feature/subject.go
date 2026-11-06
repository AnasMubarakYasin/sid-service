package feature

import (
	"context"
	"sid/service/model"
	"sid/service/repository"
)

type Subject struct {
	s  *repository.Subject
	ch chan *SubjectEvent
	sb map[any]ssub
}

func NewSubject(s *repository.Subject) *Subject {
	Subject := Subject{s, make(chan *SubjectEvent), make(map[any]ssub)}
	go Subject.event_loop()
	return &Subject
}

type ssub func(*SubjectEvent)

type ParamCreateSubject struct {
	Title       string   `bun:",notnull" form:"title" json:"title" xml:"title"`
	Description string   `bun:",notnull" form:"description" json:"description" xml:"description"`
	Video       string   `bun:",notnull" form:"video" json:"video" xml:"video"`
	Questions   []string `bun:",notnull" form:"questions" json:"questions" xml:"questions"`
}

type ParamUpdateSubject struct {
	ID          int64    `bun:",pk,autoincrement" form:"id" json:"id" xml:"id"`
	Title       string   `bun:",notnull" form:"title" json:"title" xml:"title"`
	Description string   `bun:",notnull" form:"description" json:"description" xml:"description"`
	Video       string   `bun:",notnull" form:"video" json:"video" xml:"video"`
	Questions   []string `bun:",notnull" form:"questions" json:"questions" xml:"questions"`
}

type SubjectEvent struct {
	Type string `form:"type" json:"type" xml:"type"`
	Data any    `form:"data" json:"data" xml:"data"`
}

func (f *Subject) Subscribe(c context.Context, s ssub) (e error) {

	f.sb[c] = s

	ct, e := f.Count(c)

	if e != nil {
		return e
	}

	s(&SubjectEvent{Type: "sync", Data: map[string]int{"count": ct}})

	return nil
}

func (f *Subject) Unsubscribe(c context.Context, s ssub) (e error) {

	delete(f.sb, c)

	return nil
}

func (f *Subject) event_loop() {
	for ch := range f.ch {
		for _, s := range f.sb {
			s(ch)
		}
	}
}

func (f *Subject) Count(c context.Context) (count int, e error) {

	count, e = f.s.SubjectCount(c)

	if e != nil {
		return count, e
	}

	return count, nil
}

func (f *Subject) Create(c context.Context, p *ParamCreateSubject) (r *model.Subject, e error) {
	r, e = f.s.SubjectCreate(c, &model.Subject{
		Title:       p.Title,
		Description: p.Description,
		Video:       p.Video,
		Questions:   p.Questions,
	})

	if e != nil {
		return nil, e
	}

	f.ch <- &SubjectEvent{
		Type: "create",
		Data: r,
	}

	return r, nil
}

func (f *Subject) Update(c context.Context, p *ParamUpdateSubject) (r *model.Subject, e error) {
	r, e = f.s.SubjectUpdate(c, &model.Subject{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		Video:       p.Video,
		Questions:   p.Questions,
	})

	if e != nil {
		return nil, e
	}

	return r, nil
}

func (f *Subject) Delete(c context.Context, id *int64) (r *model.Subject, e error) {
	r, e = f.s.SubjectDelete(c, id)

	if e != nil {
		return nil, e
	}

	f.ch <- &SubjectEvent{
		Type: "delete",
		Data: r,
	}

	return r, nil
}

func (f *Subject) List(c context.Context) (r *[]model.Subject, e error) {
	r, e = f.s.SubjectAll(c)

	if e != nil {
		return nil, e
	}

	return r, nil
}

func (f *Subject) Get(c context.Context, id *int64) (r *model.Subject, e error) {
	r, e = f.s.SubjectFindByID(c, id)

	if e != nil {
		return nil, e
	}

	return r, nil
}
