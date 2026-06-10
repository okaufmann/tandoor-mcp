package get_recipe_details

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/get_recipe_details"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Args struct {
	RecipeID int `json:"recipe_id" jsonschema:"The unique ID of the recipe to retrieve"`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "get_recipe_details",
		Description: "Retrieve full details of a specific recipe from Tandoor by its ID.",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing get_recipe_details. id=%d", args.RecipeID)

		res, err := get_recipe_details.Do(ctx, client, args.RecipeID)

		if err != nil {
			log.Printf("get_recipe_details error: %v", err)
			return &mcp_sdk.CallToolResult{
				IsError: true,
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("API error: %v", err)},
				},
			}, nil, nil
		}

		resBytes, _ := json.MarshalIndent(res, "", "  ")

		return &mcp_sdk.CallToolResult{
			Content: []mcp_sdk.Content{
				&mcp_sdk.TextContent{Text: string(resBytes)},
			},
		}, nil, nil
	})
}
