package create_tandoor_step

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	api_ingredient "github.com/compilercomplied/tandoor-mcp/src/tandoor/features/ingredient"
	api_step "github.com/compilercomplied/tandoor-mcp/src/tandoor/features/step"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Args struct {
	RecipeID    int    `json:"recipe" jsonschema:"ID of the recipe to add the step to."`
	Name        string `json:"name,omitempty" jsonschema:"Optional name/header of the step."`
	Instruction string `json:"instruction,omitempty" jsonschema:"Instructions/actions for the step."`
	Time        *int   `json:"time,omitempty" jsonschema:"Optional time required in minutes."`
	Order       *int   `json:"order,omitempty" jsonschema:"Optional order number for sorting steps."`
	Ingredients []int  `json:"ingredients,omitempty" jsonschema:"List of existing ingredient IDs to associate with this step."`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "create_tandoor_step",
		Description: "Create a new cooking step for a recipe in Tandoor. Associate existing ingredients by passing their IDs.",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing create_tandoor_step. recipe=%v, instruction=%v", args.RecipeID, args.Instruction)

		if args.Instruction == "" && args.Name == "" {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "Error creating step: instruction or name is required"},
				},
				IsError: true,
			}, nil, nil
		}

		ingredients := make([]api_ingredient.IngredientResponse, len(args.Ingredients))
		for i, id := range args.Ingredients {
			ing, err := api_ingredient.Get(ctx, client, id)
			if err != nil {
				return &mcp_sdk.CallToolResult{
					Content: []mcp_sdk.Content{
						&mcp_sdk.TextContent{Text: fmt.Sprintf("Error fetching ingredient %d: %v", id, err)},
					},
					IsError: true,
				}, nil, nil
			}
			ingredients[i] = *ing
		}

		res, err := api_step.Create(ctx, client, args.RecipeID, api_step.StepParam{
			Name:        args.Name,
			Instruction: args.Instruction,
			Time:        args.Time,
			Order:       args.Order,
			Ingredients: ingredients,
		})

		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error creating step: %v", err)},
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

