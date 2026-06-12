package create_food

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/food"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Args struct {
	Name           string  `json:"name"`
	PluralName     *string `json:"plural_name,omitempty"`
	Description    string  `json:"description,omitempty"`
	IgnoreShopping bool    `json:"ignore_shopping,omitempty"`
	ParentID       *int    `json:"parent_id,omitempty"`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "create_food",
		Description: "Create a new food item.",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing create_food")

		if args.Name == "" {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "Error: name is required"},
				},
				IsError: true,
			}, nil, nil
		}

		res, err := food.Create(ctx, client, food.FoodParam{
			Name:           args.Name,
			PluralName:     args.PluralName,
			Description:    args.Description,
			IgnoreShopping: args.IgnoreShopping,
			Parent:         args.ParentID,
		})
		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error creating food: %v", err)},
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
