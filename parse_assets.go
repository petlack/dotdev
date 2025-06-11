package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// GetIncludedAssets parses the HTML file and returns any JS or CSS files referenced via
// <script src="..."></script> or <link rel="stylesheet" href="...">. Returned paths
// are absolute, resolved relative to the HTML file's directory.
func GetIncludedAssets(htmlFile string) ([]string, error) {
	data, err := os.ReadFile(htmlFile)
	if err != nil {
		return nil, err
	}
	content := string(data)

	reScript := regexp.MustCompile(`(?i)<script[^>]*\bsrc=["']([^"']+)["']`)
	reLink := regexp.MustCompile(`(?i)<link[^>]*\brel=["']?stylesheet["']?[^>]*\bhref=["']([^"']+)["']`)

	var assets []string
	for _, m := range reScript.FindAllStringSubmatch(content, -1) {
		assets = append(assets, strings.TrimSpace(m[1]))
	}
	for _, m := range reLink.FindAllStringSubmatch(content, -1) {
		assets = append(assets, strings.TrimSpace(m[1]))
	}

	base := filepath.Dir(htmlFile)
	for i, p := range assets {
		if !filepath.IsAbs(p) {
			assets[i] = filepath.Join(base, p)
		}
	}
	return assets, nil
}
