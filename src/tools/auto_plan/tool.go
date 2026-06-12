package auto_plan

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	api_mealplan "github.com/compilercomplied/tandoor-mcp/src/tandoor/features/mealplan"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Args struct {
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	MealTypeID  int    `json:"meal_type_id"`
	Keywords    []int  `json:"keywords,omitempty"`
	KeywordMode *string `json:"keyword_mode,omitempty"`
	Servings    float64 `json:"servings"`
	AddShopping bool    `json:"add_shopping"`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "auto_plan",
		Description: "Automatically generate a meal plan over a range of dates.",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing auto_plan. start=%s, end=%s, meal_type_id=%d", args.StartDate, args.EndDate, args.MealTypeID)

		if args.StartDate == "" {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "Error auto-planning: start_date is required"},
				},
				IsError: true,
			}, nil, nil
		}
		if args.EndDate == "" {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "Error auto-planning: end_date is required"},
				},
				IsError: true,
			}, nil, nil
		}
		if args.MealTypeID <= 0 {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "Error auto-planning: meal_type_id must be greater than 0"},
				},
				IsError: true,
			}, nil, nil
		}
		if args.Servings <= 0 {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "Error auto-planning: servings must be greater than 0"},
				},
				IsError: true,
			}, nil, nil
		}

		startDate, err := time.Parse(time.RFC3339, args.StartDate)
		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error auto-planning: invalid start_date: %v", err)},
				},
				IsError: true,
			}, nil, nil
		}

		endDate, err := time.Parse(time.RFC3339, args.EndDate)
		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error auto-planning: invalid end_date: %v", err)},
				},
				IsError: true,
			}, nil, nil
		}

		res, err := api_mealplan.AutoPlan(ctx, client, api_mealplan.AutoMealPlanParam{
			StartDate:   startDate,
			EndDate:     endDate,
			MealTypeID:  args.MealTypeID,
			Keywords:    args.Keywords,
			KeywordMode: args.KeywordMode,
			Servings:    args.Servings,
			AddShopping: args.AddShopping,
		})

		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error auto-planning: %v", err)},
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
