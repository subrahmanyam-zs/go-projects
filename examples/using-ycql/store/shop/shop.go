package shop

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/zopsmart/gofr/examples/using-ycql/entity"
	"github.com/zopsmart/gofr/pkg/errors"
	"github.com/zopsmart/gofr/pkg/gofr"
)

// store function
type Shop struct{}

func (s Shop) Get(ctx *gofr.Context, filter entity.Shop) []entity.Shop {
	var (
		shop  entity.Shop
		shops []entity.Shop
	)

	whereCL, values := getWhereClause(filter)
	//nolint:gosec // for sql string concatenation
	query := ` select id, name, location ,state FROM shop ` + whereCL
	iter := ctx.YCQL.Session.Query(query, values...).Iter()

	for iter.Scan(&shop.ID, &shop.Name, &shop.Location, &shop.State) {
		shops = append(shops, entity.Shop{ID: shop.ID, Name: shop.Name, Location: shop.Location, State: shop.State})
	}

	return shops
}

func (s Shop) Create(ctx *gofr.Context, data entity.Shop) ([]entity.Shop, error) {
	query := "INSERT INTO shop (id, name, location, state) VALUES (?, ?, ?, ?)"

	err := ctx.YCQL.Session.Query(query, data.ID, data.Name, data.Location, data.State).Exec()
	if err != nil {
		return nil, errors.DB{Err: err}
	}

	return s.Get(ctx, entity.Shop{ID: data.ID}), nil
}

func (s Shop) Delete(ctx *gofr.Context, id string) error {
	query := "DELETE FROM shop  WHERE id = ?"

	err := ctx.YCQL.Session.Query(query, id).Exec()
	if err != nil {
		return errors.DB{Err: err}
	}

	return err
}

func (s Shop) Update(ctx *gofr.Context, data entity.Shop) ([]entity.Shop, error) {
	query := "UPDATE shop"
	setCl, values := genSetClause(&data)

	// No value is passed for update
	if values == nil {
		return s.Get(ctx, entity.Shop{ID: data.ID}), nil
	}

	query = fmt.Sprintf("%v %v where id = ?", query, setCl)
	id := strconv.Itoa(data.ID)

	values = append(values, id)

	err := ctx.YCQL.Session.Query(query, values...).Exec()
	if err != nil {
		return nil, errors.DB{Err: err}
	}

	return s.Get(ctx, entity.Shop{ID: data.ID}), nil
}

func genSetClause(p *entity.Shop) (setClause string, values []interface{}) {
	setClause = `SET`

	if p.Name != "" {
		setClause += " name = ?,"

		values = append(values, p.Name)
	}

	if p.Location != "" {
		setClause += " location = ?,"

		values = append(values, p.Location)
	}

	if p.State != "" {
		setClause += " state = ?,"

		values = append(values, p.State)
	}

	if setClause == "SET" {
		return "", nil
	}

	setClause = strings.TrimSuffix(setClause, ",")

	return setClause, values
}

func getWhereClause(p entity.Shop) (whereCl string, values []interface{}) {
	conditions := make([]string, 0)

	if p.ID != 0 {
		conditions = append(conditions, "id = ?")
		values = append(values, p.ID)
	}

	if p.Name != "" {
		conditions = append(conditions, "name = ?")
		values = append(values, p.Name)
	}

	if p.Location != "" {
		conditions = append(conditions, "location = ?")
		values = append(values, p.Location)
	}

	if p.State != "" {
		conditions = append(conditions, "state = ?")
		values = append(values, p.State)
	}

	if len(conditions) > 0 {
		whereCl = " where " + strings.Join(conditions, " AND ") + " ALLOW FILTERING"
	}

	return whereCl, values
}
