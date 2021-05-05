package person

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/zopsmart/gofr/examples/using-cassandra/entity"
	"github.com/zopsmart/gofr/pkg/errors"
	"github.com/zopsmart/gofr/pkg/gofr"
)

// store function

type Person struct{}

func (p Person) Get(ctx *gofr.Context, filter entity.Person) []*entity.Person {
	var (
		id    int
		name  string
		age   int
		state string
	)

	var persons []*entity.Person = nil

	cassDB := ctx.Cassandra.Session
	whereCL, values := getWhereClause(filter)
	query := `SELECT id, name, age ,state FROM persons`
	querystring := fmt.Sprintf("%s %s", query, whereCL)
	iter := cassDB.Query(querystring, values...).Iter()

	for iter.Scan(&id, &name, &age, &state) {
		persons = append(persons, &entity.Person{ID: id, Name: name, Age: age, State: state})
	}

	return persons
}

func (p Person) Create(ctx *gofr.Context, data entity.Person) ([]*entity.Person, error) {
	cassDB := ctx.Cassandra.Session
	query := "INSERT INTO persons (id, name, age, state) VALUES (?, ?, ?, ?)"

	err := cassDB.Query(query, data.ID, data.Name, data.Age, data.State).Exec()
	if err != nil {
		return nil, errors.DB{Err: err}
	}

	return p.Get(ctx, entity.Person{ID: data.ID}), nil
}

func (p Person) Delete(ctx *gofr.Context, id string) error {
	cassDB := ctx.Cassandra.Session
	query := "DELETE FROM persons WHERE id = ?"

	err := cassDB.Query(query, id).Exec()
	if err != nil {
		return errors.DB{Err: err}
	}

	return err
}

func (p Person) Update(ctx *gofr.Context, data entity.Person) ([]*entity.Person, error) {
	cassDB := ctx.Cassandra.Session
	query := "UPDATE persons"
	setCl, values := genSetClause(&data)

	// No value is passed for update
	if values == nil {
		return p.Get(ctx, entity.Person{ID: data.ID}), nil
	}

	query = fmt.Sprintf("%v %v WHERE id = ?", query, setCl)
	id := strconv.Itoa(data.ID)

	values = append(values, id)

	err := cassDB.Query(query, values...).Exec()
	if err != nil {
		return nil, errors.DB{Err: err}
	}

	return p.Get(ctx, entity.Person{ID: data.ID}), nil
}

func genSetClause(p *entity.Person) (setClause string, values []interface{}) {
	setClause = `SET`

	if p.Name != "" {
		setClause += " name = ?,"

		values = append(values, p.Name)
	}

	if p.Age > 0 {
		setClause += " age = ?,"

		values = append(values, p.Age)
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

func getWhereClause(p entity.Person) (whereCl string, values []interface{}) {
	whereCl = " WHERE "

	if p.ID != 0 {
		whereCl += "id = ? AND "

		values = append(values, p.ID)
	}

	if p.Name != "" {
		whereCl += "name = ? AND "

		values = append(values, p.Name)
	}

	if p.Age != 0 {
		whereCl += "age = ? AND "

		values = append(values, p.Age)
	}

	if p.State != "" {
		whereCl += "state = ? AND "

		values = append(values, p.State)
	}

	whereCl = strings.TrimSuffix(whereCl, "AND ")
	whereCl = strings.TrimSuffix(whereCl, " WHERE ")

	if whereCl != "" {
		whereCl += " ALLOW FILTERING "
	}

	return whereCl, values
}
