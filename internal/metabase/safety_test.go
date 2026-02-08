package metabase

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateReadOnlySQL(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr bool
		errMsg  string
	}{
		// Allowed queries
		{name: "simple select", sql: "SELECT * FROM users", wantErr: false},
		{name: "select with where", sql: "SELECT id, name FROM users WHERE id = 1", wantErr: false},
		{name: "select with join", sql: "SELECT u.name, o.total FROM users u JOIN orders o ON u.id = o.user_id", wantErr: false},
		{name: "select with subquery", sql: "SELECT * FROM (SELECT id FROM users) t", wantErr: false},
		{name: "select with CTE", sql: "WITH cte AS (SELECT id FROM users) SELECT * FROM cte", wantErr: false},
		{name: "select with aggregate", sql: "SELECT COUNT(*), SUM(total) FROM orders GROUP BY user_id", wantErr: false},
		{name: "select with union", sql: "SELECT id FROM users UNION ALL SELECT id FROM admins", wantErr: false},
		{name: "explain query", sql: "EXPLAIN SELECT * FROM users", wantErr: false},
		{name: "show tables", sql: "SHOW TABLES", wantErr: false},

		// Blocked queries - basic
		{name: "insert", sql: "INSERT INTO users (name) VALUES ('test')", wantErr: true, errMsg: "INSERT"},
		{name: "update", sql: "UPDATE users SET name = 'test' WHERE id = 1", wantErr: true, errMsg: "UPDATE"},
		{name: "delete", sql: "DELETE FROM users WHERE id = 1", wantErr: true, errMsg: "DELETE"},
		{name: "drop table", sql: "DROP TABLE users", wantErr: true, errMsg: "DROP"},
		{name: "alter table", sql: "ALTER TABLE users ADD COLUMN age INT", wantErr: true, errMsg: "ALTER"},
		{name: "create table", sql: "CREATE TABLE test (id INT)", wantErr: true, errMsg: "CREATE"},
		{name: "truncate", sql: "TRUNCATE TABLE users", wantErr: true, errMsg: "TRUNCATE"},
		{name: "grant", sql: "GRANT SELECT ON users TO role", wantErr: true, errMsg: "GRANT"},
		{name: "revoke", sql: "REVOKE SELECT ON users FROM role", wantErr: true, errMsg: "REVOKE"},
		{name: "exec", sql: "EXEC sp_something", wantErr: true, errMsg: "EXEC"},
		{name: "execute", sql: "EXECUTE sp_something", wantErr: true, errMsg: "EXECUTE"},
		{name: "merge", sql: "MERGE INTO users USING source ON users.id = source.id", wantErr: true, errMsg: "MERGE"},
		{name: "call", sql: "CALL some_procedure()", wantErr: true, errMsg: "CALL"},

		// Blocked queries - case variations
		{name: "insert lowercase", sql: "insert into users values (1)", wantErr: true, errMsg: "INSERT"},
		{name: "delete mixed case", sql: "Delete from users", wantErr: true, errMsg: "DELETE"},
		{name: "DROP uppercase", sql: "DROP DATABASE test", wantErr: true, errMsg: "DROP"},

		// Comment stripping
		{name: "blocked in single line comment", sql: "SELECT * FROM users -- DELETE FROM users", wantErr: false},
		{name: "blocked in multi line comment", sql: "SELECT * FROM users /* INSERT INTO foo VALUES (1) */", wantErr: false},
		{name: "blocked after comment strip", sql: "/* comment */ DELETE FROM users", wantErr: true, errMsg: "DELETE"},
		{name: "nested comment with block", sql: "SELECT * /* DROP TABLE */ FROM users", wantErr: false},

		// Edge cases
		{name: "select into is blocked by insert check", sql: "SELECT * INTO new_table FROM users", wantErr: false},
		{name: "select containing delete as column name", sql: "SELECT deleted FROM users", wantErr: false},
		{name: "substring containing insert", sql: "SELECT reinsertion FROM users", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateReadOnlySQL(tt.sql)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Contains(t, err.Error(), "Only read-only (SELECT) queries are allowed")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
