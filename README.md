# AutoCommit

CLI tool that generates commit messages using LLMs.

```bash
git add .
autocommit
# feat(api): add user pagination endpoint
```

## Install

```bash
go install github.com/josinSbazin/AutoCommit/cmd/autocommit@latest
```

Or download binary from [Releases](https://github.com/josinSbazin/AutoCommit/releases).

## Setup

```bash
export OPENAI_API_KEY=sk-...
# or
export ANTHROPIC_API_KEY=sk-ant-...
```

Run `autocommit init` for interactive setup.

## Usage

```bash
git add .
autocommit

# dry run
autocommit generate --dry-run

# install git hook for automatic generation
autocommit hook install
```

With hook installed, just run `git commit` — message will be generated automatically.

## Providers

| Provider | Env variable |
|----------|--------------|
| OpenAI | `OPENAI_API_KEY` |
| Anthropic | `ANTHROPIC_API_KEY` |
| GigaChat | `GIGACHAT_CLIENT_ID` + `GIGACHAT_CLIENT_SECRET` |
| YandexGPT | `YANDEX_API_KEY` + `YANDEX_FOLDER_ID` |
| Ollama | Local, no key needed |
| OpenAI-compatible | `AUTOCOMMIT_API_KEY` + endpoint in config |

Auto-detection: if env key exists, provider is selected automatically.

## Config

`.autocommit.yml` in project root:

```yaml
provider: openai
model: gpt-4o
style: conventional
language: en
max_subject_length: 72
include_body: true
```

## Commands

```
autocommit              Generate and commit
autocommit generate     Show message only
autocommit generate -d  Dry run
autocommit init         Setup wizard
autocommit config       Show current config
autocommit hook install Install git hook
autocommit doctor       Diagnostics
```

## Local mode

Use [Ollama](https://ollama.ai) for offline usage:

```bash
ollama pull llama3.1
autocommit init  # select Ollama
```

## Build

```bash
git clone https://github.com/josinSbazin/AutoCommit.git
cd AutoCommit
make build
```

## License

MIT
