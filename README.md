# Tandoor Recipes MCP Server (Go)

A Model Context Protocol (MCP) server written in Go that integrates with [Tandoor Recipes](https://tandoor.dev/). This server allows LLMs to query recipes in your self-hosted Tandoor instance.

## Features
The Tandoor MCP server currently supports the following features:
- ✅ **Get Recipes**: Retrieve a list of recipes from Tandoor (`get_recipes`).
- ✅ **Create Recipe**: Create a new recipe in Tandoor (`create_tandoor_recipe`).
- ✅ **Get Recipe Details**: Retrieve full details of a specific recipe (`get_recipe_details`).
- ✅ **Steps**: Manage recipe steps (`create_tandoor_step`).
- ✅ **Ingredients**: Manage ingredients (`create_ingredient`).
- ✅ **Cook Logs**: Track cooking history (`create_cook_log`, `get_cook_logs`, `create_view_log`, `get_view_logs`).
- ✅ **Imports**: Manage recipe imports (`create_recipe_import`, `get_recipe_imports`, `create_bookmarklet_import`, `get_bookmarklet_imports`, `parse_ingredients`).
- ✅ **Meal Plans**: Get existing meal plans or create new ones (`create_meal_plan`, `get_meal_plans`, `auto_plan`, `create_meal_type`, `get_meal_types`).

### Planned Features (Unimplemented)
**Meal Planning & Shopping**
- ❌ **Shopping List**: Manage shopping lists and entries (`shopping-list`, `shopping-list-entry`, `shopping-list-recipe`).
- ❌ **Supermarkets**: Manage supermarkets and categories (`supermarket`, `supermarket-category`, `supermarket-category-relation`).

**Inventory & Storage**
- ❌ **Storage**: Manage storage locations (`storage`).
- ❌ **Inventory**: Manage inventory locations, entries, and logs (`inventory-location`, `inventory-entry`, `inventory-log`).

**Food & Metadata**
- ❌ **Foods**: Query available foods and inheritance fields (`food`, `food-inherit-field`).
- ❌ **Keywords & Meal Types**: Query available keywords and meal types (`keyword`, `meal-type`).
- ❌ **Units**: Query available units and conversions (`unit`, `unit-conversion`).
- ❌ **Properties**: Manage food properties (`property`, `property-type`).

**Collections & Sharing**
- ❌ **Recipe Books**: Manage recipe books and entries (`recipe-book`, `recipe-book-entry`).
- ❌ **Invites**: Manage invite links (`invite-link`).

**System & Users**
- ❌ **Spaces**: Manage spaces and household configuration (`space`, `user-space`, `household`, `group`).
- ❌ **Users**: Manage users, preferences, and files (`user`, `user-preference`, `user-file`).
- ❌ **Search Config**: Manage search preferences and fields (`search-fields`, `search-preference`, `custom-filter`).

**Integrations & Admin**
- ❌ **Integrations**: Manage syncs, connectors, and automation (`sync`, `sync-log`, `connector-config`, `automation`).
- ❌ **Logs**: Access system logs (`import-log`, `export-log`).
- ❌ **AI & Tokens**: Manage AI features and API tokens (`ai-provider`, `ai-log`, `access-token`).
- ❌ **Server**: Manage server settings and localization (`server-settings`, `localization`).

## Supported Transports
- **stdio**: Standard input/output transport for the Claude Desktop integration.
- **sse**: Server-Sent Events transport via HTTP, useful for E2E testing and standalone debugging.

## How to Configure
You need to set up the following environment variables. You can add them to a `.env` file in the root directory:
```
TANDOOR_API_URL=http://localhost:8080/api
TANDOOR_API_TOKEN=your_token_here
TANDOOR_API_SPACE_ID=1
```

## How to Launch
To run the server locally:
```bash
go run ./src
```
To run via Docker Compose (includes a Tandoor instance and E2E tests):
```bash
docker-compose up --build tandoor-db tandoor-web tandoor-setup mcp-api -d
docker-compose run e2e-tests
```

## Getting Started

### Prerequisites

- Go 1.22 or higher.
- A running Tandoor Recipes instance.
- An API token from your Tandoor instance (find/generate this under **Settings > API** in Tandoor).

### Configuration (.env)

The server requires your Tandoor instance URL and API token. Create a `.env` file in the root of the project (or export them as system environment variables):

```env
TANDOOR_API_URL=https://your-tandoor-instance.com
TANDOOR_API_TOKEN=your_tandoor_api_token
```

### Running via Docker Compose (Recommended)

You can launch both the API and the E2E tests in fully isolated containers using Docker Compose. The `docker-compose.yml` automatically reads your `.env` file.

**Launch the API locally on port 8080 (SSE Transport):**
```bash
docker-compose up mcp-api
```

**Run the Black-Box E2E Tests:**
```bash
docker-compose up --build --exit-code-from e2e-tests
```
*(This boots the API, waits for health checks, executes the E2E tests against the container network, and exits with the test result).*

### Setup & Build (Manual / Native)

1. Clone this repository.
2. Ensure you have a `.env` file present.
3. Build and run the binary natively:
   ```bash
   go run ./src/ -transport sse -port 8080
   ```

### Running the Server Natively

#### Option A: Standalone HTTP (SSE) Server

Start the server natively as an HTTP service listening on a specific port (default is `8080`):

```bash
./tandoor-mcp -transport sse -port 8080
```

The server will expose:
- SSE connection endpoint: `http://localhost:8080/sse`
- Client message endpoint: `http://localhost:8080/sse?sessionid=<id>` (managed automatically by the transport handler)

#### Option B: Stdio Server (for Claude Desktop / Cursor)

If you are running the server locally to be consumed by Claude Desktop or Cursor, configure it to run over standard input/output (`stdio`):

```bash
./tandoor-mcp -transport stdio
```

### Configuration Examples

#### 1. Claude Desktop Configuration (`mcpServers` section)

Add the following to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "tandoor-recipes": {
      "command": "/absolute/path/to/tandoor-mcp",
      "args": ["-transport", "stdio"],
      "env": {
        "TANDOOR_API_URL": "https://your-tandoor-instance.com",
        "TANDOOR_API_TOKEN": "your_tandoor_api_token"
      }
    }
  }
}
```

#### 2. SSE Server Configuration

If your client supports HTTP SSE, run the server in `sse` mode and add the server definition to your client configuration:

```json
{
  "mcpServers": {
    "tandoor-recipes-sse": {
      "url": "http://localhost:8080/sse"
    }
  }
}
```
