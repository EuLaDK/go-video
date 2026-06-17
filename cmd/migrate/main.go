package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"sort"

	"next-video-golang/internal/config"
	"next-video-golang/internal/database"

	"github.com/jackc/pgx/v5/pgconn"
)

// main 执行本地数据库初始化；读取 migrations 和 seeds 目录下的 SQL 文件。
func main() {
	ctx := context.Background()
	cfg := config.Load()

	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer pool.Close()

	if err := executeSQLDirectory(ctx, poolExec{exec: pool.Exec}, "migrations"); err != nil {
		log.Fatalf("execute migrations: %v", err)
	}

	if err := executeSQLDirectory(ctx, poolExec{exec: pool.Exec}, "seeds"); err != nil {
		log.Fatalf("execute seeds: %v", err)
	}

	log.Println("database migrated and seeded")
}

type sqlExecutor interface {
	Exec(ctx context.Context, sql string) error
}

type poolExec struct {
	exec func(context.Context, string, ...any) (pgconn.CommandTag, error)
}

// Exec 执行 SQL 文本；ctx 为执行上下文，sql 为待执行语句。
func (executor poolExec) Exec(ctx context.Context, sql string) error {
	_, err := executor.exec(ctx, sql)
	return err
}

// executeSQLDirectory 按文件名顺序执行 SQL 目录；ctx 为执行上下文，executor 为 SQL 执行器，directory 为目录名。
func executeSQLDirectory(ctx context.Context, executor sqlExecutor, directory string) error {
	files, err := filepath.Glob(filepath.Join(directory, "*.sql"))
	if err != nil {
		return err
	}
	sort.Strings(files)

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		log.Printf("executing %s", file)
		if err := executor.Exec(ctx, string(content)); err != nil {
			return err
		}
	}

	return nil
}
