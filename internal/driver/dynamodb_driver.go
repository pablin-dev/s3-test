package driver

// Removed unused imports: "github.com/aws/aws-sdk-go-v2/aws"

// DynamoDBDriver implements storage operations using AWS DynamoDB.
type DynamoDBDriver struct {
	client    DynamoDBClientAPI // Uses the interface
	tableName string
}

// NewDynamoDBDriver creates a new instance of DynamoDBDriver with an injected client.
func NewDynamoDBDriver(client DynamoDBClientAPI, tableName string) *DynamoDBDriver {
	return &DynamoDBDriver{
		client:    client,
		tableName: tableName,
	}
}
