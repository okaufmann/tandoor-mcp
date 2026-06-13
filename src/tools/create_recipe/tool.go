package create_recipe

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	api_create_recipe "github.com/compilercomplied/tandoor-mcp/src/tandoor/features/create_recipe"
	api_ingredient "github.com/compilercomplied/tandoor-mcp/src/tandoor/features/ingredient"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Args struct {
	Name        string      `json:"name" jsonschema:"Name of the recipe"`
	Description string      `json:"description,omitempty" jsonschema:"Description of the recipe"`
	Servings    *int        `json:"servings,omitempty" jsonschema:"Number of servings"`
	WorkingTime *int        `json:"working_time,omitempty" jsonschema:"Working time in minutes"`
	WaitingTime *int        `json:"waiting_time,omitempty" jsonschema:"Waiting time in minutes"`
	Steps       []StepParam `json:"steps,omitempty" jsonschema:"Steps to create with the recipe"`
}

type StepParam struct {
	Name        string `json:"name,omitempty" jsonschema:"Optional name/header of the step."`
	Instruction string `json:"instruction,omitempty" jsonschema:"Optional action instruction text."`
	Time        *int   `json:"time,omitempty" jsonschema:"Optional time in minutes."`
	Order       *int   `json:"order,omitempty" jsonschema:"Optional order."`
	Ingredients []int  `json:"ingredients,omitempty" jsonschema:"Optional list of existing ingredient IDs to associate with this step."`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "create_recipe",
		Description: "Create a complete recipe in Tandoor including steps and associated ingredient IDs in a single request.",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing create_recipe. name=%q", args.Name)

		if args.Name == "" {
			return &mcp_sdk.CallToolResult{
				IsError: true,
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "name is required"},
				},
			}, nil, nil
		}

		steps := make([]api_create_recipe.StepParam, len(args.Steps))
		for i, stepArg := range args.Steps {
			ingredients := make([]api_ingredient.IngredientResponse, len(stepArg.Ingredients))
			for j, id := range stepArg.Ingredients {
				ing, err := api_ingredient.Get(ctx, client, id)
				if err != nil {
					return &mcp_sdk.CallToolResult{
						IsError: true,
						Content: []mcp_sdk.Content{
							&mcp_sdk.TextContent{Text: fmt.Sprintf("Error fetching ingredient %d: %v", id, err)},
						},
					}, nil, nil
				}
				ingredients[j] = *ing
			}

			steps[i] = api_create_recipe.StepParam{
				Name:        stepArg.Name,
				Instruction: stepArg.Instruction,
				Time:        stepArg.Time,
				Order:       stepArg.Order,
				Ingredients: ingredients,
			}
		}

		res, err := api_create_recipe.Create(ctx, client, api_create_recipe.CreateRecipeParams{
			Name:        args.Name,
			Description: args.Description,
			Servings:    args.Servings,
			WorkingTime: args.WorkingTime,
			WaitingTime: args.WaitingTime,
			Steps:       steps,
		})

		if err != nil {
			log.Printf("create_recipe error: %v", err)
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
