---
title: Configuration
sidebar:
  order: 5
---

The config file lives at `~/.opencodereview/config.json`. You have three ways
to edit it:

- **Interactive TUI** — `ocr config provider` / `ocr config model`, with guided menus.
- **Command line** — `ocr config set <key> <value>`, ideal for scripts and CI.
- **Manual edit (not recommended)** — the JSON file directly (it gets reformatted on the next `ocr config set` write).

## Configuring a model

### Recommended: interactive setup

```bash
ocr config provider
```

It lets you pick a built-in or custom provider, enter an API key, choose a model, saves everything to the config file, and then runs `ocr llm test` once to verify the endpoint. To switch models later:

```bash
ocr config model
```

### Non-interactive setup (CI / no-TUI environments)

Write to the same config with `ocr config set`:

```bash
ocr config set provider                    anthropic
ocr config set model                       claude-opus-4-6
ocr config set providers.anthropic.api_key sk-ant-xxxxxxxxxx
```

### Built-in providers

The following providers ship with OCR, with the Base URL and protocol
preset. Most require an API key. If `providers.<name>.api_key` is unset,
OCR falls back to the corresponding environment variable.

| Name | Protocol | Base URL | API key env var |
|---|---|---|---|
| `anthropic` | anthropic | `https://api.anthropic.com` | `ANTHROPIC_API_KEY` |
| `bedrock` | anthropic-bedrock | derived from `aws_region` | — (AWS credential chain) |
| `openai` | openai | `https://api.openai.com/v1` | `OPENAI_API_KEY` |
| `openai-responses` | openai-responses | `https://api.openai.com/v1` | `OPENAI_RESPONSES_API_KEY` |
| `openrouter` | openai | `https://openrouter.ai/api/v1` | `OPENROUTER_API_KEY` |
| `gemini` | openai | `https://generativelanguage.googleapis.com/v1beta/openai` | `GEMINI_API_KEY` |
| `dashscope` | openai | `https://dashscope.aliyuncs.com/compatible-mode/v1` | `DASHSCOPE_API_KEY` |
| `dashscope-tokenplan` | openai | `https://token-plan.cn-beijing.maas.aliyuncs.com/compatible-mode/v1` | `DASHSCOPE_TOKENPLAN_KEY` |
| `volcengine` | openai | `https://ark.cn-beijing.volces.com/api/v3` | `ARK_API_KEY` |
| `deepseek` | openai | `https://api.deepseek.com` | `DEEPSEEK_API_KEY` |
| `tencent-tokenhub` | openai | `https://tokenhub.tencentmaas.com/v1` | `TENCENT_TOKENHUB_API_KEY` |
| `hy-tokenplan` | openai | `https://api.lkeap.cloud.tencent.com/plan/v3` | `TENCENT_HUNYUAN_TOKENPLAN_KEY` |
| `iflytek` | openai | `https://spark-api-open.xf-yun.com/v1` | `SPARK_API_KEY` |
| `kimi` | openai | `https://api.moonshot.cn/v1` | `MOONSHOT_API_KEY` |
| `kimi-global` | openai | `https://api.moonshot.ai/v1` | `MOONSHOT_GLOBAL_API_KEY` |
| `z-ai` | openai | `https://open.bigmodel.cn/api/paas/v4` | `Z_AI_API_KEY` |
| `mimo` | openai | `https://api.xiaomimimo.com/v1` | `MIMO_API_KEY` |
| `minimax` | openai | `https://api.minimax.io/v1` | `MINIMAX_GLOBAL_API_KEY` |
| `minimax-cn` | openai | `https://api.minimaxi.com/v1` | `MINIMAX_API_KEY` |
| `baidu-qianfan` | openai | `https://qianfan.baidubce.com/v2` | `QIANFAN_API_KEY` |
| `siliconflow`  | openai | `https://api.siliconflow.com/v1` | `SILICONFLOW_GLOBAL_API_KEY` |
| `siliconflow-cn`  | openai | `https://api.siliconflow.cn/v1` | `SILICONFLOW_API_KEY` |
| `novita` | openai | `https://api.novita.ai/openai` | `NOVITA_API_KEY` |
| `xai` | openai | `https://api.x.ai/v1` | `XAI_API_KEY` |

Built-in providers' model lists are suggestions for `ocr config model`, not
restrictions on `--model`. An override absent from both the built-in list and
`providers.<name>.models` produces a warning on stderr; the provider validates
the model when the request is sent. Custom providers retain their existing
`--model` validation rules.

For `minimax` and `minimax-cn`, new model selections suggest
`MiniMax-M3.1-Flash` first, followed by `MiniMax-M3`. Existing saved models and
explicit legacy or custom model IDs remain accepted by OCR. Gateway presets keep their
own model IDs; a direct MiniMax model ID does not imply gateway availability.

The default MiniMax request omits `thinking` and `reasoning_effort`, allowing
[MiniMax's API defaults](https://platform.minimax.io/docs/api-reference/text-openai-api)
to apply. Flash requires adaptive thinking and defaults to `max` effort; explicit
`extra_body` settings can use `thinking: {"type":"adaptive"}` with
`reasoning_effort` set to `low`, `medium`, `high`, `xhigh`, or `max`. M3 also
supports `thinking: {"type":"disabled"}`. Saved `extra_body` settings are
preserved, so remove any `thinking: {"type":"disabled"}` override when switching
to Flash.

## Configuring the review language

`language` determines which language review comments are written in;
it defaults to English when unset:

```bash
ocr config set language 中文
ocr config set language English
```

## See Also

- [QuickStart](../quickstart/) — minimal setup and first review.
- [CLI Reference](../cli-reference/) — every flag the review command accepts.
