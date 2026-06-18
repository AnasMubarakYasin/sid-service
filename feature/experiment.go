package feature

import (
	"context"
	"log"
	"sid/service/model"
	"sid/service/repository"
	"slices"
)

type Experiment struct {
	r  *repository.Experiment
	s  *repository.Subject
	u  *repository.User
	ch chan *ExperimentEvent
	sb map[any]sub
}

func NewExperiment(r *repository.Experiment, s *repository.Subject, u *repository.User) *Experiment {
	o := Experiment{r, s, u, make(chan *ExperimentEvent), make(map[any]sub)}
	r.ExperimentDeleteAll(context.TODO())
	go o.event_loop()
	return &o
}

type sub func(*ExperimentEvent)

type ParamCreateLabExperiment struct {
	AdminID int64 `form:"admin_id" json:"admin_id" xml:"admin_id"`

	Image       string `form:"image" json:"image" xml:"image"`
	Title       string `form:"title" json:"title" xml:"title"`
	Description string `form:"description" json:"description" xml:"description"`
}

type ParamCreateExperiment struct {
	SubjectID  int64   `form:"subject_id" json:"subject_id" xml:"subject_id"`
	AdminIDS   []int64 `form:"admin_ids" json:"admin_ids" xml:"admin_ids"`
	TeacherIDS []int64 `form:"teacher_ids" json:"teacher_ids" xml:"teacher_ids"`
	StudentIDS []int64 `form:"student_ids" json:"student_ids" xml:"student_ids"`

	Image       string `form:"image" json:"image" xml:"image"`
	Title       string `form:"title" json:"title" xml:"title"`
	Description string `form:"description" json:"description" xml:"description"`
	Status      string `form:"status" json:"status" xml:"status"`

	Rules  []string       `form:"rules" json:"rules" xml:"rules"`
	Groups map[string]any `form:"groups" json:"groups" xml:"groups"`

	Questions   []string       `form:"questions" json:"questions" xml:"questions"`
	Hypotesis   []string       `form:"hypotesis" json:"hypotesis" xml:"hypotesis"`
	Simulation  map[string]any `form:"simulation" json:"simulation" xml:"simulation"`
	Data        map[string]any `form:"data" json:"data" xml:"data"`
	Explanation map[string]any `form:"explanation" json:"explanation" xml:"explanation"`
}

type ParamUpdateExperiment struct {
	ID int64 `form:"id" json:"id" xml:"id"`
}

func (f *Experiment) CreateLab(c context.Context, p *ParamCreateLabExperiment, usr *model.User) (r *model.Experiment, e error) {
	d := &model.Experiment{
		AdminIDS:    []int64{usr.ID},
		Image:       p.Image,
		Title:       p.Title,
		Description: p.Description,
		Status:      "waiting",
	}
	switch usr.Role {
	case model.UserRoleAdmin:
		break
	case model.UserRoleTeacher:
		if !slices.Contains(d.TeacherIDS, usr.ID) {
			d.TeacherIDS = append(d.TeacherIDS, usr.ID)
		}
	case model.UserRoleStudent:
		if !slices.Contains(d.StudentIDS, usr.ID) {
			d.StudentIDS = append(d.StudentIDS, usr.ID)
		}
	}

	r, e = f.r.ExperimentCreate(c, d)

	if e != nil {
		return nil, e
	}

	f.ch <- &ExperimentEvent{
		Type: "create",
		Data: r,
	}

	return r, nil
}

func (f *Experiment) DeleteLab(c context.Context, eid *int64, uid *int64) (r *model.Experiment, e error) {
	r, e = f.r.ExperimentFindByID(c, eid)

	if r == nil {
		return nil, ErrLabNExist
	}

	h := slices.ContainsFunc(r.AdminIDS, func(aid int64) bool {
		return aid == *uid
	})

	if !h {
		return nil, ErrNotAdmin
	}

	r, e = f.r.ExperimentDelete(c, eid)

	if e != nil {
		return nil, e
	}

	f.ch <- &ExperimentEvent{
		Type: "delete",
		Data: r,
	}

	return r, nil
}

func (f *Experiment) Collaboration(c context.Context, u *model.User, e *model.Experiment) (r *ExperimentCollaboration) {
	return &ExperimentCollaboration{u, e, f.u, f.s}
}

func (f *Experiment) Create(c context.Context, p *ParamCreateExperiment) (r *model.Experiment, e error) {
	r, e = f.r.ExperimentCreate(c, &model.Experiment{})

	if e != nil {
		return nil, e
	}

	return r, nil
}

func (f *Experiment) Update(c context.Context, p *ParamUpdateExperiment) (r *model.Experiment, e error) {
	r, e = f.r.ExperimentUpdate(c, &model.Experiment{})

	if e != nil {
		return nil, e
	}

	return r, nil
}

func (f *Experiment) Delete(c context.Context, id *int64) (r *model.Experiment, e error) {
	r, e = f.r.ExperimentDelete(c, id)

	if e != nil {
		return nil, e
	}

	return r, nil
}

func (f *Experiment) List(c context.Context) (r *[]model.Experiment, e error) {
	r, e = f.r.ExperimentAll(c)

	if e != nil {
		return nil, e
	}

	return r, nil
}

func (f *Experiment) Get(c context.Context, id *int64) (r *model.Experiment, e error) {
	r, e = f.r.ExperimentFindByID(c, id)

	if e != nil {
		return nil, e
	}

	return r, nil
}

func (f *Experiment) Count(c context.Context) (count int, e error) {

	count, e = f.r.ExperimentCount(c)

	if e != nil {
		return count, e
	}

	return count, nil
}

func (f *Experiment) Subscribe(c context.Context, s sub) (e error) {

	f.sb[c] = s

	ct, e := f.Count(c)

	if e != nil {
		return e
	}

	s(&ExperimentEvent{Type: "sync", Data: map[string]int{"count": ct}})

	return nil
}

func (f *Experiment) Unsubscribe(c context.Context, s sub) (e error) {

	delete(f.sb, c)

	return nil
}

func (f *Experiment) event_loop() {
	for ch := range f.ch {
		for _, s := range f.sb {
			s(ch)
		}
	}
}

var (
	ErrNotAdmin  = &Error{Code: "Not Admin", Message: "Not Admin Lab"}
	ErrLabNExist = &Error{Code: "Not Exists", Message: "Lab Not Exists"}
	ErrNSSubj    = &Error{Code: "Subject Not Set", Message: "Suject not set up yet"}
)

type ExperimentEvent struct {
	Type string `form:"type" json:"type" xml:"type"`
	Data any    `form:"data" json:"data" xml:"data"`
	// Data  *model.Experiment `form:"data,omitempty" json:"data,omitempty" xml:"data,omitempty"`
	// Count int               `form:"count,omitempty" json:"count,omitempty" xml:"count,omitempty"`
}
type ExperimentGroup struct {
	ID       int64         `form:"id" json:"id" xml:"id"`
	Name     string        `form:"name" json:"name" xml:"name"`
	Students *[]model.User `form:"students,omitempty" json:"students,omitempty" xml:"students,omitempty"`
}

type ExperimentError struct {
	Name    string `form:"name" json:"name" xml:"name"`
	Message string `form:"message" json:"message" xml:"message"`
}

func (err *ExperimentError) Error() string {
	return err.Message
}

type ExperimentCollaborationMessage struct {
	Type        string          `form:"type,omitempty" json:"type,omitempty" xml:"type,omitempty"`
	Attr        string          `form:"attr,omitempty" json:"attr,omitempty" xml:"attr,omitempty"`
	User        *model.User     `form:"user,omitempty" json:"user,omitempty" xml:"user,omitempty"`
	Teachers    *[]model.User   `form:"teachers,omitempty" json:"teachers,omitempty" xml:"teachers,omitempty"`
	Students    *[]model.User   `form:"students,omitempty" json:"students,omitempty" xml:"students,omitempty"`
	Groups      *[]any          `form:"groups,omitempty" json:"groups,omitempty" xml:"groups,omitempty"`
	Subject     *model.Subject  `form:"subject,omitempty" json:"subject,omitempty" xml:"subject,omitempty"`
	Status      *string         `form:"status,omitempty" json:"status,omitempty" xml:"status,omitempty"`
	Questions   *[]string       `form:"questions,omitempty" json:"questions,omitempty" xml:"questions,omitempty"`
	Answers     *[]string       `form:"answers,omitempty" json:"answers,omitempty" xml:"answers,omitempty"`
	Simulation  *map[string]any `form:"simulation,omitempty" json:"simulation,omitempty" xml:"simulation,omitempty"`
	Collections *[][]string     `form:"collections,omitempty" json:"collections,omitempty" xml:"collections,omitempty"`
	Conclusion  *string         `form:"conclusion,omitempty" json:"conclusion,omitempty" xml:"conclusion,omitempty"`

	// Groups   *[]ExperimentGroup `form:"groups,omitempty" json:"groups,omitempty" xml:"groups,omitempty"`
	// Experiment *ExperimentSync  `form:"experiment,omitempty" json:"experiment,omitempty" xml:"experiment,omitempty"`

	Error *Error `form:"error,omitempty" json:"error,omitempty" xml:"error,omitempty"`

	RequestPath  string         `form:"requestPath,omitempty" json:"requestPath,omitempty" xml:"requestPath,omitempty"`
	RequestBody  map[string]any `form:"requestBody,omitempty" json:"requestBody,omitempty" xml:"requestBody,omitempty"`
	PublishPath  string         `form:"publishPath,omitempty" json:"publishPath,omitempty" xml:"publishPath,omitempty"`
	PublishBody  map[string]any `form:"publishBody,omitempty" json:"publishBody,omitempty" xml:"publishBody,omitempty"`
	ResponsePath string         `form:"responsePath,omitempty" json:"responsePath,omitempty" xml:"responsePath,omitempty"`
	ResponseBody map[string]any `form:"responseBody,omitempty" json:"responseBody,omitempty" xml:"responseBody,omitempty"`
	MessagePath  string         `form:"messagePath,omitempty" json:"messagePath,omitempty" xml:"messagePath,omitempty"`
	MessageBody  map[string]any `form:"messageBody,omitempty" json:"messageBody,omitempty" xml:"messageBody,omitempty"`
}
type ExperimentSync struct {
	Teachers    *[]model.User   `form:"teachers,omitempty" json:"teachers,omitempty" xml:"teachers,omitempty"`
	Students    *[]model.User   `form:"students,omitempty" json:"students,omitempty" xml:"students,omitempty"`
	Groups      *[]any          `form:"groups,omitempty" json:"groups,omitempty" xml:"groups,omitempty"`
	Subject     *model.Subject  `form:"subject,omitempty" json:"subject,omitempty" xml:"subject,omitempty"`
	Status      *string         `form:"status,omitempty" json:"status,omitempty" xml:"status,omitempty"`
	Questions   *[]string       `form:"questions,omitempty" json:"questions,omitempty" xml:"questions,omitempty"`
	Answers     *[]string       `form:"answers,omitempty" json:"answers,omitempty" xml:"answers,omitempty"`
	Simulation  *map[string]any `form:"simulation,omitempty" json:"simulation,omitempty" xml:"simulation,omitempty"`
	Collections *[][]string     `form:"collections,omitempty" json:"collections,omitempty" xml:"collections,omitempty"`
	Conclusion  *string         `form:"conclusion,omitempty" json:"conclusion,omitempty" xml:"conclusion,omitempty"`
}

type ExperimentCollaboration struct {
	u  *model.User
	e  *model.Experiment
	ru *repository.User
	rs *repository.Subject
}

func (c *ExperimentCollaboration) Join(u *model.User) (r *ExperimentCollaborationMessage) {
	switch u.Role {
	case model.UserRoleAdmin:
		break
	case model.UserRoleTeacher:
		if !slices.Contains(c.e.TeacherIDS, u.ID) {
			c.e.TeacherIDS = append(c.e.TeacherIDS, u.ID)
		}
	case model.UserRoleStudent:
		if !slices.Contains(c.e.StudentIDS, u.ID) {
			c.e.StudentIDS = append(c.e.StudentIDS, u.ID)
		}
	}
	return &ExperimentCollaborationMessage{
		Type: "join",
		User: u,
	}
}

func (c *ExperimentCollaboration) Leave(u *model.User) (r *ExperimentCollaborationMessage) {
	switch u.Role {
	case model.UserRoleAdmin:
		break
	case model.UserRoleTeacher:
		c.e.TeacherIDS = slices.DeleteFunc(c.e.TeacherIDS, func(id int64) bool {
			return id == u.ID
		})
	case model.UserRoleStudent:
		c.e.StudentIDS = slices.DeleteFunc(c.e.StudentIDS, func(id int64) bool {
			return id == u.ID
		})
	}
	return &ExperimentCollaborationMessage{
		Type: "leave",
		User: u,
	}
}

func (c *ExperimentCollaboration) Destroy() (r *ExperimentCollaborationMessage) {
	c.ru = nil
	c.rs = nil
	return &ExperimentCollaborationMessage{
		Type: "destroy",
	}
}

func (c *ExperimentCollaboration) Kick(usr *model.User, tgt *model.User) (r *ExperimentCollaborationMessage) {
	h := slices.ContainsFunc(c.e.AdminIDS, func(aid int64) bool {
		return aid == usr.ID
	})
	if !h {
		return &ExperimentCollaborationMessage{
			Type:  "error",
			Error: ErrNotAdmin,
		}
	}
	switch tgt.Role {
	case model.UserRoleAdmin:
		break
	case model.UserRoleTeacher:
		c.e.TeacherIDS = slices.DeleteFunc(c.e.TeacherIDS, func(id int64) bool {
			return id == tgt.ID
		})
	case model.UserRoleStudent:
		c.e.StudentIDS = slices.DeleteFunc(c.e.StudentIDS, func(id int64) bool {
			return id == tgt.ID
		})
	}
	return &ExperimentCollaborationMessage{
		Type: "kick",
		User: tgt,
	}
}

func (c *ExperimentCollaboration) DownSync(attr string) (r *ExperimentCollaborationMessage) {
	message := &ExperimentCollaborationMessage{
		Type: "sync",
		Attr: attr,
	}
	if attr == "all" || attr == "users" {
		teachers, err := c.ru.UserFindByIDS(context.TODO(), c.e.TeacherIDS)
		if err != nil {
			return &ExperimentCollaborationMessage{
				Type:  "error",
				Error: &Error{Code: "sync", Message: "failed find many t"},
			}
		}
		students, err := c.ru.UserFindByIDS(context.TODO(), c.e.StudentIDS)
		if err != nil {
			return &ExperimentCollaborationMessage{
				Type:  "error",
				Error: &Error{Code: "sync", Message: "failed find many s"},
			}
		}
		message.Teachers = teachers
		message.Students = students
	}
	if attr == "all" || attr == "groups" {
		message.Groups = &c.e.Groups
	}
	if attr == "all" || attr == "subject" {
		subject, err := c.rs.SubjectFindByID(context.TODO(), &c.e.SubjectID)
		if err != nil {
			return &ExperimentCollaborationMessage{
				Type:  "error",
				Error: &Error{Code: "sync", Message: "failed find one s"},
			}
		}
		message.Subject = subject
	}
	if attr == "all" || attr == "status" {
		message.Status = &c.e.Status
	}
	if attr == "all" || attr == "questions" {
		message.Questions = &c.e.Questions
	}
	if attr == "all" || attr == "answers" {
		if len(c.e.Hypotesis) > 0 {
			message.Answers = &c.e.Hypotesis
		}
	}
	if attr == "all" || attr == "simulation" {
		if len(c.e.Simulation) > 0 {
			message.Simulation = &c.e.Simulation
		}
	}
	if attr == "all" || attr == "collections" {
		if len(c.e.Data) > 0 {
			d, ok := c.e.Data["value"].([][]string)
			if ok {
				message.Collections = &d
			}
		}
	}
	if attr == "all" || attr == "conclusion" {
		if len(c.e.Explanation) > 0 {
			d, ok := c.e.Explanation["value"].(string)
			if ok {
				message.Conclusion = &d
			}
		}
	}
	log.Printf("DWSYNC %+v", message)
	return message
}
func (c *ExperimentCollaboration) UpSync(attr string, s *ExperimentSync, u *model.User) (r *ExperimentCollaborationMessage) {
	r = &ExperimentCollaborationMessage{}
	if attr == "all" || attr == "users" {
		for _, v := range *s.Teachers {
			c.e.TeacherIDS = append(c.e.TeacherIDS, v.ID)
		}
		for _, v := range *s.Teachers {
			c.e.StudentIDS = append(c.e.StudentIDS, v.ID)
		}
	}
	if attr == "all" || attr == "groups" {
		c.e.Groups = *s.Groups
	}
	if attr == "all" || attr == "subject" {
		h := slices.ContainsFunc(c.e.AdminIDS, func(aid int64) bool {
			return aid == u.ID
		})
		if !h {
			s, e := c.rs.SubjectFindByID(context.TODO(), &c.e.SubjectID)
			if e != nil {
				return &ExperimentCollaborationMessage{
					Type:  "error",
					Error: &Error{Code: "sync", Message: "failed find one s"},
				}
			}
			r.Type = "error"
			r.Attr = "subject"
			r.Subject = s
			r.Error = ErrNotAdmin
			return r
		}
		c.e.SubjectID = s.Subject.ID
	}
	if attr == "all" || attr == "status" {
		h := slices.ContainsFunc(c.e.AdminIDS, func(aid int64) bool {
			return aid == u.ID
		})
		if !h {
			r.Type = "error"
			r.Attr = "status"
			r.Status = &c.e.Status
			r.Error = ErrNotAdmin
			return r
		}
		if c.e.SubjectID == 0 {
			r.Type = "error"
			r.Attr = "status"
			r.Status = &c.e.Status
			r.Error = ErrNSSubj
			return r
		}
		c.e.Status = *s.Status
	}
	if attr == "all" || attr == "questions" {
		c.e.Questions = *s.Questions
	}
	if attr == "all" || attr == "answers" {
		c.e.Hypotesis = *s.Answers
	}
	if attr == "all" || attr == "simulation" {
		c.e.Simulation = *s.Simulation
	}
	if attr == "all" || attr == "collections" {
		c.e.Data = map[string]any{"value": *s.Collections}
	}
	if attr == "all" || attr == "conclusion" {
		c.e.Explanation = map[string]any{"value": *s.Conclusion}
	}
	log.Printf("UPSYNC %+v", c.e)
	return nil
}

func (c *ExperimentCollaboration) Request(m *ExperimentCollaborationMessage) (r *ExperimentCollaborationMessage) {

	return r
}
func (c *ExperimentCollaboration) Publish(m *ExperimentCollaborationMessage) (r *ExperimentCollaborationMessage) {
	r = &ExperimentCollaborationMessage{}
	r.Type = "message"
	r.MessagePath = m.PublishPath
	r.MessageBody = m.PublishBody
	// switch m.PublishPath {
	// case "answers":
	// 	s, ok := m.PublishBody["content"].(string)
	// 	if !ok {
	// 		log.Printf("Failed get content %+v", m.PublishBody)
	// 		break
	// 	}
	// 	c.e.Hypotesis = []string{s}
	// }
	return r
}
