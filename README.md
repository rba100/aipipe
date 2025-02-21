# AIPipe

A command-line tool for transforming text using AI-powered formatting.

## Usage

Use the pipe operator (`|`) to send text through the `aipipe` command:

```bash
some-app | aipipe "instruction"
```

### Example

Simple reformatting, using the 'code block' flag to capture specific output from the LLM.

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

### Options

- `--cb`: outputs only the first code block emitted by the LLM, discarding all other output. Otherwise all output is emitted to std out.
- `--r1`: use a reasoning model instead, for extra oomph,
- `--fast`: use a fast-but-thick model instead, for extra speed.

## Installation

`dotnet publish aipipe.csproj -c Release -o bld --self-contained true -p:PublishSingleFile=true -p:PublishTrimmed=true -p:DebugType=None`

copy the .\bld\aipipe(.exe) to your bin folder.

Set env vars
```
GROQ_API_KEY
```

## AI Providers

Supports https://groq.com/ and https://openrouter.ai/

To use OpenRouter, use the `--or` flag. If you want this the default you'll need a config file (work in progress).