package create_unit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/unit"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Args struct {
	Name         string  `json:"name"`
	PluralName   *string `json:"plural_name,omitempty"`
	Description  *string `json:"description,omitempty"`
	BaseUnit     *string `json:"base_unit,omitempty"`
	OpenDataSlug *string `json:"open_data_slug,omitempty"`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "create_unit",
		Description: "Create a new unit of measurement.",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing create_unit")

		if args.Name == "" {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "Error: name is required"},
				},
				IsError: true,
			}, nil, nil
		}

		res, err := unit.CreateUnit(ctx, client, unit.UnitParam{
			Name:         args.Name,
			PluralName:   args.PluralName,
			Description:  args.Description,
			BaseUnit:     args.BaseUnit,
			OpenDataSlug: args.OpenDataSlug,
		})
		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error creating unit: %v", err)},
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
