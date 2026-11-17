---
title: 配置
sidebar:
  order: 5
---

配置文件在 `~/.opencodereview/config.json`，你有三种方式编辑它：

- **交互式 TUI** —— `ocr config provider` / `ocr config model`，带引导菜单。
- **命令行** —— `ocr config set <key> <value>`，适合脚本与 CI。
- **手动编辑（不推荐）** —— 该 JSON 文件（下次 `ocr config set` 写入时会重新格式化）。

## 配置模型

### 推荐：交互式设置

```bash
ocr config provider
```

它会让你选择一个内置或自定义 provider、填入 API key、挑选 model，保存到配置文件后自动运行一次 `ocr llm test` 验证端点。之后想换模型：

```bash
ocr config model
```

### 非交互设置（CI / 无 TUI 环境）

用 `ocr config set` 写入同一份配置：

```bash
ocr config set provider                    anthropic
ocr config set model                       claude-opus-4-6
ocr config set providers.anthropic.api_key sk-ant-xxxxxxxxxx
```

### 内置 provider

下列 provider 随 OCR 发布，并已预置 Base URL 与协议。多数 provider 需要 API key；
若 `providers.<name>.api_key` 未设置，会自动回退到对应的环境变量。

| 名称 | 协议 | Base URL | API key 环境变量 |
|---|---|---|---|
| `anthropic` | anthropic | `https://api.anthropic.com` | `ANTHROPIC_API_KEY` |
| `bedrock` | anthropic-bedrock | 由 `aws_region` 决定 | —（AWS 凭证链） |
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
| `siliconflow` | openai | `https://api.siliconflow.com/v1` | `SILICONFLOW_GLOBAL_API_KEY` |
| `siliconflow-cn`  | openai | `https://api.siliconflow.cn/v1` | `SILICONFLOW_API_KEY` |
| `novita` | openai | `https://api.novita.ai/openai` | `NOVITA_API_KEY` |
| `xai` | openai | `https://api.x.ai/v1` | `XAI_API_KEY` |

内置 provider 的模型列表只为 `ocr config model` 提供选择建议，不限制 `--model`。
如果指定的模型既不在内置列表中，也不在 `providers.<name>.models` 中，OCR 会向
stderr 输出警告；发送请求时由 provider 验证模型。自定义 provider 仍遵循原有的
`--model` 校验规则。

`minimax` 和 `minimax-cn` 的新模型选择优先推荐 `MiniMax-M3.1-Flash`，其次为
`MiniMax-M3`。OCR 仍接受已保存的模型及显式指定的旧版或自定义模型 ID。
网关使用各自的模型 ID；MiniMax 直连模型 ID 并不代表网关已支持该模型。

默认请求不设置 `thinking` 和 `reasoning_effort`，由
[MiniMax API 默认值](https://platform.minimax.io/docs/api-reference/text-openai-api)决定。
Flash 必须使用自适应思考，默认 effort 为 `max`。在 `extra_body` 中可显式设置
`thinking: {"type":"adaptive"}`，并将 `reasoning_effort` 设为 `low`、`medium`、
`high`、`xhigh` 或 `max`。M3 还支持 `thinking: {"type":"disabled"}`。
已有 `extra_body` 配置会保留，因此切换到 Flash 前应移除禁用思考的设置。

## 配置评审语言

`language` 决定评审评论用哪种语言输出，未设置时默认英文：

```bash
ocr config set language 中文
ocr config set language English
```

## 另见

- [快速开始](../quickstart/)——最小化设置与首次评审。
- [CLI 参考](../cli-reference/)——review 命令接受的每个参数。
