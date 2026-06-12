package create_meal_type

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	api_mealtype "github.com/compilercomplied/tandoor-mcp/src/tandoor/features/mealtype"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Args struct {
	Name  string  `json:"name"`
	Order *int    `json:"order,omitempty"`
	Time  *string `json:"time,omitempty"`
	Color *string `json:"color,omitempty"`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "create_meal_type",
		Description: "Create a new meal type (e.g. Lunch, Dinner, Breakfast) in Tandoor.",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing create_meal_type. name=%s", args.Name)

		if args.Name == "" {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "Error creating meal type: name is required"},
				},
				IsError: true,
			}, nil, nil
		}

		res, err := api_mealtype.Create(ctx, client, api_mealtype.MealTypeParam{
			Name:  args.Name,
			Order: args.Order,
			Time:  args.Time,
			Color: args.Color,
		})

		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error creating meal type: %v", err)},
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
