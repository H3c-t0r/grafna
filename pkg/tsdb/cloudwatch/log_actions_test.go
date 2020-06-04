package cloudwatch

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana/pkg/components/simplejson"
	"github.com/grafana/grafana/pkg/tsdb"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (m mockedLogs) GetQueryResultsWithContext(ctx context.Context, input *cloudwatchlogs.GetQueryResultsInput, option ...request.Option) (*cloudwatchlogs.GetQueryResultsOutput, error) {
	return &m.queryResults, nil
}

//***
// LogActions Tests
//***

func TestHandleDescribeLogGroups_WhenLogGroupNamePrefixIsEmpty(t *testing.T) {
	logs := mockedLogs{
		logGroups: cloudwatchlogs.DescribeLogGroupsOutput{
			LogGroups: []*cloudwatchlogs.LogGroup{
				{
					LogGroupName: aws.String("group_a"),
				},
				{
					LogGroupName: aws.String("group_b"),
				},
				{
					LogGroupName: aws.String("group_c"),
				},
			},
		},
	}
	executor := &CloudWatchExecutor{
		DataSource: mockDatasource(),
		clients: &mockClients{
			logs: logs,
		},
	}

	params := simplejson.NewFromAny(map[string]interface{}{
		"limit": 50,
	})

	frame, err := executor.handleDescribeLogGroups(context.Background(), logs, params)

	expectedField := data.NewField("logGroupName", nil, []*string{aws.String("group_a"), aws.String("group_b"), aws.String("group_c")})
	expectedFrame := data.NewFrame("logGroups", expectedField)

	assert.Equal(t, nil, err)
	assert.Equal(t, expectedFrame, frame)
}

func TestHandleDescribeLogGroups_WhenLogGroupNamePrefixIsNotEmpty(t *testing.T) {
	logs := mockedLogs{
		logGroups: cloudwatchlogs.DescribeLogGroupsOutput{
			LogGroups: []*cloudwatchlogs.LogGroup{
				{
					LogGroupName: aws.String("group_a"),
				},
				{
					LogGroupName: aws.String("group_b"),
				},
				{
					LogGroupName: aws.String("group_c"),
				},
			},
		},
	}
	executor := &CloudWatchExecutor{
		DataSource: mockDatasource(),
		clients: &mockClients{
			logs: logs,
		},
	}

	params := simplejson.NewFromAny(map[string]interface{}{
		"logGroupNamePrefix": "g",
	})

	frame, err := executor.handleDescribeLogGroups(context.Background(), logs, params)

	expectedField := data.NewField("logGroupName", nil, []*string{aws.String("group_a"), aws.String("group_b"), aws.String("group_c")})
	expectedFrame := data.NewFrame("logGroups", expectedField)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedFrame, frame)
}

func TestHandleGetLogGroupFields_WhenLogGroupNamePrefixIsNotEmpty(t *testing.T) {
	logs := mockedLogs{
		logGroupFields: cloudwatchlogs.GetLogGroupFieldsOutput{
			LogGroupFields: []*cloudwatchlogs.LogGroupField{
				{
					Name:    aws.String("field_a"),
					Percent: aws.Int64(100),
				},
				{
					Name:    aws.String("field_b"),
					Percent: aws.Int64(30),
				},
				{
					Name:    aws.String("field_c"),
					Percent: aws.Int64(55),
				},
			},
		},
	}
	executor := &CloudWatchExecutor{
		DataSource: mockDatasource(),
		clients: &mockClients{
			logs: logs,
		},
	}

	params := simplejson.NewFromAny(map[string]interface{}{
		"logGroupName": "group_a",
		"limit":        50,
	})

	frame, err := executor.handleGetLogGroupFields(context.Background(), logs, params, "A")

	expectedNameField := data.NewField("name", nil, []*string{aws.String("field_a"), aws.String("field_b"), aws.String("field_c")})
	expectedPercentField := data.NewField("percent", nil, []*int64{aws.Int64(100), aws.Int64(30), aws.Int64(55)})
	expectedFrame := data.NewFrame("A", expectedNameField, expectedPercentField)
	expectedFrame.RefID = "A"

	assert.Equal(t, nil, err)
	assert.Equal(t, expectedFrame, frame)
}

func TestExecuteStartQuery(t *testing.T) {
	logs := mockedLogs{}
	executor := &CloudWatchExecutor{
		DataSource: mockDatasource(),
		clients: &mockClients{
			logs: logs,
		},
	}

	timeRange := &tsdb.TimeRange{
		From: "1584873443000",
		To:   "1584700643000",
	}

	params := simplejson.NewFromAny(map[string]interface{}{
		"region":      "default",
		"limit":       50,
		"queryString": "fields @message",
	})

	response, err := executor.executeStartQuery(context.Background(), logs, params, timeRange)

	var expectedResponse *cloudwatchlogs.StartQueryOutput = nil

	assert.Equal(t, expectedResponse, response)
	assert.Equal(t, fmt.Errorf("invalid time range: start time must be before end time"), err)

}

func TestHandleStartQuery(t *testing.T) {
	logs := mockedLogs{}
	executor := &CloudWatchExecutor{
		DataSource: mockDatasource(),
		clients: &mockClients{
			logs: logs,
		},
	}

	timeRange := &tsdb.TimeRange{
		From: "1584700643000",
		To:   "1584873443000",
	}

	params := simplejson.NewFromAny(map[string]interface{}{
		"region":      "default",
		"limit":       50,
		"queryString": "fields @message",
	})

	frame, err := executor.handleStartQuery(context.Background(), logs, params, timeRange, "A")

	expectedField := data.NewField("queryId", nil, []string{"abcd-efgh-ijkl-mnop"})
	expectedFrame := data.NewFrame("A", expectedField)
	expectedFrame.RefID = "A"
	expectedFrame.Meta = &data.FrameMeta{
		Custom: map[string]interface{}{
			"Region": "default",
		},
	}

	assert.Equal(t, nil, err)
	assert.Equal(t, expectedFrame, frame)
}

func TestHandleStopQuery(t *testing.T) {
	logs := mockedLogs{}
	executor := &CloudWatchExecutor{
		DataSource: mockDatasource(),
		clients: &mockClients{
			logs: logs,
		},
	}

	params := simplejson.NewFromAny(map[string]interface{}{
		"queryId": "abcd-efgh-ijkl-mnop",
	})

	frame, err := executor.handleStopQuery(context.Background(), logs, params)

	expectedField := data.NewField("success", nil, []bool{true})
	expectedFrame := data.NewFrame("StopQueryResponse", expectedField)

	assert.Equal(t, nil, err)
	assert.Equal(t, expectedFrame, frame)
}

func TestHandleGetQueryResults(t *testing.T) {
	logs := mockedLogs{
		queryResults: cloudwatchlogs.GetQueryResultsOutput{
			Results: [][]*cloudwatchlogs.ResultField{
				{
					{
						Field: aws.String("@timestamp"),
						Value: aws.String("2020-03-20 10:37:23.000"),
					},
					{
						Field: aws.String("field_b"),
						Value: aws.String("b_1"),
					},
					{
						Field: aws.String("@ptr"),
						Value: aws.String("abcdefg"),
					},
				},

				{
					{
						Field: aws.String("@timestamp"),
						Value: aws.String("2020-03-20 10:40:43.000"),
					},
					{
						Field: aws.String("field_b"),
						Value: aws.String("b_2"),
					},
					{
						Field: aws.String("@ptr"),
						Value: aws.String("hijklmnop"),
					},
				},
			},

			Statistics: &cloudwatchlogs.QueryStatistics{
				BytesScanned:   aws.Float64(512),
				RecordsMatched: aws.Float64(256),
				RecordsScanned: aws.Float64(1024),
			},

			Status: aws.String("Complete"),
		},
	}
	executor := &CloudWatchExecutor{
		DataSource: mockDatasource(),
		clients: &mockClients{
			logs: logs,
		},
	}

	params := simplejson.NewFromAny(map[string]interface{}{
		"queryId": "abcd-efgh-ijkl-mnop",
	})

	frame, err := executor.handleGetQueryResults(context.Background(), logs, params, "A")
	require.NoError(t, err)
	timeA, err := time.Parse("2006-01-02 15:04:05.000", "2020-03-20 10:37:23.000")
	require.NoError(t, err)
	timeB, err := time.Parse("2006-01-02 15:04:05.000", "2020-03-20 10:40:43.000")
	require.NoError(t, err)
	expectedTimeField := data.NewField("@timestamp", nil, []*time.Time{
		aws.Time(timeA), aws.Time(timeB),
	})
	expectedTimeField.SetConfig(&data.FieldConfig{DisplayName: "Time"})

	expectedFieldB := data.NewField("field_b", nil, []*string{
		aws.String("b_1"), aws.String("b_2"),
	})

	expectedFrame := data.NewFrame("A", expectedTimeField, expectedFieldB)
	expectedFrame.RefID = "A"

	expectedFrame.Meta = &data.FrameMeta{
		Custom: map[string]interface{}{
			"Status": "Complete",
			"Statistics": cloudwatchlogs.QueryStatistics{
				BytesScanned:   aws.Float64(512),
				RecordsMatched: aws.Float64(256),
				RecordsScanned: aws.Float64(1024),
			},
		},
	}

	assert.Equal(t, nil, err)
	assert.ElementsMatch(t, expectedFrame.Fields, frame.Fields)
	assert.Equal(t, expectedFrame.Meta, frame.Meta)
}
