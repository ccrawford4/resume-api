# Resume Analysis API

A Go-based HTTP server that analyzes resumes against job descriptions using OpenAI's GPT model. The API fetches PDF resumes from Google Cloud Storage, extracts their content, and provides detailed analysis including match percentage, strengths, improvements, and missing skills.

## Features

- **PDF Resume Parsing**: Extracts text content from PDF files stored in Google Cloud Storage
- **AI-Powered Analysis**: Uses OpenAI GPT to analyze resume-job description matches
- **Batch Processing**: Analyze multiple resumes at once or individual resumes
- **Detailed Insights**: Provides strengths, improvements, and missing skills analysis
- **RESTful API**: Clean HTTP endpoints with JSON request/response format

## Prerequisites

- Go 1.21 or higher
- Google Cloud Platform account with a Storage bucket
- OpenAI API key
- Service account with GCS access

## Setup

### 1. Clone and Install Dependencies

```bash
git clone git@github.com:ccrawford4/resume-api.git
cd resume-api
go mod tidy
```

### 2. Google Cloud Storage Setup

1. **Create a Service Account:**
   - Go to [Google Cloud Console](https://console.cloud.google.com/)
   - Navigate to IAM & Admin > Service Accounts
   - Create a new service account
   - Assign the `Storage Object Viewer` role (or higher if needed)
   - Download the JSON key file

2. **Set up Authentication:**
   ```bash
   # Copy your service account JSON to keys.json in the root directory
   cp /path/to/your/service-account.json keys.json
   
   # Set environment variables
   export GCS_BUCKET_NAME="your-resume-bucket-name"
   export OPENAI_API_KEY="your-openai-api-key"
   ```

### 3. Prepare Your Resume Bucket

- Upload PDF resume files to your Google Cloud Storage bucket
- Ensure the service account has read access to the bucket

### 4. Run the Server

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### 1. Health Check

**GET** `/`

Returns a simple health check message.

**Response:**
```
hello world
```

### 2. Analyze All Resumes

**POST** `/`

Analyzes all PDF resumes in the GCS bucket against a job description.

**Request Body:**
```json
{
  "jobDescription": "Your detailed job description here..."
}
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "1",
      "name": "Software_Engineer_Resume_2024.pdf",
      "uploadDate": "2024-01-15",
      "fileUrl": "/api/resume-analysis/resumes/1/preview",
      "matchPercentage": 87,
      "insights": {
        "strengths": [
          "Strong React and TypeScript experience",
          "Previous experience with modern frontend frameworks",
          "Good understanding of state management"
        ],
        "improvements": [
          "Could highlight more testing experience",
          "Add specific examples of performance optimization"
        ],
        "missingSkills": [
          "Next.js experience not mentioned",
          "Limited backend integration examples"
        ]
      }
    }
  ],
  "totalResumes": 3
}
```

### 3. Analyze Single Resume

**POST** `/new-resume`

Analyzes a specific resume file against a job description.

**Request Body:**
```json
{
  "resumeName": "Software_Engineer_Resume_2024.pdf",
  "jobDescription": "Your detailed job description here..."
}
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "1",
      "name": "Software_Engineer_Resume_2024.pdf",
      "uploadDate": "2024-01-15",
      "fileUrl": "/api/resume-analysis/resumes/1/preview",
      "matchPercentage": 87,
      "insights": {
        "strengths": [...],
        "improvements": [...],
        "missingSkills": [...]
      }
    }
  ],
  "totalResumes": 1
}
```

## Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `GCS_BUCKET_NAME` | Google Cloud Storage bucket name | Yes | "user-resumes-hs-hackathon" |
| `OPENAI_API_KEY` | OpenAI API key for GPT analysis | Yes | - |

## Project Structure

```
resume-api/
├── main.go              # HTTP server and routes
├── go.mod               # Go module dependencies
├── go.sum               # Dependency checksums
├── keys.json            # Google Cloud service account credentials
├── gcloud/
│   └── client.go        # Google Cloud Storage client
├── openai/
│   └── client.go        # OpenAI API client
└── parser/
    └── parser.go        # PDF parsing utilities
```

## Error Handling

The API returns appropriate HTTP status codes:

- `200`: Success
- `400`: Bad request (invalid JSON, missing required fields)
- `500`: Internal server error (GCS/OpenAI connection issues)

Error responses include:
```json
{
  "success": false,
  "error": "Error description"
}
```

## Development

### Adding New Features

1. **New Endpoints**: Add routes in `main.go`
2. **GCS Operations**: Extend `gcloud/client.go`
3. **AI Analysis**: Modify `openai/client.go`
4. **File Parsing**: Update `parser/parser.go`

### Testing

```bash
# Test health endpoint
curl http://localhost:8080/

# Test resume analysis
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{"jobDescription": "Software Engineer position..."}'
```

## Troubleshooting

### Common Issues

1. **"could not find default credentials"**
   - Ensure `keys.json` exists in the project root
   - Check that the service account has proper permissions

2. **"bucket doesn't exist"**
   - Check the bucket name in `GCS_BUCKET_NAME`
   - Ensure the service account has access to the bucket

3. **"OpenAI API error"**
   - Verify `OPENAI_API_KEY` is set and valid
   - Check OpenAI API quota and billing

4. **"PDF parsing error"**
   - Ensure files in GCS are valid PDFs
   - Check file permissions and accessibility

## Security Notes

- **Never commit `keys.json`** to version control
- Use environment variables for sensitive configuration
- Consider using Google Cloud's workload identity for production deployments
