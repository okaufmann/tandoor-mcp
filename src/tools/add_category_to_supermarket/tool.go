package add_category_to_supermarket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/supermarket"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Args struct {
	SupermarketID int `json:"supermarket_id"`
	CategoryID    int `json:"category_id"`
	Order         int `json:"order"`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "add_category_to_supermarket",
		Description: "Link a supermarket category to a supermarket with a specific order.",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing add_category_to_supermarket")

		if args.SupermarketID == 0 {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "Error: supermarket_id is required"},
				},
				IsError: true,
			}, nil, nil
		}
		if args.CategoryID == 0 {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "Error: category_id is required"},
				},
				IsError: true,
			}, nil, nil
		}

		category, err := supermarket.GetCategory(ctx, client, args.CategoryID)
		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error retrieving category details: %v", err)},
				},
				IsError: true,
			}, nil, nil
		}

		res, err := supermarket.CreateRelation(ctx, client, supermarket.SupermarketCategoryRelationParam{
			Supermarket: args.SupermarketID,
			Category: supermarket.CategoryRef{
				ID:   args.CategoryID,
				Name: category.Name,
			},
			Order: args.Order,
		})
		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error adding category to supermarket: %v", err)},
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
