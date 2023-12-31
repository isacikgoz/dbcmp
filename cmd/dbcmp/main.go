package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/isacikgoz/dbcmp/internal/store"
	"github.com/spf13/cobra"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "dbcmp",
		Short: "dbcmp is your go-to database content comparison tool",
		Long:  "dbcmp is a powerful and efficient command line tool designed to simplify the process of comparing content between two databases.",
		RunE:  runRootCmdFn,
	}

	rootCmd.PersistentFlags().String("source", "", "source database dsn")
	rootCmd.PersistentFlags().String("target", "", "target database dsn")
	rootCmd.Flags().StringSlice("exclude", []string{}, "exclude tables from comparison, takes comma-separated values.")
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runRootCmdFn(cmd *cobra.Command, args []string) error {
	source, err := cmd.Flags().GetString("source")
	if err != nil {
		return err
	}

	target, err := cmd.Flags().GetString("target")
	if err != nil {
		return err
	}

	excl, err := cmd.Flags().GetStringSlice("exclude")
	if err != nil {
		return err
	}

	diffs, err := store.Compare(source, target, store.CompareOptions{
		ExcludePatterns: excl,
	})
	if err != nil {
		return fmt.Errorf("error during comparison: %w", err)
	}

	if len(diffs) == 0 {
		fmt.Println("Database values are same.")
		return nil
	}

	fmt.Printf("Database values differ. Tables: %s\n", strings.Join(diffs, ", "))

	return nil
}
