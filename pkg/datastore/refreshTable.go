package datastore

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	msSQL = "mssql"
	mySQL = "mysql"
	pgSQL = "postgres"
)

type Seeder struct {
	*DataStore

	path         string
	dialect      string
	ResetCounter bool
}

func NewSeeder(db *DataStore, directoryPath string) *Seeder {
	v := db.GORM()
	dialect := ""

	if v != nil {
		dialect = db.GORM().Dialector.Name()
	}

	return &Seeder{DataStore: db, path: directoryPath, dialect: dialect}
}

/* RefreshTables : The function will first clear the tables and then populate it with the test data for each table.
The tables will have to be passed in the reverse order in which the dependency flows,i.e, the child first and then the parent */
func (d *Seeder) RefreshTables(t tester, tableNames ...string) {
	for _, tableName := range tableNames {
		d.ClearTable(t, tableName)
	}

	for index := len(tableNames) - 1; index >= 0; index-- {
		tableName := tableNames[index]

		records, err := d.getRecords(tableName)
		if err != nil {
			t.Error(err)
			return
		}

		d.populateTable(t, tableName, records)
	}
}

func (d *Seeder) ClearTable(t tester, tableName string) {
	_, err := d.DB().Exec(`DELETE` + ` FROM ` + tableName)
	if err != nil {
		t.Error(err)
		return
	}
}

func (d *Seeder) populateTable(t tester, tableName string, records [][]string) {
	d.resetIdentitySequence(t, tableName, true)
	txn := d.GORM().Begin()

	var err error

	// this indicates if a table has identity column or not
	identityInsert := false

	if d.dialect == "sqlserver" {
		identityInsert, err = getIdentityInsert(txn, tableName)
		if err != nil {
			_ = txn.Rollback()

			t.Error(err)

			return
		}
	}

	query := d.getQueryFromRecords(records, tableName)

	err = txn.Exec(query).Error
	if err != nil {
		_ = txn.Rollback()

		t.Error(err)

		return
	}

	if d.dialect == msSQL && identityInsert {
		err = txn.Exec(`SET ` + `IDENTITY_INSERT ` + tableName + ` OFF`).Error
		if err != nil {
			_ = txn.Rollback()

			t.Error(err)

			return
		}
	}

	_ = txn.Commit()

	// identity sequence has to be set only after test data has been added in case of postgres
	d.resetIdentitySequence(t, tableName, false)
}

// resets identity in case of mssql and sequence in case of postgres
func (d *Seeder) resetIdentitySequence(t tester, tableName string, beforeTransaction bool) {
	if !d.ResetCounter {
		return
	}

	var q string
	// in case of mysql and mssql, resetting identity to 0 at beginning works but in case of pgsql, this has to be done
	// after the data has been inserted
	switch beforeTransaction {
	case false:
		if d.dialect == pgSQL {
			//nolint
			q = `SELECT pg_catalog.setval(pg_get_serial_sequence('` + tableName + `', 'id'), (SELECT MAX(id) FROM ` + tableName + `));`
		}
	default:
		if d.dialect == mySQL {
			q = `ALTER TABLE ` + tableName + ` AUTO_INCREMENT = 0;`
		}

		if d.dialect == msSQL {
			q = `DBCC CHECKIDENT (` + tableName + `, RESEED, 0)`
		}
	}

	if err := d.GORM().Exec(q).Error; err != nil {
		t.Errorf("Unable to reset identity. got err: %v", err)
	}
}

// getIdentityInsert checks if the MSSQL table has an identity column, if yes, it will turn IDENTITY_INSERT to ON in order to insert
// values to the identity columns
func getIdentityInsert(txn *gorm.DB, tableName string) (bool, error) {
	var name string

	// query the information schema to identify if the tables has an identity
	_ = txn.Raw(`SELECT TABLE_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE 
		COLUMNPROPERTY(object_id(TABLE_SCHEMA+'.'+TABLE_NAME), COLUMN_NAME, 'IsIdentity') = 1 ORDER BY TABLE_NAME`).Scan(&name)

	identityInsert := false

	if name == tableName {
		identityInsert = true
	}

	if identityInsert {
		err := txn.Exec(`SET` + ` IDENTITY_INSERT ` + tableName + ` ON`).Error

		if err != nil {
			return identityInsert, err
		}
	}

	return identityInsert, nil
}

func (d *Seeder) getRecords(tableName string) ([][]string, error) {
	fileLocation := d.path + "/" + tableName + ".csv"

	fileLoc, err := os.Open(fileLocation)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(fileLoc)

	return reader.ReadAll()
}

func (d *Seeder) getQueryFromRecords(records [][]string, tableName string) string {
	columns := records[0]
	query := "insert into " + tableName + " (" + strings.Join(columns, ",") + ") values"

	var values []string

	for i := 1; i < len(records); i++ {
		var rows []string

		for j := range records[i] {
			if !strings.EqualFold(records[i][j], "NULL") {
				rows = append(rows, "'"+records[i][j]+"'")
			} else {
				rows = append(rows, records[i][j])
			}
		}

		values = append(values, "("+strings.Join(rows, ",")+")")
	}

	query += strings.Join(values, ",")

	return query
}

func (d *Seeder) getCassandraQueryFromRecords(records [][]string, tableName string) (query string, rows []interface{}, err error) {
	columns := records[0]
	columnsStr := " (\"" + strings.Join(columns, "\",\"") + "\")"

	qRows := make([]string, len(columns))
	for i := range records[0] {
		qRows[i] = "?"
	}

	qRowsStr := strings.Join(qRows, ",")
	query = "BEGIN BATCH"

	schema, err := d.Cassandra.Session.KeyspaceMetadata(d.Cassandra.config.Keyspace)
	if err != nil {
		return "", nil, err
	}
	// types contains mapping of columns to their types
	types := make(map[string]string)
	for _, col := range schema.Tables[tableName].Columns {
		types[col.Name] = col.Validator
	}

	rows, err = marshalRecords(records, types)
	if err != nil {
		return "", nil, err
	}

	for i := 0; i < len(records)-1; i++ {
		query += " insert into " + tableName + columnsStr + " values(" + qRowsStr + ");"
	}

	query += " APPLY BATCH"

	return query, rows, nil
}

// nolint:gocyclo,gocognit // cannot break down the function further
func marshalRecords(records [][]string, types map[string]string) ([]interface{}, error) {
	columns := records[0]
	rows := make([]interface{}, len(columns)*(len(records)-1))
	// type casting all values by columns
	for col := range columns {
		switch types[columns[col]] {
		case "double":
			const bitSize = 64
			// marshaling whole row for double type
			for row := 1; row < len(records); row++ {
				f, err := strconv.ParseFloat(records[row][col], bitSize)
				if err != nil {
					return nil, err
				}

				rows[col+(row-1)*(len(columns))] = f
			}
		case "timestamp":
			const layout = "2006-01-02 15:04:05"

			for row := 1; row < len(records); row++ {
				t, err := time.Parse(layout, records[row][col])
				if err != nil {
					return nil, err
				}

				rows[col+(row-1)*(len(columns))] = t
			}
		case "time":
			for row := 1; row < len(records); row++ {
				t, err := time.ParseDuration(records[row][col])
				if err != nil {
					return nil, err
				}

				rows[col+(row-1)*(len(columns))] = t
			}
		case "boolean":
			for row := 1; row < len(records); row++ {
				t, err := strconv.ParseBool(records[row][col])
				if err != nil {
					return nil, err
				}

				rows[col+(row-1)*(len(columns))] = t
			}
		default:
			// marshaling whole row for other types
			for row := 1; row < len(records); row++ {
				rows[col+(row-1)*(len(columns))] = records[row][col]
			}
		}
	}

	return rows, nil
}

// check the type is int or float type or not
func check(s string) bool {
	if s == "int" || s == "float" {
		return true
	}

	return false
}

func (d *Seeder) getYCQLQueryFromRecords(records [][]string, tableName string) string {
	columns := records[0]
	n := len(columns)
	columnsStr := " (\"" + strings.Join(columns, "\",\"") + "\")"

	query := "BEGIN TRANSACTION  "
	insertStmt := "insert into  " + tableName + columnsStr + " VALUES"

	fieldTypes := make([]string, n)
	i := 0

	field := ""

	keyspace := d.YCQL.Cluster.Keyspace

	iter := d.YCQL.Session.Query("SELECT   type  FROM system_schema.columns WHERE  keyspace_name =?"+
		" AND table_name = ?; ", keyspace, tableName).Iter()

	// through this we can get field type type of table so, that accordingly we implement query
	for iter.Scan(&field) {
		fieldTypes[i] = field
		i++
	}

	for i := 1; i < len(records); i++ {
		var rows []string

		for j := range records[i] {
			// check the field type is int or not
			if check(fieldTypes[j]) {
				rows = append(rows, records[i][j])
			} else {
				rows = append(rows, "'"+records[i][j]+"'")
			}
		}

		query += insertStmt + "(" + strings.Join(rows, ",") + ");"
	}

	query += " END TRANSACTION ;"

	return query
}

func (d *Seeder) AssertVersion(t tester, version string) {
	var ver, query string

	switch d.dialect {
	case mySQL:
		query = "SELECT @@version as version"

	case pgSQL:
		query = "SHOW server_version"

	case msSQL:
		query = "SELECT @@MICROSOFTVERSION / 0x01000000 AS MajorVersionNumber"
	}

	err := d.DB().QueryRow(query).Scan(&ver)
	if err != nil {
		t.Error(err)
	}

	if version != ver {
		t.Errorf("Version Mismatch. Required Version: %s. Version in use: %s", version, ver)
		return
	}
}

func (d *Seeder) AssertRowCount(t tester, tableName string, count int) {
	var ct int

	query := `SELECT COUNT(*)` + `FROM ` + tableName

	err := d.DB().QueryRow(query).Scan(&ct)
	if err != nil {
		t.Error(err)
	}

	if ct != count {
		t.Errorf("incorrect number of records. expected: %d got: %d", count, ct)
		return
	}
}

func (d *Seeder) RefreshMongoCollections(t tester, collectionNames ...string) {
	for i := range collectionNames {
		collectionName := collectionNames[i]
		fileLoc := d.path + "/" + collectionName + ".json"

		file, err := os.ReadFile(fileLoc)
		if err != nil {
			t.Error(err)
			return
		}

		var data []interface{}

		err = json.Unmarshal(file, &data)
		if err != nil {
			t.Error(err)
			return
		}

		collection := d.MongoDB.Collection(collectionName)

		err = collection.Drop(context.TODO())
		if err != nil {
			t.Error(err)
			return
		}

		_, err = collection.InsertMany(context.TODO(), data)
		if err != nil {
			t.Error(err)
		}
	}
}

func (d *Seeder) RefreshCassandra(t tester, tableNames ...string) {
	for i := range tableNames {
		tableName := tableNames[i]

		err := d.Cassandra.Session.Query(`TRUNCATE ` + tableName).Exec()
		if err != nil {
			t.Error(err)
			return
		}

		records, err := d.getRecords(tableName)
		if err != nil {
			t.Error(err)
			return
		}

		query, rows, err := d.getCassandraQueryFromRecords(records, tableName)
		if err != nil {
			t.Error(err)
			return
		}

		err = d.Cassandra.Session.Query(query, rows...).Exec()
		if err != nil {
			t.Error(err)
			return
		}
	}
}

func (d *Seeder) RefreshYCQL(t tester, tableNames ...string) {
	for i := range tableNames {
		tableName := tableNames[i]

		err := d.YCQL.Session.Query(`TRUNCATE ` + tableName).Exec()
		if err != nil {
			t.Error(err)
			return
		}

		records, err := d.getRecords(tableName)
		if err != nil {
			t.Error(err)
			return
		}

		q := d.getYCQLQueryFromRecords(records, tableName)

		err = d.YCQL.Session.Query(q).Exec()
		if err != nil {
			t.Error(err)
		}
	}
}

// nolint:gocognit // cannot break down the function further
func (d *Seeder) RefreshRedis(t tester, tableNames ...string) {
	for i := range tableNames {
		tableName := tableNames[i]

		records, err := d.getRecords(tableName)
		if err != nil {
			// if <tableName>.csv not found then looking for <tableName>.json
			err = d.setRedisHashMaps(tableName)
			if err != nil {
				t.Error(err)
			}

			return
		}

		const recordLimit = 2

		for r := range records {
			if len(records[r]) != recordLimit {
				t.Error("The csv input for redis should have data in the format - key,value")
				return
			}

			d.Redis.Set(context.Background(), records[r][0], records[r][1], 0)
		}

		_ = d.setRedisHashMaps(tableName)
	}
}

func (d *Seeder) setRedisHashMaps(tableName string) error {
	fileLoc := d.path + "/" + tableName + ".json"

	file, err := os.ReadFile(fileLoc)
	if err != nil {
		return err
	}

	var data []map[string]interface{}

	err = json.Unmarshal(file, &data)
	if err != nil {
		return err
	}

	keys, err := d.Redis.Keys(context.Background(), tableName+":*").Result()
	if err != nil {
		return err
	}

	d.Redis.Del(context.Background(), keys...)

	for i := range data {
		hKey := tableName + ":" + strconv.Itoa(i)

		for k, v := range data[i] {
			d.Redis.HSet(context.Background(), hKey, k, v)
		}
	}

	return nil
}

func (d *Seeder) RefreshDynamoDB(t tester, tableNames ...string) {
	for _, tableName := range tableNames {
		fileLoc := fmt.Sprintf("%s/%s.json", d.path, tableName)

		raw, err := ioutil.ReadFile(fileLoc)
		if err != nil {
			t.Errorf("Got error reading file: %s", err)

			return
		}

		putItem(d, t, tableName, raw)
	}
}

func putItem(d *Seeder, t tester, tableName string, raw []byte) {
	var items []map[string]interface{}

	err := json.Unmarshal(raw, &items)
	if err != nil {
		t.Error(err)

		return
	}

	for _, item := range items {
		av, err := dynamodbattribute.MarshalMap(item)
		if err != nil {
			t.Errorf("Got error while marshaling map: %s", err)
		}

		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(tableName),
		}

		_, err = d.DynamoDB.PutItem(input)
		if err != nil {
			t.Errorf("Got error while calling PutItem: %s", err)
		}
	}
}
