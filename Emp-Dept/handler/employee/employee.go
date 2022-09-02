package employee

import (
	"fmt"
	"strconv"

	"github.com/google/uuid"

	"developer.zopsmart.com/go/gofr/Emp-Dept/entities"
	"developer.zopsmart.com/go/gofr/Emp-Dept/service"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Handler struct {
	service service.Employee
}

func New(s service.Employee) Handler {
	return Handler{service: s}
}

func (h Handler) Post(ctx *gofr.Context) (interface{}, error) {
	var emp entities.Employee

	err := ctx.Bind(&emp)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"invalid parameters"}}
	}

	res, err := h.service.Post(ctx, emp)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h Handler) Put(ctx *gofr.Context) (interface{}, error) {
	var dataToUpdate entities.Employee

	id := ctx.PathParam("id")

	empID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"invalid id"}}
	}

	err = ctx.Bind(&dataToUpdate)
	if err != nil {
		fmt.Println(err)
		return nil, errors.InvalidParam{Param: []string{"invalid details"}}
	}

	res, err := h.service.Put(ctx, empID, dataToUpdate)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h Handler) Delete(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")

	empID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"invalid id"}}
	}

	res, err := h.service.Delete(ctx, empID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h Handler) Get(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")

	empID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"invalid id"}}
	}

	res, err := h.service.Get(ctx, empID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h Handler) GetAll(ctx *gofr.Context) (interface{}, error) {
	name := ctx.Param("name")
	include := ctx.Param("includeDepartment")

	includeDepartment, err := strconv.ParseBool(include)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"unconvertable type to bool"}}
	}

	res, err := h.service.GetAll(ctx, name, includeDepartment)
	if err != nil {
		return nil, err
	}

	return res, err
}
