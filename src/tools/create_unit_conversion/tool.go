package create_unit_conversion

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/food"
	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/unit"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Args struct {
	BaseAmount      float64 `json:"base_amount"`
	BaseUnitID      int     `json:"base_unit_id"`
	ConvertedAmount float64 `json:"converted_amount"`
	ConvertedUnitID int     `json:"converted_unit_id"`
	FoodID          *int    `json:"food_id,omitempty"`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "create_unit_conversion",
		Description: "Create a new unit conversion relation (optionally tied to a specific food).",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing create_unit_conversion")

		if args.BaseAmount <= 0 {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "Error: base_amount must be greater than 0"},
				},
				IsError: true,
			}, nil, nil
		}
		if args.BaseUnitID == 0 {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "Error: base_unit_id is required"},
				},
				IsError: true,
			}, nil, nil
		}
		if args.ConvertedAmount <= 0 {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "Error: converted_amount must be greater than 0"},
				},
				IsError: true,
			}, nil, nil
		}
		if args.ConvertedUnitID == 0 {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "Error: converted_unit_id is required"},
				},
				IsError: true,
			}, nil, nil
		}

		baseUnit, err := unit.GetUnit(ctx, client, args.BaseUnitID)
		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error retrieving base unit: %v", err)},
				},
				IsError: true,
			}, nil, nil
		}

		convertedUnit, err := unit.GetUnit(ctx, client, args.ConvertedUnitID)
		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error retrieving converted unit: %v", err)},
				},
				IsError: true,
			}, nil, nil
		}

		var foodRef *unit.FoodRef
		if args.FoodID != nil && *args.FoodID > 0 {
			f, err := food.Get(ctx, client, *args.FoodID)
			if err != nil {
				return &mcp_sdk.CallToolResult{
					Content: []mcp_sdk.Content{
						&mcp_sdk.TextContent{Text: fmt.Sprintf("Error retrieving food: %v", err)},
					},
					IsError: true,
				}, nil, nil
			}
			foodRef = &unit.FoodRef{
				ID:   f.ID,
				Name: f.Name,
			}
		}

		res, err := unit.CreateUnitConversion(ctx, client, unit.UnitConversionParam{
			BaseAmount: args.BaseAmount,
			BaseUnit: unit.UnitRef{
				ID:   baseUnit.ID,
				Name: baseUnit.Name,
			},
			ConvertedAmount: args.ConvertedAmount,
			ConvertedUnit: unit.UnitRef{
				ID:   convertedUnit.ID,
				Name: convertedUnit.Name,
			},
			Food: foodRef,
		})
		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error creating unit conversion: %v", err)},
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
