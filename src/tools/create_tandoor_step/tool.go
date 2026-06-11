package create_tandoor_step

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	api_step "github.com/compilercomplied/tandoor-mcp/src/tandoor/features/step"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Args struct {
	RecipeID    int    `json:"recipe"`
	Name        string `json:"name,omitempty"`
	Instruction string `json:"instruction,omitempty"`
	Time        *int   `json:"time,omitempty"`
	Order       *int   `json:"order,omitempty"`
	Ingredients []int  `json:"ingredients,omitempty"`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "create_tandoor_step",
		Description: "Create a new cooking step for a recipe in Tandoor.",
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

		ingredients := args.Ingredients
		if ingredients == nil {
			ingredients = []int{}
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
