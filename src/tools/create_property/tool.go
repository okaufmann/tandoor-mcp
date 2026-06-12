package create_property

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/property"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Args struct {
	PropertyAmount float64 `json:"property_amount"`
	PropertyTypeID int     `json:"property_type_id"`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "create_property",
		Description: "Create a new food property relation (amount tied to a property type).",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing create_property")

		if args.PropertyAmount <= 0 {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "Error: property_amount must be greater than 0"},
				},
				IsError: true,
			}, nil, nil
		}
		if args.PropertyTypeID == 0 {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "Error: property_type_id is required"},
				},
				IsError: true,
			}, nil, nil
		}

		propType, err := property.GetPropertyType(ctx, client, args.PropertyTypeID)
		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error retrieving property type details: %v", err)},
				},
				IsError: true,
			}, nil, nil
		}

		res, err := property.CreateProperty(ctx, client, property.PropertyParam{
			PropertyAmount: &args.PropertyAmount,
			PropertyType: property.PropertyTypeRef{
				ID:   propType.ID,
				Name: propType.Name,
			},
		})
		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error creating property: %v", err)},
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
