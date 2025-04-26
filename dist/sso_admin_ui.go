package dist

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed sso_admin_ui/*
var sso_content embed.FS

func SSOAdminUI() {
	// Create the dist directory if it does not exist
	err := os.MkdirAll("dist", os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory dist:", err)
		return
	}

	// Walk through the embedded files and copy them to the target directory
	err = fs.WalkDir(sso_content, "sso_admin_ui", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Read the file content
		fileContent, err := sso_content.ReadFile(path)
		if err != nil {
			return err
		}

		// Create the target file path, preserving the directory structure
		targetPath := filepath.Join("dist", path)

		// Create the target directory if it does not exist
		targetDir := filepath.Dir(targetPath)
		err = os.MkdirAll(targetDir, os.ModePerm)
		if err != nil {
			return err
		}

		// Write the file content to the target path
		err = os.WriteFile(targetPath, fileContent, os.ModePerm)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error copying files:", err)
	}
}
