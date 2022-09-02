package department

import (
	"strconv"

	"developer.zopsmart.com/go/gofr/Emp-Dept/entities"
	"developer.zopsmart.com/go/gofr/Emp-Dept/service"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Handler struct {
	service service.Department
}

func New(s service.Department) Handler {
	return Handler{service: s}
}
func (h Handler) Post(ctx *gofr.Context) (interface{}, error) {
	var dept entities.Department

	err := ctx.Bind(&dept)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"invalid details"}}
	}

	res, err := h.service.Post(ctx, dept)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h Handler) Put(ctx *gofr.Context) (interface{}, error) {
	var dataToUpdate entities.Department

	id := ctx.Request().URL.Path[12:]

	deptID, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"invalid id"}}
	}

	err = ctx.Bind(&dataToUpdate)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"invalid details"}}
	}

	res, err := h.service.Put(ctx, deptID, dataToUpdate)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h Handler) Delete(ctx *gofr.Context) (interface{}, error) {
	id := ctx.Request().URL.Path[12:]

	deptID, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"Invalid id"}}
	}

	res, err := h.service.Delete(ctx, deptID)
	if err != nil {
		return nil, err
	}

	return res, nil
}
