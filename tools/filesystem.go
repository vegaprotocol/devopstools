package tools

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}

	return err == nil
}

func CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, fmt.Errorf("failed to stat source file: %w", err)
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, fmt.Errorf("failed to read source file: %w", err)
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	if err != nil {
		return 0, fmt.Errorf("failed to copy content of the source file to destination: %w", err)
	}
	return nBytes, nil
}

// CopyDir copies the content of src to dst. src should be a full path.
func CopyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error from previous walk: %w", err)
		}

		// copy to this path
		outpath := filepath.Join(dst, strings.TrimPrefix(path, src))

		if info.IsDir() {
			os.MkdirAll(outpath, info.Mode())
			return nil // means recursive
		}

		// handle irregular files
		if !info.Mode().IsRegular() {
			switch info.Mode().Type() & os.ModeType {
			case os.ModeSymlink:
				link, err := os.Readlink(path)
				if err != nil {
					return fmt.Errorf("failed to read link: %w", err)
				}
				return os.Symlink(link, outpath)
			}
			return nil
		}

		// copy contents of regular file efficiently

		// open input
		in, _ := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open the file(%s): %w", path, err)
		}
		defer in.Close()

		// create output
		fh, err := os.Create(outpath)
		if err != nil {
			return fmt.Errorf("failed to create outpath(%s): %w", outpath, err)
		}
		defer fh.Close()

		// make it the same
		err = fh.Chmod(info.Mode())
		if err != nil {
			return fmt.Errorf("failed to chmod file: %w", err)
		}

		// copy content
		_, err = io.Copy(fh, in)
		if err != nil {
			return fmt.Errorf("failed to copy file content: %w", err)
		}
		return nil
	})
}
