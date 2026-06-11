package create_ingredient

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	api_ingredient "github.com/compilercomplied/tandoor-mcp/src/tandoor/features/ingredient"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Args are the MCP tool arguments for create_ingredient.
type Args struct {
	FoodName string  `json:"food_name"`
	UnitName string  `json:"unit_name"`
	Amount   float64 `json:"amount"`
	Note     string  `json:"note,omitempty"`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "create_ingredient",
		Description: "Create a new ingredient in Tandoor (POST /api/ingredient/). Food and unit are created automatically if they do not yet exist. Returns the created ingredient.",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing create_ingredient. food=%q, unit=%q, amount=%v", args.FoodName, args.UnitName, args.Amount)

		if args.FoodName == "" {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "validation error: food_name is required"},
				},
				IsError: true,
			}, nil, nil
		}
		if args.UnitName == "" {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "validation error: unit_name is required"},
				},
				IsError: true,
			}, nil, nil
		}
		if args.Amount <= 0 {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "validation error: amount must be greater than 0"},
				},
				IsError: true,
			}, nil, nil
		}

		res, err := api_ingredient.Create(ctx, client, api_ingredient.IngredientParam{
			Food:   api_ingredient.FoodRef{Name: args.FoodName},
			Unit:   api_ingredient.UnitRef{Name: args.UnitName},
			Amount: args.Amount,
			Note:   args.Note,
		})
		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error creating ingredient: %v", err)},
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
