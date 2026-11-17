---
title: Конфигурация
sidebar:
  order: 5
---

Файл конфигурации находится по пути `~/.opencodereview/config.json`. Изменить
его можно тремя способами:

- **Интерактивный TUI** — `ocr config provider` / `ocr config model` с пошаговыми меню.
- **Командная строка** — `ocr config set <key> <value>`, удобно для скриптов и CI.
- **Редактирование вручную (не рекомендуется)** — напрямую изменить JSON-файл (при следующей записи через `ocr config set` его форматирование изменится).

## Настройка модели

### Рекомендуемый способ: интерактивная настройка

```bash
ocr config provider
```

Команда позволяет выбрать встроенного или пользовательского провайдера,
ввести API-ключ и выбрать модель, сохраняет всё в файл конфигурации, а затем
один раз запускает `ocr llm test` для проверки эндпоинта. Чтобы позже сменить
модель:

```bash
ocr config model
```

### Неинтерактивная настройка (CI / среды без TUI)

Запишите параметры в тот же файл конфигурации с помощью `ocr config set`:

```bash
ocr config set provider                    anthropic
ocr config set model                       claude-opus-4-6
ocr config set providers.anthropic.api_key sk-ant-xxxxxxxxxx
```

### Встроенные провайдеры

Перечисленные ниже провайдеры поставляются вместе с OCR; для них заранее
заданы базовый URL и протокол. Большинству нужен API-ключ. Если
`providers.<name>.api_key` не задан, OCR использует соответствующую переменную
окружения.

| Имя | Протокол | Базовый URL | Переменная окружения для API-ключа |
|---|---|---|---|
| `anthropic` | anthropic | `https://api.anthropic.com` | `ANTHROPIC_API_KEY` |
| `bedrock` | anthropic-bedrock | определяется `aws_region` | — (цепочка учётных данных AWS) |
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

Списки моделей встроенных провайдеров служат подсказкой для `ocr config model` и
не ограничивают `--model`. Если модель отсутствует и во встроенном списке, и в
`providers.<name>.models`, OCR выводит предупреждение в stderr; провайдер
проверит модель при отправке запроса. Для пользовательских провайдеров действуют
прежние правила проверки `--model`.

При выборе новой модели для `minimax` и `minimax-cn` первой предлагается
`MiniMax-M3.1-Flash`, затем `MiniMax-M3`. OCR по-прежнему принимает сохранённые модели
и явно заданные устаревшие или пользовательские идентификаторы. У шлюзов свои
идентификаторы: наличие модели в прямом API MiniMax не означает её доступность через шлюз.

По умолчанию запрос не содержит `thinking` и `reasoning_effort`, поэтому применяются
[значения по умолчанию API MiniMax](https://platform.minimax.io/docs/api-reference/text-openai-api).
Flash требует адаптивного рассуждения; значение effort по умолчанию — `max`.
В `extra_body` можно явно задать `thinking: {"type":"adaptive"}` и `reasoning_effort`
со значением `low`, `medium`, `high`, `xhigh` или `max`. M3 также поддерживает
`thinking: {"type":"disabled"}`. Сохранённый `extra_body` не меняется автоматически:
перед переходом на Flash удалите настройку, отключающую рассуждение.

## Настройка языка ревью

Параметр `language` определяет, на каком языке будут написаны комментарии
ревью; если он не задан, используется английский:

```bash
ocr config set language 中文
ocr config set language English
```

## См. также

- [Быстрый старт](../quickstart/) — минимальная настройка и первое ревью.
- [Справочник CLI](../cli-reference/) — все флаги, принимаемые командой ревью.
