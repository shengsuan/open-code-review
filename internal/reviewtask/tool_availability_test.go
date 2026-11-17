// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 alibaba/open-code-review Contributors

package reviewtask

import (
	"context"
	"testing"

	"github.com/alibaba/open-code-review/internal/llm"
	"github.com/alibaba/open-code-review/internal/session"
	"github.com/alibaba/open-code-review/internal/tool"
)

type countingToolProvider struct {
	name  string
	calls int
}

func (p *countingToolProvider) Tool() tool.Tool {
	if tool.IsReserved(p.name) {
		return tool.OfName(p.name)
	}
	return tool.Dynamic(p.name)
}
func (p *countingToolProvider) Execute(context.Context, map[string]any) (string, error) {
	p.calls++
	return "executed", nil
}

func toolCallResponse(name, args string) *llm.ChatResponse {
	response := taskDoneResponseWithArguments(args)
	response.Choices[0].Message.ToolCalls[0].Function.Name = name
	return response
}

func TestRunMainTaskRejectsUnadvertisedTools(t *testing.T) {
	for _, tc := range []struct {
		name, arguments string
		noTools         bool
	}{
		{"file_read", `{}`, false},
		{"file_read_diff", `{}`, false},
		{"external_lookup", `{}`, false},
		{"task_done", `{}`, false},
		{"task_done", `{"state":"FAILED"}`, false},
		{"code_comment", `{"comments":[{"content":"issue","existing_code":"x"}]}`, false},
		{"file_read", `{}`, true},
	} {
		t.Run(tc.name+tc.arguments, func(t *testing.T) {
			client := &fakeClient{responses: []*llm.ChatResponse{
				toolCallResponse(tc.name, tc.arguments), toolCallResponse("allowed_tool", `{}`),
			}}
			deps := newTestDeps(client)
			deps.Template.MaxToolRequestTimes = 2
			deps.Session = session.New(t.TempDir(), "main", "fake", session.SessionOptions{})
			deps.MainToolDefs = []llm.ToolDef{{Type: "function", Function: llm.FunctionDef{Name: "allowed_tool"}}}
			if tc.noTools {
				deps.MainToolDefs = nil
			}
			rejected := &countingToolProvider{name: tc.name}
			allowed := &countingToolProvider{name: "allowed_tool"}
			deps.Tools.Register(rejected)
			deps.Tools.Register(allowed)
			runner := NewRunner(deps)
			completed, _, err := runner.RunMainTask(context.Background(), []llm.Message{llm.NewTextMessage("user", "review")}, "main.go")
			if err != nil || completed || len(client.requests) != 2 {
				t.Fatalf("completed=%v, err=%v, requests=%d", completed, err, len(client.requests))
			}
			messages := client.requests[1].Messages
			if got := messages[len(messages)-1].ExtractText(); got != tool.NotAvailableMsg {
				t.Fatalf("rejected tool result = %q", got)
			}
			if rejected.calls != 0 || len(runner.CollectPendingComments()) != 0 {
				t.Fatal("unadvertised tool executed or collected comments")
			}
			wantAllowed := 1
			if tc.noTools {
				wantAllowed = 0
			}
			if allowed.calls != wantAllowed {
				t.Fatalf("allowed tool calls = %d, want %d", allowed.calls, wantAllowed)
			}
		})
	}
}
