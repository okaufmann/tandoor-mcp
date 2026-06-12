package update_shopping_list_item

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/shoppinglist"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Args struct {
	ItemID  int     `json:"item_id"`
	Amount  *string `json:"amount,omitempty"`
	Checked *bool   `json:"checked,omitempty"`
	UnitID  *int    `json:"unit_id,omitempty"`
	Note    *string `json:"note,omitempty"`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "update_shopping_list_item",
		Description: "Update an existing shopping list item (e.g., check/uncheck, change amount).",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing update_shopping_list_item for item ID %d", args.ItemID)

		if args.ItemID == 0 {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "Error: item_id is required"},
				},
				IsError: true,
			}, nil, nil
		}

		params := shoppinglist.UpdateEntryParam{
			Checked: args.Checked,
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

		if args.UnitID != nil {
			params.Unit = &shoppinglist.UnitResponse{
				ID: *args.UnitID,
			}
		}

		res, err := shoppinglist.PatchEntry(ctx, client, args.ItemID, params)
		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error updating shopping list item: %v", err)},
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
