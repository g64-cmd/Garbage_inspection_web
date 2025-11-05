//go:build !dev

package main

import "patrol-cloud/internal/db"

// createTempUser 是一个空函数，在生产构建中不执行任何操作。
func createTempUser(repo db.Repository) {
	// 在生产环境中不创建临时用户
}
