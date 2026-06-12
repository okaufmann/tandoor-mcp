package get_inventory_entries

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/inventory"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Args struct {
	InventoryLocationID *int  `json:"inventory_location_id,omitempty"`
	FoodID              *int  `json:"food_id,omitempty"`
	Empty               *bool `json:"empty,omitempty"`
	Page                *int  `json:"page,omitempty"`
	PageSize            *int  `json:"page_size,omitempty"`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "get_inventory_entries",
		Description: "Retrieve a paginated list of inventory entries, optionally filtered by location, food, or empty status.",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing get_inventory_entries")

		res, err := inventory.ListEntries(ctx, client, inventory.ListEntriesParams{
			InventoryLocationID: args.InventoryLocationID,
			FoodID:              args.FoodID,
			Empty:               args.Empty,
			Page:                args.Page,
			PageSize:            args.PageSize,
		})
		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error listing inventory entries: %v", err)},
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
