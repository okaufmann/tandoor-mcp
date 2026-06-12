package get_shopping_list

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/shoppinglist"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Args struct {
	Checked string `json:"checked,omitempty"` // "true", "false", "both", "recent"
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "get_shopping_list",
		Description: "Retrieve the current shopping list items.",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing get_shopping_list with checked filter: %s", args.Checked)

		pageSize := 100
		res, err := shoppinglist.ListEntries(ctx, client, shoppinglist.ListEntriesParams{
			PageSize: &pageSize,
		})
		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error retrieving shopping list items: %v", err)},
				},
				IsError: true,
			}, nil, nil
		}

		filter := args.Checked
		if filter == "" {
			filter = "recent" // default
		}

		var filtered []shoppinglist.ShoppingListEntryResponse
		for _, entry := range res.Results {
			switch filter {
			case "true":
				if entry.Checked {
					filtered = append(filtered, entry)
				}
			case "false":
				if !entry.Checked {
					filtered = append(filtered, entry)
				}
			case "both":
				filtered = append(filtered, entry)
			case "recent":
				if !entry.Checked {
					filtered = append(filtered, entry)
				}
			default:
				filtered = append(filtered, entry)
			}
		}

		b, _ := json.MarshalIndent(filtered, "", "  ")
		return &mcp_sdk.CallToolResult{
			Content: []mcp_sdk.Content{
				&mcp_sdk.TextContent{Text: string(b)},
			},
		}, nil, nil
	})
}
