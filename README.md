# AIPipe

A command-line tool for transforming text or otherwise making adhoc LLM calls.

## Usage

Use the pipe operator (`|`) to send text through the `aipipe` command:

```bash
some-app | aipipe "instruction"
```

...or just call it by itself for shell based one-off llm calls.

### Example

Simple reformatting

Input:
```bash
echo "Robin Anderson 1 High Street CB1 1AA" | aipipe "format as JSON" --cb
```

Output:
```json
{
    "name": "Robin Anderson",
    "address": {
        "street": "1 High Street",
        "postcode": "CB1 1AA"
    }
}
```

The 'code block' flag `--cb` is best for when you want something specifically formatted, rather than just for you to read yourself. Without it the LLM might write "Sure, here's your thing..." which you might not want to pipe into another application.

### Options

- `-c / --cb`: outputs only the first code block emitted by the LLM, discarding all other output. Otherwise all output is emitted to std out.
- `-p / --pretty`: use console colours to highlight markdown.
- `-s / --stream`: stream the output for faster perceived response
- `-r / --reasoning`: use a reasoning model instead, for extra oomph.
- `-f / --fast`: use a fast-but-thick model instead, for extra speed.

## Installation

build.ps1 (windows)
build.sh  (!windows)

copy the binary produced to your bin folder.

Set env vars
```
GROQ_API_KEY
# OR
OPENAI_API_KEY
```

Groq is used in preference OpenAI if both api keys are defined, since this application is meant for speed.

However, you can override with any openai compatible provider:

```
AIPIPE_API_KEY=xxx
AIPIPE_ENDPOINT=https://some-provider.example.com/v1
```

as well as storing stuff in `~/.aipipe/config.yaml`

```yaml
apiKey: xxx
endpoint: https://...
defaultMode: gpt-5o
reasoningModel: 6o-mini
fastModel: llama-7.1-1b-nano
```

## Syntax highlighting

`-p` mode will make markdown formatted output more colourful, as well as applying syntax highlighting to the contents of codeblocks.