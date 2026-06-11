package create_bookmarklet_import

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	api_bookmarklet "github.com/compilercomplied/tandoor-mcp/src/tandoor/features/bookmarklet_import"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Args struct {
	Url  *string `json:"url,omitempty"`
	Html string  `json:"html"`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "create_bookmarklet_import",
		Description: "Create a new bookmarklet import from a recipe URL and its HTML content.",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing create_bookmarklet_import")

		if args.Html == "" {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "Error creating bookmarklet import: html content is required"},
				},
				IsError: true,
			}, nil, nil
		}

		res, err := api_bookmarklet.Create(ctx, client, api_bookmarklet.BookmarkletImportParam{
			Url:  args.Url,
			Html: args.Html,
		})

		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error creating bookmarklet import: %v", err)},
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
