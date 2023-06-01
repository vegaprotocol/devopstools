package tools

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
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

func CopyDirectory(scrDir, dest string) error {
	entries, err := os.ReadDir(scrDir)
	if err != nil {
		return fmt.Errorf("failed to read source dir(%s): %w", scrDir, err)
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(scrDir, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			return fmt.Errorf("failed to stat source path(%s): %w", sourcePath, err)
		}

		stat, ok := fileInfo.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("failed to get raw syscall.Stat_t data for '%s'", sourcePath)
		}

		switch fileInfo.Mode() & os.ModeType {
		case os.ModeDir:
			if err := CreateIfNotExists(destPath, 0755); err != nil {
				return fmt.Errorf("failed to create destination path(%s) if not exists: %w", destPath, err)
			}
			if err := CopyDirectory(sourcePath, destPath); err != nil {
				return fmt.Errorf("failed to copy directory from '%s' to '%s': %w", sourcePath, destPath, err)
			}
		case os.ModeSymlink:
			if err := CopySymLink(sourcePath, destPath); err != nil {
				return fmt.Errorf("failed to copy symlink from %s to %s: %w", sourcePath, destPath, err)
			}
		default:
			if err := Copy(sourcePath, destPath); err != nil {
				return fmt.Errorf("failed to copy regular file from '%s' to '%s': %w", sourcePath, destPath, err)
			}
		}

		if err := os.Lchown(destPath, int(stat.Uid), int(stat.Gid)); err != nil {
			return fmt.Errorf("failed to chwon path %s: %w", destPath, err)
		}

		fInfo, err := entry.Info()
		if err != nil {
			return fmt.Errorf("failed to get info about file(%s): %w", entry.Name(), err)
		}

		isSymlink := fInfo.Mode()&os.ModeSymlink != 0
		if !isSymlink {
			if err := os.Chmod(destPath, fInfo.Mode()); err != nil {
				return fmt.Errorf("failed to chown for path(%s): %w", destPath, err)
			}
		}
	}
	return nil
}

func Copy(srcFile, dstFile string) error {
	out, err := os.Create(dstFile)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	defer out.Close()

	in, err := os.Open(srcFile)
	defer in.Close()
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}

	_, err = io.Copy(out, in)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

func Exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func CreateIfNotExists(dir string, perm os.FileMode) error {
	if Exists(dir) {
		return nil
	}

	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}

	return nil
}

func CopySymLink(source, dest string) error {
	link, err := os.Readlink(source)
	if err != nil {
		return fmt.Errorf("failed to read link: %w", err)
	}
	return os.Symlink(link, dest)
}
