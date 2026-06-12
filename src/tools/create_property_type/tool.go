package create_property_type

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
	Name         string  `json:"name"`
	Unit         *string `json:"unit,omitempty"`
	Description  *string `json:"description,omitempty"`
	Order        *int    `json:"order,omitempty"`
	OpenDataSlug *string `json:"open_data_slug,omitempty"`
	FdcID        *int    `json:"fdc_id,omitempty"`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "create_property_type",
		Description: "Create a new property type.",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing create_property_type")

		if args.Name == "" {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "Error: name is required"},
				},
				IsError: true,
			}, nil, nil
		}

		orderVal := 0
		if args.Order != nil {
			orderVal = *args.Order
		}

		res, err := property.CreatePropertyType(ctx, client, property.PropertyTypeParam{
			Name:         args.Name,
			Unit:         args.Unit,
			Description:  args.Description,
			Order:        orderVal,
			OpenDataSlug: args.OpenDataSlug,
			FdcID:        args.FdcID,
		})
		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error creating property type: %v", err)},
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
