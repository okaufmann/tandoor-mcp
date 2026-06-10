package get_recipes

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	api_get_recipes "github.com/compilercomplied/tandoor-mcp/src/tandoor/features/get_recipes"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Args struct {
	Query    string `json:"query,omitempty" jsonschema:"Search term for recipe names."`
	Foods    []int  `json:"foods,omitempty" jsonschema:"Array of Food IDs (match ANY)."`
	Keywords []int  `json:"keywords,omitempty" jsonschema:"Array of Keyword IDs (match ANY)."`
	Limit    *int   `json:"limit,omitempty" jsonschema:"Max number of recipes to return (default: 10)."`
	Rating   *int   `json:"rating,omitempty" jsonschema:"Minimum rating (0-5)."`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "get_recipes",
		Description: "Search for recipes in Tandoor based on various criteria.",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing get_recipes. query=%q, foods=%v, keywords=%v, limit=%v, rating=%v",
			args.Query, args.Foods, args.Keywords, args.Limit, args.Rating)

		res, err := api_get_recipes.Search(ctx, client, api_get_recipes.GetRecipesParams{
			Query:    args.Query,
			Search:   args.Query,
			Foods:    args.Foods,
			Keywords: args.Keywords,
			Limit:    args.Limit,
			Rating:   args.Rating,
		})
		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error fetching recipes: %v", err)},
				},
				IsError: true,
			}, nil, nil
		}

		if len(res.Results) == 0 {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "No recipes found matching the criteria."},
				},
			}, nil, nil
		}

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Found %d recipes:\n\n", len(res.Results)))
		sb.WriteString("| ID | Name | Description | Prep Time | Cook Time | Servings |\n")
		sb.WriteString("|---|---|---|---|---|---|\n")
		for _, r := range res.Results {
			desc := r.Description
			desc = strings.ReplaceAll(desc, "\r\n", " ")
			desc = strings.ReplaceAll(desc, "\n", " ")
			if len(desc) > 100 {
				desc = desc[:97] + "..."
			}
			prep := "-"
			if r.PrepTime != nil {
				prep = fmt.Sprintf("%d mins", *r.PrepTime)
			}
			cook := "-"
			if r.CookTime != nil {
				cook = fmt.Sprintf("%d mins", *r.CookTime)
			}
			serv := "-"
			if r.Servings != nil {
				serv = fmt.Sprintf("%d", *r.Servings)
			}
			sb.WriteString(fmt.Sprintf("| %d | %s | %s | %s | %s | %s |\n", r.ID, r.Name, desc, prep, cook, serv))
		}

		return &mcp_sdk.CallToolResult{
			Content: []mcp_sdk.Content{
				&mcp_sdk.TextContent{Text: sb.String()},
			},
		}, nil, nil
	})
}
