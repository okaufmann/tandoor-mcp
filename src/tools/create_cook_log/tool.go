package create_cook_log

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	api_cooklog "github.com/compilercomplied/tandoor-mcp/src/tandoor/features/cooklog"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Args struct {
	Recipe    int     `json:"recipe"`
	Servings  *int    `json:"servings,omitempty"`
	Rating    *int    `json:"rating,omitempty"`
	Comment   *string `json:"comment,omitempty"`
	CreatedAt *string `json:"created_at,omitempty"`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "create_cook_log",
		Description: "Log that a recipe was cooked, recording servings, rating, comment, and date/time.",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing create_cook_log. recipe=%v", args.Recipe)

		if args.Recipe <= 0 {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "Error creating cook log: recipe is required and must be greater than 0"},
				},
				IsError: true,
			}, nil, nil
		}

		var parsedTime *time.Time
		if args.CreatedAt != nil && *args.CreatedAt != "" {
			t, err := time.Parse(time.RFC3339, *args.CreatedAt)
			if err != nil {
				return &mcp_sdk.CallToolResult{
					Content: []mcp_sdk.Content{
						&mcp_sdk.TextContent{Text: fmt.Sprintf("Error creating cook log: invalid created_at date-time format (expected RFC3339/ISO-8601, e.g. 2026-06-11T21:00:00Z): %v", err)},
					},
					IsError: true,
				}, nil, nil
			}
			parsedTime = &t
		}

		res, err := api_cooklog.Create(ctx, client, api_cooklog.CookLogParam{
			Recipe:    args.Recipe,
			Servings:  args.Servings,
			Rating:    args.Rating,
			Comment:   args.Comment,
			CreatedAt: parsedTime,
		})

		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error creating cook log: %v", err)},
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
