# Development Loop
Every feature or bug fix implementation must follow this strict development loop:
- **Implementation**: Write the feature code or fix following the architecture patterns (vertical slicing, thin client, generic abstractions).
- **Testing**: All E2E tests must pass successfully using the `mise run test-e2e-docker` pipeline. A new E2E test **must** be added whenever changing behavior or adding a new feature.
- **Documentation**: Update `README.md` if necessary (e.g., when adding new features or modifying the setup process).

# Architecture & Design Decisions

This document outlines the core architectural patterns and design decisions for the Tandoor MCP project. Agents modifying or extending this codebase must adhere strictly to these principles.

## 1. Code Splitting & Vertical Slicing
We utilize a **Vertical Slicing** pattern for feature implementation to minimize cross-domain concerns. 
- Features are encapsulated in directories under `src/tandoor/features/<feature_name>`.
- **`dto.go`**: Contains the Data Transfer Objects (Request/Response models) strictly specific to that feature.
- **`route.go`**: Contains the client method and feature-specific business logic (e.g., query building).
- We maintain a strict boundary between the internal Tandoor API features and the external MCP tools. MCP tools exist in `src/tools/<tool_name>/tool.go` and act as lightweight glue layers that register with the SDK and forward traffic to the corresponding internal API features.

## 2. Tandoor HTTP Client Structure
The HTTP Client acts as a very thin, generic envelope without complex boilerplate or reflection.
- **`tandoor.Client`**: Holds solely the `http.Client`, the Base URL, and injects authentication headers (via `Bearer` token middleware inside `authTransport`). It handles **no business logic**.
- **Generics over Reflection**: Responses are deserialized using a top-level generic function: `func Request[T any](ctx context.Context, c *Client, method, path string, query url.Values, body any) (*T, error)`.
- **Query Building**: A generic `QueryBuilder` struct exists to fluently chain URL parameters (`qb.Add("key", value)`).
- **Error Sinking**: Errors from the API are captured and handled dynamically in a single dedicated sink (`checkResponse`) inside `client.go`, leaving the calling methods extremely clean and readable.

## 3. Black-Box End-to-End (E2E) Testing
E2E tests in this project must be strictly **Black Box**, meaning zero mocks are permitted for external APIs or internal structures.
- **Global Package State**: Tests utilize a generalized `TestMain` function in `main_test.go` to ensure `SetupFixture` and `Teardown` happen exactly once per package execution, making testing lightning fast.
- **Native Booting**: The fixture natively compiles the project into a binary and executes it dynamically with SSE transport over a free port (`8081`).
- **Real World Replication**: The E2E tests connect via an HTTP SSE transport exactly as Claude Desktop would. The compiled server securely reads `.env` variables (e.g. `TANDOOR_API_URL` and `TANDOOR_API_TOKEN`) to hit the *real* live Tandoor instance. 
- **Generic MCP Call Tool**: The `infra.Client` abstracts the underlying raw map arguments allowing tests to trigger MCP operations completely generically using `infra.CallTool[T](...)`.
- **AAA Pattern**: All tests must strictly adhere to the Arrange, Act, Assert pattern.

## 4. Configuration
- Configuration logic is isolated to `src/config.go`.
- `LoadEnv()` natively parses the `.env` file into `os.Setenv` values.
- `GetEnv(key)` ensures fast failing at startup by panicking aggressively if environment values are missing. No silent failures.

## 5. Reusing Established Abstractions
Agents must actively identify and reuse existing, established abstractions in the codebase to prevent drift and duplication. The project already encapsulates repetitive operations in clean, generalized functions. Do not invent new structures when a suitable abstraction exists.

**Example 1: Tandoor HTTP Client Generics**
Instead of manually decoding JSON or handling HTTP errors per feature, reuse the generalized `tandoor.Request` pattern:
```go
res, err := tandoor.Request[GetRecipesResponse](ctx, client, "GET", "/api/recipe/", qb.Values(), nil)
```

**Example 2: Reusable Test Assertions**
To keep black-box tests clean and focused on the AAA (Arrange, Act, Assert) pattern, do not hand-roll error parsing. Use `AssertToolSuccess` to automatically unwrap SDK errors:
```go
res, err := infra.CallTool(ctx, fixture.Client, "get_recipes", get_recipes.GetRecipesParams{
    Query: "Tandoori Chicken",
})
AssertToolSuccess(t, res, err)
```

## 6. Manual Testing with curl
For manual testing and verifying MCP protocol behavior directly, a helper script is provided at `scripts/curl_mcp.sh`.
- This script connects to the MCP server running locally (e.g. on `http://localhost:8081`) and performs raw JSON-RPC requests over the SSE transport.
- Agents should use/edit this script when debugging or validating tool inputs/outputs instead of hand-rolling custom Python scripts or complex testing setups.

