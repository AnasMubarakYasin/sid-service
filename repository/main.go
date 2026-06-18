package repository

import "github.com/uptrace/bun"

type Core struct {
	User       *User
	Profile    *Profile
	Subject    *Subject
	Experiment *Experiment
}

func New(db *bun.DB) *Core {
	return &Core{
		User:       NewUser(db),
		Profile:    NewProfile(db),
		Subject:    NewSubject(db),
		Experiment: NewExperiment(db),
	}
}
