package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type MCP struct {
}

func NewMCPServer(ctx context.Context) {

	mcpServer := server.NewMCPServer(
		"bono-tools",
		"1.0.0",
	)

	runQueryTool := mcp.NewTool("run_query",
		mcp.WithDescription("Execute a read-only SELECT query against the ClickHouse analytics database. Only SELECT/WITH/EXPLAIN statements are allowed. Results are returned as JSON."),
		mcp.WithString("sql",
			mcp.Required(),
			mcp.Description("The ClickHouse SQL query to execute. Must be a single SELECT statement."),
		),
	)

}
