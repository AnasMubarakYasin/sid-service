package handler

import (
	"sid/service/feature"

	"github.com/gin-gonic/gin"
)

type Subject struct {
	f *feature.Subject
}

func NewSubject(feature *feature.Subject) *Subject {
	return &Subject{f: feature}
}

func (h *Subject) Create(c *gin.Context) {
	p := &feature.ParamCreateSubject{}

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

func (h *Subject) Update(c *gin.Context) {
	p := &ParamID{}

	if err := c.ShouldBindUri(p); err != nil {
		_ = c.Error(ErrBadRequest)
		return
	}

	pp := &feature.ParamUpdateSubject{}

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

func (h *Subject) Delete(c *gin.Context) {
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

func (h *Subject) List(c *gin.Context) {
	r, e := h.f.List(c)

	if e != nil {
		Fail(c, 500, e)
		return
	}

	OKM(c, r, &Meta{})
}

func (h *Subject) Get(c *gin.Context) {
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
