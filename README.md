# OpenAI PDF SDK

A simple Go SDK and CLI tool to extract text from PDF files and interact with the OpenAI API. This tool allows you to extract text from PDFs and upload files or send extracted text to OpenAI for processing (e.g., summarization).

## Features

- **PDF Text Extraction**: Extract plain text from PDF files.
- **OpenAI File Upload**: Upload files to OpenAI for use with Assistants or Fine-tuning.
- **Chat Completion**: Send extracted text to OpenAI's GPT models (e.g., for summarization).

## Prerequisites

- Go 1.21 or higher
- An OpenAI API Key

## Installation

```bash
git clone <repository-url>
cd openai-pdf-sdk
go mod download
```

## Usage

### Environment Setup

Set your OpenAI API key as an environment variable:

```bash
export OPENAI_API_KEY="your-api-key-here"
```

### Running the CLI

The project includes a CLI in `cmd/main.go`. You can run it directly:

```bash
go run cmd/main.go -file path/to/your/document.pdf
```

### Flags

- `-file`: (Required) Path to the PDF file.
- `-purpose`: (Optional) Purpose of the file upload (default: "assistants"). Common values: "assistants", "fine-tune".

### Example

```bash
go run cmd/main.go -file ./sample.pdf -purpose assistants
```

This command will:
1. Extract text from `sample.pdf`.
2. Print a preview of the extracted text.
3. Upload the file to OpenAI.
4. (Demo) Send a snippet of the text to OpenAI for summarization.

## Project Structure

- `cmd/`: Contains the main application entry point.
- `pkg/openai/`: Client for interacting with the OpenAI API.
- `pkg/pdf/`: Utilities for PDF text extraction.

## Testing

To run the tests, use:

```bash
go test ./...
```

Note: Integration tests in `pkg/openai` require the `OPENAI_API_KEY` environment variable to be set.
