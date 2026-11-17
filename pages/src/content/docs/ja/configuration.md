---
title: 設定
sidebar:
  order: 5
---

設定ファイルは `~/.opencodereview/config.json` にあります。編集方法は 3 つあります。

- **インタラクティブ TUI** —— `ocr config provider` / `ocr config model`。ガイド付きメニューが表示されます。
- **コマンドライン** —— `ocr config set <key> <value>`。スクリプトや CI に適しています。
- **手動編集（非推奨）** —— この JSON ファイルを直接編集（次回の `ocr config set` 書き込み時に再フォーマットされます）。

## モデルを設定する

### 推奨：インタラクティブ設定

```bash
ocr config provider
```

組み込みまたはカスタムの provider を選択し、API key を入力し、model を選び、すべてを設定ファイルに保存したうえで、`ocr llm test` を 1 回実行してエンドポイントを検証します。あとで model を切り替えるには：

```bash
ocr config model
```

### 非インタラクティブ設定（CI / TUI なし環境）

`ocr config set` で同じ設定に書き込みます。

```bash
ocr config set provider                    anthropic
ocr config set model                       claude-opus-4-6
ocr config set providers.anthropic.api_key sk-ant-xxxxxxxxxx
```

### 組み込み provider

以下の provider が OCR に同梱されており、Base URL とプロトコルがプリセット
されています。多くは API key が必要です。`providers.<name>.api_key` が未設定の
場合は対応する環境変数を使用します。

| 名称 | プロトコル | Base URL | API key 環境変数 |
|---|---|---|---|
| `anthropic` | anthropic | `https://api.anthropic.com` | `ANTHROPIC_API_KEY` |
| `bedrock` | anthropic-bedrock | `aws_region` から決定 | —（AWS 認証情報チェーン） |
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

組み込み provider のモデル一覧は `ocr config model` の選択候補であり、`--model`
の制限には使われません。組み込み一覧にも `providers.<name>.models` にもないモデルを
指定すると、OCR は stderr に警告を出します。モデルはリクエスト送信時に provider が
検証します。カスタム provider には従来の `--model` 検証ルールが適用されます。

`minimax` と `minimax-cn` の新規モデル選択では、`MiniMax-M3.1-Flash`、
`MiniMax-M3` の順に候補を表示します。保存済みモデルや明示的な旧モデル・カスタムモデルの
ID は引き続き OCR で受け付けます。ゲートウェイには固有のモデル ID があり、
MiniMax の直接接続用 ID が使えるとは限りません。

既定のリクエストは `thinking` と `reasoning_effort` を省略し、
[MiniMax API の既定値](https://platform.minimax.io/docs/api-reference/text-openai-api)を使います。
Flash は adaptive thinking が必須で、effort の既定値は `max` です。
`extra_body` では `thinking: {"type":"adaptive"}` と、`low`、`medium`、`high`、
`xhigh`、`max` のいずれかの `reasoning_effort` を指定できます。M3 は
`thinking: {"type":"disabled"}` にも対応します。保存済みの `extra_body` は保持されるため、
Flash に切り替える前に thinking を無効にする設定を削除してください。

## レビュー言語を設定する

`language` はレビューコメントをどの言語で出力するかを決めます。未設定の場合は
デフォルトで英語になります。

```bash
ocr config set language 中文
ocr config set language English
```

## 関連項目

- [クイックスタート](../quickstart/)——最小限のセットアップと初回のレビュー。
- [CLI リファレンス](../cli-reference/)——review コマンドが受け入れる各引数。
