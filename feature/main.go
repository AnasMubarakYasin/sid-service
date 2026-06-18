package feature

import "sid/service/repository"

type Core struct {
	User       *User
	Profile    *Profile
	Subject    *Subject
	Experiment *Experiment
}

func New(repository *repository.Core) *Core {
	return &Core{
		User:       NewUser(repository.User, repository.Profile),
		Profile:    NewProfile(repository.Profile),
		Subject:    NewSubject(repository.Subject),
		Experiment: NewExperiment(repository.Experiment, repository.Subject, repository.User),
	}
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}
