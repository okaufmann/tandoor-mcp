package parse_ingredients

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	api_parser "github.com/compilercomplied/tandoor-mcp/src/tandoor/features/ingredient_parser"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Args struct {
	Ingredient  string   `json:"ingredient,omitempty"`
	Ingredients []string `json:"ingredients,omitempty"`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "parse_ingredients",
		Description: "Parse raw ingredient text strings into structured food, unit, and quantity objects.",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing parse_ingredients")

		if args.Ingredient == "" && len(args.Ingredients) == 0 {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "Error parsing ingredients: either ingredient or ingredients must be provided"},
				},
				IsError: true,
			}, nil, nil
		}

		res, err := api_parser.Parse(ctx, client, api_parser.IngredientParserRequest{
			Ingredient:  args.Ingredient,
			Ingredients: args.Ingredients,
		})

		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error parsing ingredients: %v", err)},
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
