package metabase

import (
	"fmt"
	"regexp"
	"strings"
)

// blockedOperations lists SQL operations that are not allowed through the MCP server.
var blockedOperations = []string{
	"INSERT",
	"UPDATE",
	"DELETE",
	"DROP",
	"ALTER",
	"CREATE",
	"TRUNCATE",
	"GRANT",
	"REVOKE",
	"EXEC",
	"EXECUTE",
	"MERGE",
	"CALL",
}

// singleLineComment matches -- style comments.
var singleLineComment = regexp.MustCompile(`--[^\n]*`)

// multiLineComment matches /* ... */ style comments.
var multiLineComment = regexp.MustCompile(`/\*[\s\S]*?\*/`)

// buildBlockPattern builds a regex for a given SQL keyword.
// Matches the keyword at word boundary, case-insensitive.
func buildBlockPattern(op string) *regexp.Regexp {
	return regexp.MustCompile(`(?i)\b` + op + `\b`)
}

var blockedPatterns []*regexp.Regexp

func init() {
	blockedPatterns = make([]*regexp.Regexp, len(blockedOperations))
	for i, op := range blockedOperations {
		blockedPatterns[i] = buildBlockPattern(op)
	}
}

// ValidateReadOnlySQL checks that a SQL query does not contain write operations.
// It returns an error if a write operation is detected.
func ValidateReadOnlySQL(sql string) error {
	cleaned := stripComments(sql)
	for i, pattern := range blockedPatterns {
		if pattern.MatchString(cleaned) {
			return fmt.Errorf("query contains blocked operation: %s. Only read-only (SELECT) queries are allowed", blockedOperations[i])
		}
	}
	return nil
}

// stripComments removes SQL comments from the query.
func stripComments(sql string) string {
	result := multiLineComment.ReplaceAllString(sql, " ")
	result = singleLineComment.ReplaceAllString(result, " ")
	return strings.TrimSpace(result)
}
