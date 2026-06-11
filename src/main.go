package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_ingredient"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_tandoor_recipe"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_tandoor_step"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_recipe_details"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_recipes"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	LoadEnv()

	transportFlag := flag.String("transport", "sse", "Transport to use: 'sse' (HTTP) or 'stdio'")
	hostFlag := flag.String("host", "0.0.0.0", "Host to listen on (only for SSE transport)")
	portFlag := flag.String("port", "8080", "Port to listen on (only for SSE transport)")
	flag.Parse()

	apiURL := GetEnv("TANDOOR_API_URL")
	apiToken := GetEnv("TANDOOR_API_TOKEN")

	client := tandoor.NewClient(apiURL, apiToken)

	server := mcp_sdk.NewServer(
		&mcp_sdk.Implementation{
			Name:    "tandoor-recipes-mcp",
			Version: "1.0.0",
		},
		nil,
	)

	get_recipes.Register(server, client)
	create_tandoor_recipe.Register(server, client)
	get_recipe_details.Register(server, client)
	create_tandoor_step.Register(server, client)
	create_ingredient.Register(server, client)

	switch *transportFlag {
	case "stdio":
		log.Println("Starting MCP server on stdio transport...")
		if err := server.Run(context.Background(), &mcp_sdk.StdioTransport{}); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	case "sse":
		addr := fmt.Sprintf("%s:%s", *hostFlag, *portFlag)
		log.Printf("Starting MCP server on SSE transport at http://%s ...", addr)
		log.Println("SSE Endpoint: /sse")
		log.Println("Message Endpoint: /sse?sessionid=<id> (automatically handled)")

		handler := mcp_sdk.NewSSEHandler(func(request *http.Request) *mcp_sdk.Server {
			return server
		}, nil)

		http.Handle("/sse", handler)
		http.Handle("/", handler)

		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	default:
		log.Fatalf("Unknown transport %q. Supported transports: stdio, sse", *transportFlag)
	}
}
