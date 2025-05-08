package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/github/github-mcp-server/pkg/translations"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// GetGist creates a tool to get the details of a specific gist in GitHub.
func GetGist(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_gist",
			mcp.WithDescription(t("TOOL_GET_GIST_DESCRIPTION", "Get details of a specific gist in GitHub.")),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        t("TOOL_GET_GIST_USER_TITLE", "Get gist details"),
				ReadOnlyHint: true,
			}),
			mcp.WithNumber("gist_id",
				mcp.Required(),
				mcp.Description("The id of the gist to retrieve"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			gistID, err := requiredParam[string](request, "gist_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := getClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get GitHub client: %w", err)
			}
			issue, resp, err := client.Gists.Get(ctx, gistID)
			if err != nil {
				return nil, fmt.Errorf("failed to get gist: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to get gist: %s", string(body))), nil
			}

			r, err := json.Marshal(issue)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal gist: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

// ListGists creates a tool to list the gists of the authenticated user.
func ListGists(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("list_gists",
			mcp.WithDescription(t("TOOL_LIST_DESCRIPTION", "List the gists of the authenticated user.")),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        t("TOOL_LIST_GISTS_USER_TITLE", "List gists"),
				ReadOnlyHint: true,
			}),
			WithPagination(),
		),
		func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			client, err := getClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get GitHub client: %w", err)
			}
			gists, resp, err := client.Gists.ListStarred(ctx, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to list gists: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to list gists: %s", string(body))), nil
			}

			r, err := json.Marshal(gists)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal gists: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

// ListStarredGists creates a tool to list the starred gists of the authenticated user.
func ListStarredGists(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("list_starred_gists",
			mcp.WithDescription(t("TOOL_LIST_STARRED_DESCRIPTION", "List the starred gists of the authenticated user.")),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        t("TOOL_LIST_STARRED_GISTS_USER_TITLE", "List starred gists"),
				ReadOnlyHint: true,
			}),
			WithPagination(),
		),
		func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			client, err := getClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get GitHub client: %w", err)
			}
			gists, resp, err := client.Gists.ListStarred(ctx, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to list starred gists: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to list starred gists: %s", string(body))), nil
			}

			r, err := json.Marshal(gists)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal starred gists: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}
