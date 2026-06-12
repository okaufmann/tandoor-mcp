package create_meal_plan

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
	Title        *string  `json:"title,omitempty"`
	RecipeID     *int     `json:"recipe_id,omitempty"`
	RecipeName   *string  `json:"recipe_name,omitempty"`
	Servings     float64  `json:"servings"`
	Note         *string  `json:"note,omitempty"`
	FromDate     string   `json:"from_date"`
	ToDate       *string  `json:"to_date,omitempty"`
	MealTypeID   int      `json:"meal_type_id"`
	MealTypeName string   `json:"meal_type_name"`
	AddShopping  *bool    `json:"add_shopping,omitempty"`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "create_meal_plan",
		Description: "Create a new meal plan entry for a recipe or custom note on a given date/time.",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing create_meal_plan. from_date=%s, servings=%v", args.FromDate, args.Servings)

		if args.FromDate == "" {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "Error creating meal plan: from_date is required"},
				},
				IsError: true,
			}, nil, nil
		}
		if args.MealTypeID <= 0 {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "Error creating meal plan: meal_type_id must be greater than 0"},
				},
				IsError: true,
			}, nil, nil
		}
		if args.MealTypeName == "" {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "Error creating meal plan: meal_type_name is required"},
				},
				IsError: true,
			}, nil, nil
		}
		if args.Servings <= 0 {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "Error creating meal plan: servings must be greater than 0"},
				},
				IsError: true,
			}, nil, nil
		}

		fromDate, err := time.Parse(time.RFC3339, args.FromDate)
		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error creating meal plan: invalid from_date format: %v", err)},
				},
				IsError: true,
			}, nil, nil
		}

		var toDate *time.Time
		if args.ToDate != nil && *args.ToDate != "" {
			t, err := time.Parse(time.RFC3339, *args.ToDate)
			if err != nil {
				return &mcp_sdk.CallToolResult{
					Content: []mcp_sdk.Content{
						&mcp_sdk.TextContent{Text: fmt.Sprintf("Error creating meal plan: invalid to_date format: %v", err)},
					},
					IsError: true,
				}, nil, nil
			}
			toDate = &t
		}

		var recipeParam *api_mealplan.RecipeOverviewParam
		if args.RecipeID != nil {
			recipeName := ""
			if args.RecipeName != nil {
				recipeName = *args.RecipeName
			}
			recipeParam = &api_mealplan.RecipeOverviewParam{
				ID:   *args.RecipeID,
				Name: recipeName,
			}
		}

		res, err := api_mealplan.Create(ctx, client, api_mealplan.MealPlanParam{
			Title:       args.Title,
			Recipe:      recipeParam,
			Servings:    args.Servings,
			Note:        args.Note,
			FromDate:    fromDate,
			ToDate:      toDate,
			AddShopping: args.AddShopping,
			MealType: api_mealplan.MealTypeParam{
				ID:   args.MealTypeID,
				Name: args.MealTypeName,
			},
		})

		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error creating meal plan: %v", err)},
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
