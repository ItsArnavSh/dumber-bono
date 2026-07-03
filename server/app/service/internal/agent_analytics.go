package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

const defaultRowLimit = 100

// AnalyticsSQL runs an AI-agent-provided read-only query against ClickHouse
// and returns the results as a JSON string. It enforces read-only validation
// app-side; the actual security boundary is the ClickHouse user's readonly
// setting (see ai_agent user grants).
func (r *Repository) AnalyticsSQL(ctx context.Context, query string) (string, error) {
	if err := ValidateReadOnlyQuery(query); err != nil {
		return "", fmt.Errorf("query rejected: %w", err)
	}

	query = ensureLimit(query, defaultRowLimit)

	rows, err := r.chdb.Query(ctx, query)
	if err != nil {
		return "", fmt.Errorf("execute query: %w", err)
	}
	defer rows.Close()

	columnTypes := rows.ColumnTypes()
	columnNames := rows.Columns()

	results := make([]map[string]any, 0)

	for rows.Next() {
		scanTargets := make([]any, len(columnTypes))
		for i, ct := range columnTypes {
			scanTargets[i] = reflect.New(ct.ScanType()).Interface()
		}

		if err := rows.Scan(scanTargets...); err != nil {
			return "", fmt.Errorf("scan row: %w", err)
		}

		row := make(map[string]any, len(columnNames))
		for i, name := range columnNames {
			val := reflect.ValueOf(scanTargets[i]).Elem().Interface()
			row[name] = val
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("row iteration: %w", err)
	}

	out, err := json.Marshal(results)
	if err != nil {
		return "", fmt.Errorf("marshal results: %w", err)
	}

	return string(out), nil
}

// ensureLimit appends a LIMIT clause if the query doesn't already specify one.
func ensureLimit(query string, limit int) string {
	trimmed := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(query), ";"))
	if strings.Contains(strings.ToUpper(trimmed), "LIMIT") {
		return trimmed
	}
	return fmt.Sprintf("%s LIMIT %d", trimmed, limit)
}

var forbiddenKeywords = []string{
	"INSERT", "UPDATE", "DELETE", "ALTER", "DROP", "TRUNCATE",
	"CREATE", "GRANT", "REVOKE", "RENAME", "ATTACH", "DETACH",
	"OPTIMIZE", "SYSTEM", "KILL", "EXCHANGE",
}

var stringLiteralPattern = regexp.MustCompile(`'(?:[^'\\]|\\.)*'`)
var commentPattern = regexp.MustCompile(`--[^\n]*|/\*[\s\S]*?\*/`)

// ValidateReadOnlyQuery rejects any query that isn't a plain read.
// This is a fast-fail filter, NOT the source of truth — the ClickHouse
// user's readonly=1 setting is what actually enforces this server-side.
func ValidateReadOnlyQuery(sql string) error {
	trimmed := strings.TrimSpace(sql)
	if trimmed == "" {
		return fmt.Errorf("empty query")
	}

	// Strip string literals and comments so keywords inside them
	// (e.g. a column value containing the word "DROP") don't false-positive.
	stripped := commentPattern.ReplaceAllString(trimmed, "")
	stripped = stringLiteralPattern.ReplaceAllString(stripped, "''")

	// Reject multiple statements — one query in, one query out.
	withoutTrailingSemi := strings.TrimSuffix(strings.TrimSpace(stripped), ";")
	if strings.Contains(withoutTrailingSemi, ";") {
		return fmt.Errorf("multiple statements are not allowed")
	}

	upper := strings.ToUpper(stripped)

	// Must start with SELECT, WITH (CTE), or EXPLAIN.
	firstWord := strings.Fields(upper)
	if len(firstWord) == 0 {
		return fmt.Errorf("empty query")
	}
	switch firstWord[0] {
	case "SELECT", "WITH", "EXPLAIN":
		// allowed
	default:
		return fmt.Errorf("only SELECT/WITH/EXPLAIN statements are allowed, got: %s", firstWord[0])
	}

	for _, kw := range forbiddenKeywords {
		if matchesWholeWord(upper, kw) {
			return fmt.Errorf("query contains forbidden keyword: %s", kw)
		}
	}

	return nil
}

func matchesWholeWord(text, word string) bool {
	pattern := `\b` + regexp.QuoteMeta(word) + `\b`
	matched, _ := regexp.MatchString(pattern, text)
	return matched
}
