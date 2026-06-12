package get_meal_plans

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	api_mealplan "github.com/compilercomplied/tandoor-mcp/src/tandoor/features/mealplan"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Args struct {
	FromDate    *string `json:"from_date,omitempty"` // YYYY-MM-DD
	ToDate      *string `json:"to_date,omitempty"`   // YYYY-MM-DD
	MealTypeIDs []int   `json:"meal_type_ids,omitempty"`
	Page        *int    `json:"page,omitempty"`
	PageSize    *int    `json:"page_size,omitempty"`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "get_meal_plans",
		Description: "Get a list of planned meals, optionally filtered by dates and meal type IDs.",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing get_meal_plans")

		res, err := api_mealplan.List(ctx, client, api_mealplan.ListParams{
			FromDate: args.FromDate,
			ToDate:   args.ToDate,
			MealType: args.MealTypeIDs,
			Page:     args.Page,
			PageSize: args.PageSize,
		})

		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error listing meal plans: %v", err)},
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
