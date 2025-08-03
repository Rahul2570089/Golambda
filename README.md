# GoLambda

A lightweight serverless function platform built in Go that allows you to register, compile, and execute Go functions dynamically via HTTP endpoints or cron schedules.

## Features

- **Dynamic Function Registration**: Register Go functions at runtime via REST API
- **HTTP Triggers**: Execute functions via HTTP endpoints
- **Cron Scheduling**: Schedule functions to run at specific intervals using cron expressions
- **Real-time Compilation**: Automatically compiles user-submitted Go code into executable binaries
- **Comprehensive Logging**: Detailed logging with structured fields using logrus
- **Function Metadata**: Maintains registry of registered functions with metadata

## Architecture

```
golambda/
├── main.go                # Main server and HTTP handlers
├── manager/
│   ├── execute.go         # Binary execution logic
│   └── register.go        # Function registration and compilation
├── models/
│   └── models.go          # Data models for functions and metadata
├── orchestrator/
│   ├── api.go             # HTTP route registration
│   └── cron.go            # Cron job scheduling
├── plugins/               # Compiled function binaries (.exe files)
├── user_functions/        # User-submitted Go source files
├── registry/
│   └── registry.json      # Function metadata registry
├── go.mod
└── go.sum
```

## Dependencies

- **github.com/robfig/cron/v3**: Cron job scheduling
- **github.com/sirupsen/logrus**: Structured logging
- **golang.org/x/sys**: System-level operations

## Getting Started

### Prerequisites

- Go 1.24.4 or later
- Windows environment (binaries are compiled as `.exe` files)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd golambda
```

2. Install dependencies:
```bash
go mod download
```

3. Create required directories:
```bash
mkdir -p plugins user_functions registry
```

4. Run the server:
```bash
go run main.go
```

The server will start on port 8080.

## API Endpoints

### Health Check
```
GET /
```
Returns server status and confirms the service is running.

### Register Function
```
PUT /register
```

Register a new Go function with HTTP or cron trigger.

**Request Body:**
```json
{
    "name": "function_name",
    "trigger": "http", // or cron expression like "0 */5 * * * *"
    "code": "base64_encoded_go_code"
}
```

**Example:**
```json
{
    "name": "hello",
    "trigger": "http",
    "code": "cGFja2FnZSBtYWluCgppbXBvcnQgImZtdCIKCmZ1bmMgbWFpbigpIHsKICAgIGZtdC5QcmludGxuKCJIZWxsbywgV29ybGQhIikKfQ=="
}
```

### Execute Function (HTTP Trigger)
```
GET /{function_name}
```
Executes a function registered with HTTP trigger.

## Usage Examples

### 1. Register an HTTP Function

```bash
# Create a simple Hello World function
echo 'package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}' > hello.go

# Base64 encode the function
CODE=$(base64 -w 0 hello.go)

# Register the function
curl -X PUT http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"hello\",
    \"trigger\": \"http\",
    \"code\": \"$CODE\"
  }"
```

### 2. Execute the Function

```bash
curl http://localhost:8080/hello
```

Response:
```json
{
    "message": "Hello, World!\n"
}
```

### 3. Register a Cron Function

```bash
# Register a function that runs every 5 minutes
curl -X PUT http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"scheduler\",
    \"trigger\": \"0 */5 * * * *\",
    \"code\": \"$CODE\"
  }"
```

## Directory Structure Details

- **`plugins/`**: Contains compiled executable binaries of registered functions
- **`user_functions/`**: Stores the original Go source code files submitted by users
- **`registry/`**: Contains `registry.json` with function metadata including name, trigger type, and binary path

## Logging

The application uses structured logging with the following information:
- HTTP request details (method, path)
- Function registration events
- Function execution results
- Error conditions with detailed context

## Error Handling

- **Invalid HTTP methods**: Returns 405 Method Not Allowed
- **Malformed JSON**: Returns 400 Bad Request
- **Compilation errors**: Returns 501 Not Implemented
- **Execution failures**: Returns 501 Not Implemented with error details

## Cron Expression Format

Supports standard cron expressions with seconds precision:

```
┌─────────────────────  seconds (0-59)
│ ┌───────────────────  minute (0-59)
│ │ ┌─────────────────  hour (0-23)
│ │ │ ┌───────────────  day of month (1-31)
│ │ │ │ ┌─────────────  month (1-12)
│ │ │ │ │ ┌───────────  day of week (0-6)
│ │ │ │ │ │
* * * * * *
```

Examples:
- `0 0 * * * *` - Every hour
- `0 */5 * * * *` - Every 5 minutes
- `0 0 9 * * MON-FRI` - 9 AM on weekdays

## Future Improvements

- **Cross-platform Support**: Linux and macOS compatibility
- **Authentication**: JWT-based API authentication and rate limiting
- **Function Versioning**: Manage multiple versions of functions
- **Resource Limits**: CPU, memory, and execution time constraints
- **Database Registry**: Replace JSON file with proper database
- **Multi-language Support**: Python, Node.js, Rust runtimes
- **Containerization**: Docker isolation for security and resource management
- **CLI Tool**: Command-line interface for function management
- **Hot Reloading**: Update functions without server restart
