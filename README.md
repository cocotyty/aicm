# AI Commit Messages Generator

AI Commit Messages Generator is a powerful tool that automatically generates meaningful commit messages for your code changes.

NOTICE: This project is AI-generated and may have unexpected behaviors. Please use with caution.

## Configuration
Configure your settings using the following commands:
```
aicm config set LLM_API_KEY=your_api_key
aicm config set LLM_MODEL=your_model
aicm config set LLM_API_URL=your_api_url
```

The configuration will be stored in `~/.aicm/config.json`.

## Usage

Simply run:
```
aicm
```

## How it works

The tool follows these steps:

1. Detects code changes in the current directory using:
```
git status
```
2. Reads modified files using:
```
git diff --cached [file_name]
```
3. Generates commit messages through LLM in two phases:

### Phase 1: Code Changes Description

```go
type Request struct {
    CodeFileChanges []CodeFileChanges `json:"code_file_changes"`
}

type CodeFileChanges struct {
    FileName string `json:"file_name"`
    Content string `json:"content"`
    Diff string `json:"diff"`
}
```

```go
type Response struct {
   map[string]string  `json:"descriptions"` // key: file name, value: description
}
```

```prompt
You are a code changes analyzer. Please generate descriptions for the following changes:

{{.Request}}

Response format:
```json
{{.Response}}
```

### Phase 2: Commit Message Generation

```json
{
    "removed": [
        "file_name"
    ],
    "added": {
        "file_name" : {
            "description": "",
        }
    },
    "changed": {
        "file_name" : {
            "description": "",
        }
    }
}
```

```prompt
You are a commit message generator. Follow these guidelines:
- Format: "type: description"
- Types:
  * feat: New feature
  * fix: Bug fix
  * docs: Documentation changes
  * style: Code formatting, no logic changes
  * refactor: Code restructuring
  * perf: Performance improvements
  * test: Test additions
  * chore: Maintenance tasks

Generate a commit message for:
{{code_changes}}
```

