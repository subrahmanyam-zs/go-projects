package datastore

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/net/context"

	"developer.zopsmart.com/go/gofr/pkg"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
	"developer.zopsmart.com/go/gofr/pkg/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/prometheus/client_golang/prometheus"
)

// nolint:gochecknoglobals // dynamodbStats has to be a global variable for prometheus
var (
	dynamodbStats = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "zs_dynamodb_stats",
		Help:    "Histogram for DynamoDB",
		Buckets: []float64{.001, .003, .005, .01, .025, .05, .1, .2, .3, .4, .5, .75, 1, 2, 3, 5, 10, 30},
	}, []string{"type", "host", "table"})

	_ = prometheus.Register(dynamodbStats)
)

// DynamoDBConfig configuration for DynamoDB connection
type DynamoDBConfig struct {
	Region            string
	Endpoint          string
	AccessKeyID       string
	SecretAccessKey   string
	ConnRetryDuration int
}

type DynamoDB struct {
	*dynamodb.DynamoDB
	logger log.Logger
	config DynamoDBConfig
}

// NewDynamoDB connects to DynamoDB and returns the connection
func NewDynamoDB(logger log.Logger, c DynamoDBConfig) (DynamoDB, error) {
	sessionConfig := &aws.Config{
		Region:      aws.String(c.Region),
		Logger:      logger,
		Endpoint:    aws.String(c.Endpoint),
		Credentials: credentials.NewStaticCredentials(c.AccessKeyID, c.SecretAccessKey, ""),
	}

	sess, err := session.NewSession(sessionConfig)
	if err != nil {
		return DynamoDB{}, err
	}

	db := dynamodb.New(sess)

	// check the db connection
	_, err = db.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		return DynamoDB{}, err
	}

	return DynamoDB{DynamoDB: db, logger: logger, config: c}, nil
}

// HealthCheck checks health of the Dya
func (d *DynamoDB) HealthCheck() types.Health {
	resp := types.Health{
		Name:   pkg.DynamoDB,
		Status: pkg.StatusDown,
	}

	// check if DynamoDB instance has been created during initialization
	if d.DynamoDB == nil {
		return resp
	}

	_, err := d.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		return resp
	}

	resp.Status = pkg.StatusUp

	return resp
}

func (d *DynamoDB) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	begin := time.Now()

	out, err := d.DynamoDB.PutItem(input)

	duration := time.Since(begin)
	query := genPutItemQuery(input)

	d.monitorQuery(query, begin, duration)

	return out, err
}

func (d *DynamoDB) PutItemRequest(input *dynamodb.PutItemInput) (*request.Request, *dynamodb.PutItemOutput) {
	begin := time.Now()

	req, out := d.DynamoDB.PutItemRequest(input)

	duration := time.Since(begin)
	query := genPutItemQuery(input)

	d.monitorQuery(query, begin, duration)

	return req, out
}

func (d *DynamoDB) PutItemWithContext(ctx context.Context, input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	begin := time.Now()

	out, err := d.DynamoDB.PutItemWithContext(ctx, input)

	duration := time.Since(begin)
	query := genPutItemQuery(input)

	d.monitorQuery(query, begin, duration)

	return out, err
}

func (d *DynamoDB) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	begin := time.Now()

	out, err := d.DynamoDB.GetItem(input)

	duration := time.Since(begin)
	query := genGetItemQuery(input)

	d.monitorQuery(query, begin, duration)

	return out, err
}

func (d *DynamoDB) GetItemRequest(input *dynamodb.GetItemInput) (*request.Request, *dynamodb.GetItemOutput) {
	begin := time.Now()

	req, out := d.DynamoDB.GetItemRequest(input)

	duration := time.Since(begin)
	query := genGetItemQuery(input)

	d.monitorQuery(query, begin, duration)

	return req, out
}

func (d *DynamoDB) GetItemWithContext(ctx context.Context, input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	begin := time.Now()

	out, err := d.DynamoDB.GetItemWithContext(ctx, input)

	duration := time.Since(begin)
	query := genGetItemQuery(input)

	d.monitorQuery(query, begin, duration)

	return out, err
}

func (d *DynamoDB) DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	begin := time.Now()

	out, err := d.DynamoDB.DeleteItem(input)

	duration := time.Since(begin)
	query := genDeleteItemQuery(input)

	d.monitorQuery(query, begin, duration)

	return out, err
}

func (d *DynamoDB) DeleteItemRequest(input *dynamodb.DeleteItemInput) (*request.Request, *dynamodb.DeleteItemOutput) {
	begin := time.Now()

	req, out := d.DynamoDB.DeleteItemRequest(input)

	duration := time.Since(begin)
	query := genDeleteItemQuery(input)

	d.monitorQuery(query, begin, duration)

	return req, out
}

func (d *DynamoDB) DeleteItemWithContext(ctx context.Context, input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	begin := time.Now()

	out, err := d.DynamoDB.DeleteItemWithContext(ctx, input)

	duration := time.Since(begin)
	query := genDeleteItemQuery(input)

	d.monitorQuery(query, begin, duration)

	return out, err
}

func (d *DynamoDB) UpdateItem(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	begin := time.Now()

	out, err := d.DynamoDB.UpdateItem(input)

	duration := time.Since(begin)
	query := genUpdateItemQuery(input)

	d.monitorQuery(query, begin, duration)

	return out, err
}

func (d *DynamoDB) UpdateItemRequest(input *dynamodb.UpdateItemInput) (*request.Request, *dynamodb.UpdateItemOutput) {
	begin := time.Now()

	req, out := d.DynamoDB.UpdateItemRequest(input)

	duration := time.Since(begin)
	query := genUpdateItemQuery(input)

	d.monitorQuery(query, begin, duration)

	return req, out
}

func (d *DynamoDB) UpdateItemWithContext(ctx context.Context, input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	begin := time.Now()

	out, err := d.DynamoDB.UpdateItemWithContext(ctx, input)

	duration := time.Since(begin)
	query := genUpdateItemQuery(input)

	d.monitorQuery(query, begin, duration)

	return out, err
}

func (d *DynamoDB) monitorQuery(query []string, begin time.Time, duration time.Duration) {
	var (
		ql    QueryLogger
		table string
	)

	queryOperation := query[0]
	lenQuery := len(query)

	table = query[lenQuery-1]
	query[lenQuery-1] = fmt.Sprintf("on table %v", table)

	ql.Query = make([]string, 1)

	if lenQuery > 1 {
		ql.Query[0] = fmt.Sprintf("%v - with %v", query[0], strings.Join(query[1:], ", "))
	}

	ql.Duration = duration.Microseconds()
	ql.StartTime = begin
	ql.DataStore = pkg.DynamoDB
	ql.Hosts = d.config.Endpoint

	// log the query

	if d.logger != nil {
		d.logger.Debug(ql)
	}

	// push stats to metrics server
	dynamodbStats.WithLabelValues(queryOperation, ql.Hosts, table).Observe(duration.Seconds())
}

func getAttributeNames(mp map[string]*dynamodb.AttributeValue) string {
	var names string

	for key := range mp {
		names += key + ", "
	}

	names = strings.TrimSuffix(names, ", ")

	return fmt.Sprintf("{%v}", names)
}

func getTableNameString(tableName *string) string {
	var name string

	if tableName != nil {
		name = *tableName
	}

	return name
}

func genPutItemQuery(input *dynamodb.PutItemInput) []string {
	query := []string{"PutItem"}

	query = append(query, fmt.Sprintf("Item Fields %v", getAttributeNames(input.Item)))

	if input.ConditionExpression != nil {
		query = append(query, fmt.Sprintf("ConditionExpression %v", *input.ConditionExpression))
	}

	query = append(query, getTableNameString(input.TableName))

	return query
}

func genGetItemQuery(input *dynamodb.GetItemInput) []string {
	query := []string{"GetItem"}

	if input.AttributesToGet != nil {
		var sub string

		for _, v := range input.AttributesToGet {
			sub += *v + ", "
		}

		sub = strings.TrimSuffix(sub, ", ")

		query = append(query, fmt.Sprintf("AttributesToGet {%v}", sub))
	}

	query = append(query, fmt.Sprintf("Key %v", getAttributeNames(input.Key)), getTableNameString(input.TableName))

	return query
}

func genDeleteItemQuery(input *dynamodb.DeleteItemInput) []string {
	query := []string{"DeleteItem"}

	if input.ConditionExpression != nil {
		query = append(query, fmt.Sprintf("ConditionExpression %v", *input.ConditionExpression))
	}

	query = append(query, fmt.Sprintf("Key %v", getAttributeNames(input.Key)), getTableNameString(input.TableName))

	return query
}

func genUpdateItemQuery(input *dynamodb.UpdateItemInput) []string {
	query := []string{"UpdateItem"}

	if input.AttributeUpdates != nil {
		var attributes string

		for key := range input.AttributeUpdates {
			attributes += key + ", "
		}

		attributes = strings.TrimSuffix(attributes, ", ")

		query = append(query, fmt.Sprintf("AttributesToUpdate {%v}", attributes))
	}

	if input.UpdateExpression != nil {
		query = append(query, fmt.Sprintf("UpdateExpression %v", *input.UpdateExpression))
	}

	if input.ConditionExpression != nil {
		query = append(query, fmt.Sprintf("ConditionExpression %v", *input.ConditionExpression))
	}

	query = append(query, fmt.Sprintf("Key %v", getAttributeNames(input.Key)), getTableNameString(input.TableName))

	return query
}
