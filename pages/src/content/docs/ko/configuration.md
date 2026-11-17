---
title: 설정
sidebar:
  order: 5
---

설정 파일은 `~/.opencodereview/config.json`에 있습니다. 편집하는 방법은 세 가지입니다:

- **대화형 TUI** — `ocr config provider` / `ocr config model`. 메뉴가 안내합니다.
- **커맨드라인** — `ocr config set <key> <value>`. 스크립트와 CI에 적합합니다.
- **직접 편집(권장하지 않음)** — JSON 파일을 직접 수정합니다(다음 `ocr config set` 기록 때 다시 포맷됩니다).

## 모델 설정 {#configuring-a-model}

### 권장: 대화형 설정 {#recommended-interactive-setup}

```bash
ocr config provider
```

내장 프로바이더나 커스텀 프로바이더를 고르고 API 키를 입력한 뒤 모델을 선택하면, 명령이 모든 값을 설정 파일에 저장하고 `ocr llm test`를 한 번 실행해 엔드포인트를 검증합니다. 나중에 모델을 바꾸려면:

```bash
ocr config model
```

### 비대화형 설정 (CI / TUI가 없는 환경) {#non-interactive-setup-ci-no-tui-environments}

`ocr config set`으로 같은 설정 파일에 기록합니다:

```bash
ocr config set provider                    anthropic
ocr config set model                       claude-opus-4-6
ocr config set providers.anthropic.api_key sk-ant-xxxxxxxxxx
```

### 내장 프로바이더 {#built-in-providers}

다음 프로바이더는 Base URL과 프로토콜이 미리 설정된 채 OCR에 내장되어 있습니다. 대부분은 API 키가 필요하며, `providers.<name>.api_key`가 비어 있으면 OCR은 해당 환경 변수로 대체합니다.

| 이름 | 프로토콜 | Base URL | API 키 환경 변수 |
|---|---|---|---|
| `anthropic` | anthropic | `https://api.anthropic.com` | `ANTHROPIC_API_KEY` |
| `bedrock` | anthropic-bedrock | `aws_region`에서 결정 | — (AWS 자격 증명 체인) |
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

내장 프로바이더의 모델 목록은 `ocr config model`에서 선택할 때 제안하는 목록이며
`--model`을 제한하지 않습니다. 지정한 모델이 내장 목록과
`providers.<name>.models` 모두에 없으면 OCR은 stderr에 경고를 출력합니다.
모델의 유효성은 요청을 보낼 때 프로바이더가 확인합니다. 사용자 정의 프로바이더에는
기존 `--model` 검증 규칙이 적용됩니다.

`minimax`와 `minimax-cn`의 새 모델 선택에서는 `MiniMax-M3.1-Flash`,
`MiniMax-M3` 순으로 추천합니다. OCR은 저장된 모델과 명시적으로 지정한 이전 모델 또는
사용자 정의 모델 ID를 계속 허용합니다. 게이트웨이는 자체 모델 ID를 사용하므로,
MiniMax 직접 연결용 ID가 게이트웨이에서도 제공된다는 의미는 아닙니다.

기본 요청은 `thinking`과 `reasoning_effort`를 생략하여
[MiniMax API 기본값](https://platform.minimax.io/docs/api-reference/text-openai-api)을 사용합니다.
Flash는 adaptive thinking이 필수이며 effort 기본값은 `max`입니다. `extra_body`에서
`thinking: {"type":"adaptive"}`와 `low`, `medium`, `high`, `xhigh`, `max` 중 하나의
`reasoning_effort`를 명시할 수 있습니다. M3는 `thinking: {"type":"disabled"}`도
지원합니다. 저장된 `extra_body`는 유지되므로 Flash로 전환하기 전에 thinking을
비활성화하는 설정을 제거하세요.

## 리뷰 언어 설정 {#configuring-the-review-language}

`language`는 리뷰 코멘트를 어떤 언어로 쓸지 정합니다. 비어 있으면 기본값은 영어입니다:

```bash
ocr config set language 中文
ocr config set language English
```

## 관련 문서 {#see-also}

- [빠른 시작](../quickstart/) — 최소 설정과 첫 리뷰.
- [CLI 레퍼런스](../cli-reference/) — review 명령이 받는 모든 플래그.
