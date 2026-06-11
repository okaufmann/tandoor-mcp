package get_cook_logs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	api_cooklog "github.com/compilercomplied/tandoor-mcp/src/tandoor/features/cooklog"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Args struct {
	Recipe   *int `json:"recipe,omitempty"`
	Page     *int `json:"page,omitempty"`
	PageSize *int `json:"page_size,omitempty"`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "get_cook_logs",
		Description: "Get a list of logged cooking events, optionally filtered by recipe.",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing get_cook_logs")

		res, err := api_cooklog.List(ctx, client, api_cooklog.ListParams{
			Recipe:   args.Recipe,
			Page:     args.Page,
			PageSize: args.PageSize,
		})

		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error listing cook logs: %v", err)},
				},
				IsError: true,
			}, nil, nil
		}

		b, _ := json.MarshalIndent(res, "", "  ")

		return &mcp_sdk.CallToolResult{
			Content: []mcp_sdk.Content{
				&mcp_sdk.TextContent{Text: string(b)},
			},
		}, nil, nil
	})
}
