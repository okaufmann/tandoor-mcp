package get_keywords

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/keyword"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Args struct {
	Query    *string `json:"query,omitempty"`
	Root     *int    `json:"root,omitempty"`
	Tree     *int    `json:"tree,omitempty"`
	Page     *int    `json:"page,omitempty"`
	PageSize *int    `json:"page_size,omitempty"`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "get_keywords",
		Description: "List or search for keywords.",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing get_keywords")

		res, err := keyword.List(ctx, client, keyword.ListParams{
			Query:    args.Query,
			Root:     args.Root,
			Tree:     args.Tree,
			Page:     args.Page,
			PageSize: args.PageSize,
		})
		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error listing keywords: %v", err)},
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
