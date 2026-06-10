package create_tandoor_recipe

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	api_create_recipe "github.com/compilercomplied/tandoor-mcp/src/tandoor/features/create_recipe"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Args struct {
	Name        string `json:"name" jsonschema:"Name of the recipe"`
	Description string `json:"description,omitempty" jsonschema:"Description of the recipe"`
	Servings    *int   `json:"servings,omitempty" jsonschema:"Number of servings"`
	WorkingTime *int   `json:"working_time,omitempty" jsonschema:"Working time in minutes"`
	WaitingTime *int   `json:"waiting_time,omitempty" jsonschema:"Waiting time in minutes"`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "create_tandoor_recipe",
		Description: "Create a new recipe in Tandoor. Returns the created recipe ID and details.",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing create_tandoor_recipe. name=%q", args.Name)

		if args.Name == "" {
			return &mcp_sdk.CallToolResult{
				IsError: true,
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "name is required"},
				},
			}, nil, nil
		}

		res, err := api_create_recipe.Create(ctx, client, api_create_recipe.CreateRecipeParams{
			Name:        args.Name,
			Description: args.Description,
			Servings:    args.Servings,
			WorkingTime: args.WorkingTime,
			WaitingTime: args.WaitingTime,
			Steps:       []api_create_recipe.StepParam{},
		})

		if err != nil {
			log.Printf("create_tandoor_recipe error: %v", err)
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
