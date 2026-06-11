package create_recipe_import

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	api_recipe_import "github.com/compilercomplied/tandoor-mcp/src/tandoor/features/recipe_import"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Args struct {
	StorageName string  `json:"storage_name"`
	StoragePath string  `json:"storage_path,omitempty"`
	ImportName  string  `json:"import_name"`
	FileUID     *string `json:"file_uid,omitempty"`
	FilePath    *string `json:"file_path,omitempty"`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "create_recipe_import",
		Description: "Create a new recipe import job associated with a storage location configuration.",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing create_recipe_import. import_name=%s, storage_name=%s", args.ImportName, args.StorageName)

		if args.StorageName == "" {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "Error creating recipe import: storage_name is required"},
				},
				IsError: true,
			}, nil, nil
		}
		if args.ImportName == "" {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "Error creating recipe import: import_name is required"},
				},
				IsError: true,
			}, nil, nil
		}

		var fileUIDVal string
		if args.FileUID != nil {
			fileUIDVal = *args.FileUID
		}
		var filePathVal string
		if args.FilePath != nil {
			filePathVal = *args.FilePath
		}

		res, err := api_recipe_import.Create(ctx, client, api_recipe_import.RecipeImportParam{
			Name:     args.ImportName,
			FileUID:  fileUIDVal,
			FilePath: filePathVal,
			Storage: api_recipe_import.StorageParam{
				Name: args.StorageName,
				Path: args.StoragePath,
			},
		})

		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error creating recipe import: %v", err)},
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
