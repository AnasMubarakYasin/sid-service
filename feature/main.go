package feature

import "sid/service/repository"

type Core struct {
	User       *User
	Profile    *Profile
	Subject    *Subject
	Experiment *Experiment
	Resource   *Resource
}

func New(repository *repository.Core) *Core {
	usr := NewUser(repository.User, repository.Profile)
	pro := NewProfile(repository.Profile)
	sub := NewSubject(repository.Subject)
	exp := NewExperiment(repository.Experiment, repository.Subject, repository.User)

	return &Core{
		User:       usr,
		Profile:    pro,
		Subject:    sub,
		Experiment: exp,
		Resource:   NewResource(usr, pro, exp, sub),
	}
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}
