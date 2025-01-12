package debug

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/spf13/cobra"
)

func newDynamoCommand() *cobra.Command {
	dynamoCmd := &cobra.Command{
		Use:   "dynamo",
		Short: "Debug DynamoDB operations",
		Long:  `Commands for debugging DynamoDB operations like listing, deleting, and canceling jobs.`,
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List jobs in DynamoDB",
		Long:  `List all jobs in DynamoDB with optional filtering by status and time range.`,
		RunE:  listJobs,
	}
	listCmd.Flags().StringP("status", "s", "", "Filter by status (e.g., FAILED, COMPLETED)")
	listCmd.Flags().StringP("since", "", "", "List jobs since timestamp (e.g., 2024-01-12)")
	listCmd.Flags().StringP("until", "", "", "List jobs until timestamp")

	deleteCmd := &cobra.Command{
		Use:   "delete [jobId]",
		Short: "Delete a job from DynamoDB",
		Long:  `Delete a specific job or multiple jobs from DynamoDB.`,
		RunE:  deleteJob,
	}
	deleteCmd.Flags().StringP("status", "s", "", "Delete all jobs with this status")
	deleteCmd.Flags().StringP("older-than", "", "", "Delete jobs older than timestamp")

	dynamoCmd.AddCommand(listCmd, deleteCmd)
	return dynamoCmd
}

func listJobs(cmd *cobra.Command, args []string) error {
	resources, err := LoadResources(cmd)
	if err != nil {
		return fmt.Errorf("failed to load resources: %w", err)
	}

	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(resources.Region),
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %w", err)
	}

	status, _ := cmd.Flags().GetString("status")
	since, _ := cmd.Flags().GetString("since")
	until, _ := cmd.Flags().GetString("until")

	var sinceTime, untilTime time.Time
	if since != "" {
		sinceTime, err = time.Parse("2006-01-02", since)
		if err != nil {
			return fmt.Errorf("invalid since date format: %w", err)
		}
	}
	if until != "" {
		untilTime, err = time.Parse("2006-01-02", until)
		if err != nil {
			return fmt.Errorf("invalid until date format: %w", err)
		}
	}

	input := &dynamodb.ScanInput{
		TableName: aws.String(resources.JobsTable),
	}

	var filterExpressions []string
	expressionValues := make(map[string]*dynamodb.AttributeValue)
	expressionNames := make(map[string]*string)

	if status != "" {
		filterExpressions = append(filterExpressions, "#status = :status")
		expressionValues[":status"] = &dynamodb.AttributeValue{S: aws.String(status)}
		expressionNames["#status"] = aws.String("Status")
	}

	if since != "" {
		filterExpressions = append(filterExpressions, "SubmittedAt >= :since")
		expressionValues[":since"] = &dynamodb.AttributeValue{S: aws.String(sinceTime.Format(time.RFC3339))}
	}

	if until != "" {
		filterExpressions = append(filterExpressions, "SubmittedAt <= :until")
		expressionValues[":until"] = &dynamodb.AttributeValue{S: aws.String(untilTime.Format(time.RFC3339))}
	}

	if len(filterExpressions) > 0 {
		input.FilterExpression = aws.String(strings.Join(filterExpressions, " AND "))
		input.ExpressionAttributeValues = expressionValues
		if len(expressionNames) > 0 {
			input.ExpressionAttributeNames = expressionNames
		}
	}

	// Create DynamoDB client
	client := dynamodb.New(sess)

	result, err := client.Scan(input)
	if err != nil {
		return fmt.Errorf("failed to scan DynamoDB: %w", err)
	}

	for _, item := range result.Items {
		var jobID, status, submittedAt, errorMsg string

		if v, ok := item["JobID"]; ok && v.S != nil {
			jobID = *v.S
		}
		if v, ok := item["Status"]; ok && v.S != nil {
			status = *v.S
		}
		if v, ok := item["SubmittedAt"]; ok && v.S != nil {
			submittedAt = *v.S
		}
		if v, ok := item["Error"]; ok && v.S != nil {
			errorMsg = *v.S
		}

		fmt.Printf("JobID: %s\n  Status: %s\n  SubmittedAt: %s\n", jobID, status, submittedAt)
		if errorMsg != "" {
			fmt.Printf("  Error: %s\n", errorMsg)
		}
		fmt.Println()
	}

	return nil
}

func deleteJob(cmd *cobra.Command, args []string) error {
	resources, err := LoadResources(cmd)
	if err != nil {
		return fmt.Errorf("failed to load resources: %w", err)
	}

	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(resources.Region),
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %w", err)
	}

	client := dynamodb.New(sess)

	status, _ := cmd.Flags().GetString("status")
	olderThan, _ := cmd.Flags().GetString("older-than")

	if len(args) > 0 {
		// Delete specific job
		jobID := args[0]
		input := &dynamodb.DeleteItemInput{
			TableName: aws.String(resources.JobsTable),
			Key: map[string]*dynamodb.AttributeValue{
				"JobID": {S: aws.String(jobID)},
			},
		}

		_, err := client.DeleteItem(input)
		if err != nil {
			return fmt.Errorf("failed to delete job %s: %w", jobID, err)
		}
		fmt.Printf("Successfully deleted job %s\n", jobID)
		return nil
	}

	// Batch delete based on filters
	input := &dynamodb.ScanInput{
		TableName: aws.String(resources.JobsTable),
	}

	var filterExpressions []string
	expressionValues := make(map[string]*dynamodb.AttributeValue)
	expressionNames := make(map[string]*string)

	if status != "" {
		filterExpressions = append(filterExpressions, "#status = :status")
		expressionValues[":status"] = &dynamodb.AttributeValue{S: aws.String(status)}
		expressionNames["#status"] = aws.String("Status")
	}

	if olderThan != "" {
		olderThanTime, err := time.Parse("2006-01-02", olderThan)
		if err != nil {
			return fmt.Errorf("invalid older-than date format: %w", err)
		}

		filterExpressions = append(filterExpressions, "SubmittedAt <= :olderThan")
		expressionValues[":olderThan"] = &dynamodb.AttributeValue{S: aws.String(olderThanTime.Format(time.RFC3339))}
	}

	if len(filterExpressions) > 0 {
		input.FilterExpression = aws.String(strings.Join(filterExpressions, " AND "))
		input.ExpressionAttributeValues = expressionValues
		if len(expressionNames) > 0 {
			input.ExpressionAttributeNames = expressionNames
		}
	}

	var deletedCount int
	result, err := client.Scan(input)
	if err != nil {
		return fmt.Errorf("failed to scan DynamoDB: %w", err)
	}

	for _, item := range result.Items {
		jobID := aws.StringValue(item["JobID"].S)
		deleteInput := &dynamodb.DeleteItemInput{
			TableName: aws.String(resources.JobsTable),
			Key: map[string]*dynamodb.AttributeValue{
				"JobID": {S: aws.String(jobID)},
			},
		}

		_, err := client.DeleteItem(deleteInput)
		if err != nil {
			fmt.Printf("Failed to delete job %s: %v\n", jobID, err)
			continue
		}
		deletedCount++
	}

	fmt.Printf("Successfully deleted %d jobs\n", deletedCount)
	return nil
}

func joinWithAND(expressions []string) string {
	if len(expressions) == 0 {
		return ""
	}
	if len(expressions) == 1 {
		return expressions[0]
	}
	return fmt.Sprintf("%s", strings.Join(expressions, " AND "))
}
