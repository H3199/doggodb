package cli_test

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func runCLI(t *testing.T, input string) string {
	t.Helper()

	// Prepare the command to run the CLI.
	cmd := exec.Command("go", "run", "../../cmd/client/cli.go")
	cmd.Stdin = strings.NewReader(input)

	// Capture the output.
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the CLI and check for errors.
	err := cmd.Run()
	if err != nil {
		t.Fatalf("CLI failed with error: %v\nStderr: %s", err, stderr.String())
	}

	return stdout.String()
}

func TestCLIIntegration(t *testing.T) {
	// Test case: Create a table.
	t.Run("CreateTable", func(t *testing.T) {
		fmt.Println("Running Test: CreateTable")
		input := `
CREATE TABLE users;
exit;
`
		output := runCLI(t, input)
		require.Contains(t, output, "Table 'users' created successfully")
	})

	// Test case: Insert a row.
	t.Run("InsertRow", func(t *testing.T) {
		fmt.Println("Running Test: InsertRow")
		input := `
INSERT INTO users (id, name, email) VALUES (1, 'Alice', 'alice@example.com');
exit;
`
		output := runCLI(t, input)
		require.Contains(t, output, "Row inserted successfully")
	})

	// Test case: Select rows.
	t.Run("SelectRows", func(t *testing.T) {
		fmt.Println("Running Test: SelectRows")
		input := `
SELECT * FROM users;
exit;
`
		output := runCLI(t, input)
		require.Contains(t, output, "Alice")
		require.Contains(t, output, "alice@example.com")
	})

	// Test case: Update rows.
	t.Run("UpdateRows", func(t *testing.T) {
		fmt.Println("Running Test: UpdateRows")
		input := `
UPDATE users SET email = 'alice@newdomain.com' WHERE id = 1;
SELECT * FROM users;
exit;
`
		output := runCLI(t, input)
		require.Contains(t, output, "Row(s) updated: 1")
		require.Contains(t, output, "alice@newdomain.com")
	})

	// Test case: Delete rows.
	t.Run("DeleteRows", func(t *testing.T) {
		fmt.Println("Running Test: DeleteRows")
		input := `
DELETE FROM users WHERE id = 1;
SELECT * FROM users;
exit;
`
		output := runCLI(t, input)
		require.Contains(t, output, "Row(s) deleted: 1")
		require.NotContains(t, output, "Alice")
	})
}
