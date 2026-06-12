package remove_shopping_list_item

import (
	"context"
	"fmt"
	"log"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/shoppinglist"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Args struct {
	ItemID int `json:"item_id"`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "remove_shopping_list_item",
		Description: "Remove an item from the shopping list.",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing remove_shopping_list_item for item ID %d", args.ItemID)

		if args.ItemID == 0 {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "Error: item_id is required"},
				},
				IsError: true,
			}, nil, nil
		}

		err := shoppinglist.DeleteEntry(ctx, client, args.ItemID)
		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error removing shopping list item: %v", err)},
				},
				IsError: true,
			}, nil, nil
		}

		return &mcp_sdk.CallToolResult{
			Content: []mcp_sdk.Content{
				&mcp_sdk.TextContent{Text: fmt.Sprintf("Successfully removed shopping list item %d", args.ItemID)},
			},
		}, nil, nil
	})
}
