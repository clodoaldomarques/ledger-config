package dynamodb

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/clodoaldomarques/ledger-config/internal/domain/ledger"
)

type Repository struct {
	client    *dynamodb.Client
	tableName string
}

func NewRepository() *Repository {
	client := dynamodb.NewFromConfig(configure())
	return &Repository{
		client:    client,
		tableName: "ConfigTable",
	}
}

func (r Repository) SaveScript(ctx context.Context, s ledger.Config) error {
	st := buildConfigTable(s)

	item, err := attributevalue.MarshalMap(st)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}

	_, err = r.client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put item in dynamodb: %w", err)
	}

	return nil
}

func (r Repository) UpdateScript(ctx context.Context, s ledger.Config) error {
	st := buildConfigTable(s)

	itemMap, err := attributevalue.MarshalMap(st)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	delete(itemMap, "org_id")
	delete(itemMap, "config_id")
	delete(itemMap, "filters")
	delete(itemMap, "updated_at")

	updateExpression := "SET"
	exprAttrValues := make(map[string]types.AttributeValue)
	exprAttrNames := make(map[string]string)

	i := 1
	for field, value := range itemMap {
		valPlaceholder := fmt.Sprintf(":val%d", i)
		namePlaceholder := fmt.Sprintf("#field%d", i)

		updateExpression += fmt.Sprintf(" %s = %s,", namePlaceholder, valPlaceholder)
		exprAttrNames[namePlaceholder] = field
		exprAttrValues[valPlaceholder] = value
		i++
	}

	// Adicionar timestamp de atualização
	updateExpression += " updated_at = :now"
	exprAttrValues[":now"] = &types.AttributeValueMemberS{
		Value: time.Now().UTC().Format(time.RFC3339Nano),
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"org_id":    &types.AttributeValueMemberS{Value: s.OrgID},
			"script_id": &types.AttributeValueMemberS{Value: s.ConfigID},
		},
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeNames:  exprAttrNames,
		ExpressionAttributeValues: exprAttrValues,
		ReturnValues:              types.ReturnValueAllNew,
	}
	_, err = r.client.UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put item in dynamodb: %w", err)
	}

	return nil
}

func (r Repository) FindScriptByID(ctx context.Context, orgID string, scriptID string) (ledger.Config, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"org_id":    &types.AttributeValueMemberS{Value: orgID},
			"script_id": &types.AttributeValueMemberS{Value: scriptID},
		},
	}

	result, err := r.client.GetItem(ctx, input)
	if err != nil {
		return ledger.Config{}, fmt.Errorf("failed to get item: %w", err)
	}

	if result.Item == nil {
		return ledger.Config{}, ErrScriptNotFound{}
	}

	var script Config
	if err := attributevalue.UnmarshalMap(result.Item, &script); err != nil {
		return ledger.Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return script.toEntity(), nil
}

func (r Repository) FindScriptByLevel(ctx context.Context, level string, eventTypeID string, orgID string, programID *int64) (ledger.Config, error) {
	filters := buildFilters(level, orgID, eventTypeID, *programID)
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("GSI-Index"),
		KeyConditionExpression: aws.String("org_id = :ORGID AND filters = :FILTERS"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":ORGID":   &types.AttributeValueMemberS{Value: orgID},
			":FILTERS": &types.AttributeValueMemberS{Value: filters},
		},
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int32(1),
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return ledger.Config{}, fmt.Errorf("failed to get item: %w", err)
	}

	if len(result.Items) == 0 {
		return ledger.Config{}, ErrScriptNotFound{}
	}

	var script Config
	if err := attributevalue.UnmarshalMap(result.Items[0], &script); err != nil {
		return ledger.Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return script.toEntity(), nil
}

func (r Repository) FindAllScripts(ctx context.Context, orgID string, programID *int64) ([]ledger.Config, error) {
	input := buildInputQuery(r.tableName, orgID)
	if programID != nil {
		input = buildInputQueryWithFilters(r.tableName, orgID, *programID)
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, ErrScriptNotFound{}
	}

	scripts := make([]ledger.Config, 0, len(result.Items))

	for _, m := range result.Items {
		var script Config
		if err := attributevalue.UnmarshalMap(m, &script); err != nil {
			return nil, fmt.Errorf("failed to unmarshal script: %w", err)
		}
		scripts = append(scripts, script.toEntity())
	}

	return scripts, nil
}

func buildInputQuery(tableName string, orgID string) *dynamodb.QueryInput {
	return &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		IndexName:              aws.String("GSI-Index"),
		KeyConditionExpression: aws.String("org_id = :ORG_ID"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":ORG_ID": &types.AttributeValueMemberS{Value: orgID},
		},
		ScanIndexForward: aws.Bool(false),
	}
}

func buildInputQueryWithFilters(tableName string, orgID string, programID int64) *dynamodb.QueryInput {
	filters := buildAllQuery(string(ledger.ProgramLevel), orgID, &programID)
	return &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		IndexName:              aws.String("GSI-Index"),
		KeyConditionExpression: aws.String("org_id = :ORG_ID AND begins_with(filters, :event_prefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":ORG_ID":       &types.AttributeValueMemberS{Value: orgID},
			":event_prefix": &types.AttributeValueMemberS{Value: filters},
		},
		ScanIndexForward: aws.Bool(false),
	}
}

func (r Repository) Close() {

}
