package update_inventory_entry

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/inventory"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Args struct {
	EntryID             int     `json:"entry_id"`
	Amount              *string `json:"amount,omitempty"`
	InventoryLocationID *int    `json:"inventory_location_id,omitempty"`
	Note                *string `json:"note,omitempty"`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "update_inventory_entry",
		Description: "Update an existing inventory entry.",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing update_inventory_entry for entry ID %d", args.EntryID)

		if args.EntryID == 0 {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "Error: entry_id is required"},
				},
				IsError: true,
			}, nil, nil
		}

		params := inventory.InventoryEntryUpdateParam{
			Note: args.Note,
		}

		if args.Amount != nil {
			amt, err := parseAmount(*args.Amount)
			if err != nil {
				return &mcp_sdk.CallToolResult{
					Content: []mcp_sdk.Content{
						&mcp_sdk.TextContent{Text: fmt.Sprintf("Error parsing amount: %v", err)},
					},
					IsError: true,
				}, nil, nil
			}
			params.Amount = &amt
		}

		if args.InventoryLocationID != nil {
			loc, err := inventory.GetLocation(ctx, client, *args.InventoryLocationID)
			if err != nil {
				return &mcp_sdk.CallToolResult{
					Content: []mcp_sdk.Content{
						&mcp_sdk.TextContent{Text: fmt.Sprintf("Error retrieving inventory location: %v", err)},
					},
					IsError: true,
				}, nil, nil
			}
			params.InventoryLocation = loc
		}

		res, err := inventory.PatchEntry(ctx, client, args.EntryID, params)
		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error updating inventory entry: %v", err)},
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

func parseAmount(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("amount is empty")
	}
	if val, err := strconv.ParseFloat(s, 64); err == nil {
		return val, nil
	}
	if parts := strings.Split(s, "/"); len(parts) == 2 {
		num, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		den, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if err1 == nil && err2 == nil && den != 0 {
			return num / den, nil
		}
	}
	return 0, fmt.Errorf("invalid amount format: %q", s)
}
