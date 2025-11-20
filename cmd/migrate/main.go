package main

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"test-go/internal/dbConn"
	"time"
)

type Migration struct {
	Name string
	Path string
}

func ensureMigrationsTable(db *sql.DB) error {
	const query = `
CREATE TABLE IF NOT EXISTS schema_migrations (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;
`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}
	return nil
}

func getAppliedMigrations(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query(`SELECT name FROM schema_migrations`)
	if err != nil {
		return nil, fmt.Errorf("select schema_migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("scan migration name: %w", err)
		}
		applied[name] = true
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return applied, nil
}

func listMigrations(dir string) ([]Migration, error) {
	var migrations []Migration

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(d.Name()) == ".sql" {
			migrations = append(migrations, Migration{
				Name: d.Name(),
				Path: path,
			})
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk migrations dir: %w", err)
	}

	// Сортируем по имени файла: 001_..., 002_..., 010_...
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Name < migrations[j].Name
	})

	return migrations, nil
}

func applyMigration(db *sql.DB, m Migration) error {
	log.Printf("applying migration: %s", m.Name)

	sqlBytes, err := os.ReadFile(m.Path)
	if err != nil {
		return fmt.Errorf("read migration file %s: %w", m.Path, err)
	}
	sqlText := string(sqlBytes)

	// Запускаем в транзакции
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx for migration %s: %w", m.Name, err)
	}

	if _, err := tx.Exec(sqlText); err != nil {
		tx.Rollback()
		return fmt.Errorf("exec migration %s: %w", m.Name, err)
	}

	if _, err := tx.Exec(`INSERT INTO schema_migrations (name, applied_at) VALUES (?, ?)`, m.Name, time.Now()); err != nil {
		tx.Rollback()
		return fmt.Errorf("insert into schema_migrations for %s: %w", m.Name, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migration %s: %w", m.Name, err)
	}

	log.Printf("migration applied: %s", m.Name)
	return nil
}

func main() {
	db, err := dbConn.NewDB()
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}
	defer db.Close()

	if err := ensureMigrationsTable(db); err != nil {
		log.Fatalf("failed to ensure schema_migrations: %v", err)
	}

	applied, err := getAppliedMigrations(db)
	if err != nil {
		log.Fatalf("failed to load applied migrations: %v", err)
	}

	migrations, err := listMigrations("migrations")
	if err != nil {
		log.Fatalf("failed to list migrations: %v", err)
	}

	if len(migrations) == 0 {
		log.Println("no migrations found")
		return
	}

	for _, m := range migrations {
		if applied[m.Name] {
			log.Printf("skip already applied migration: %s", m.Name)
			continue
		}
		if err := applyMigration(db, m); err != nil {
			log.Fatalf("failed to apply migration %s: %v", m.Name, err)
		}
	}

	log.Println("all migrations applied successfully")
}
