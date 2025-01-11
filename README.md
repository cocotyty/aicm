# AI Commit Messages Generator

AI Commit Messages Generator is a tool that helps you generate commit messages for your code changes.

## Configuration
```
aicm config set LLM_API_KEY=your_api_key
aicm config set LLM_MODEL=your_model
aicm config set LLM_API_URL=your_api_url
```

The configuration will be saved in the `~/.aicm/config.json` file.

## Usage

```
aicm
```

## How it works

First, it will read the code changes from the current directory.
```
git status
```
then all the files that are modified will be read.

```
git diff --cached [file_name]
```

Then it will generate the commit message by the LLM.

We will have two steps:

1. Generate the code file changes description.
2. Generate the commit message.

### 1. Generate the code file changes description.

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
   map[string]string  `json:"descriptions"` // key is the file name, value is the description
}
```

```prompt
You are a code file changes description generator.

The response should be in the format of:
```json
{{.Response}}

Please generate the description for the following code file changes:

{{.Request}}

```


### 2. Generate the commit message.


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
You are a commit message generator.
The commit message should be in the format of "type: description".
type should be one of the following:
- feat: a new feature
- fix: a bug fix
- docs: documentation only changes
- style: formatting, missing semi colons, etc; no production code change
- refactor: refactoring production code
- perf: performance improvements
- test: adding missing tests
- chore: updating build tasks, package manager, etc; no production code change
Please generate a commit message for the following code changes:

{{code_changes}}
```
