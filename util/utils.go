package util

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/adrg/xdg"
	"github.com/tigerinus/libvirt-go-demo/config"
)

func GetUserPkgData(fileName string) string {
	dir := filepath.Join(xdg.DataHome, config.PackageTarname)

	EnsureDirectory(dir)

	return filepath.Join(dir, fileName)
}

func GetUserPkgCache(fileName string) string {
	dir := filepath.Join(xdg.CacheHome, config.PackageTarname)

	EnsureDirectory(dir)

	return filepath.Join(dir, fileName)
}

func EnsureDirectory(dir string) {
	if err := os.MkdirAll(dir, 0x744); err != nil {
		fmt.Printf("warning: %s", err.Error())
	}
}

func ReplaceRegex(str, old, replacement string) string {
	regex, err := regexp.Compile(old)
	if err != nil {
		fmt.Println(err.Error())
		return str
	}
	return regex.ReplaceAllString(str, replacement)
}
