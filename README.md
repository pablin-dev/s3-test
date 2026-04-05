# S3 YAML Processor

## Use Cases
1. Upload YAML file to S3.
2. Download all YAML files from S3, pick the latest version per ID, compile the expression, and upload the result to DynamoDB.

## Running
### Prerequisites
- Go 1.22+

### Running Tests
```bash
go test ./...
```

### Build
```bash
go build -o app ./cmd/app
```