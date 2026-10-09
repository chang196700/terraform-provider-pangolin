// release builds the same six Registry archives locally without publishing.
// Signing is a separate step with the owner's existing GPG key.
package main

import (
	"archive/zip"
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

func main() {
	version := flag.String("version", "0.1.0", "Semantic provider version without v")
	flag.Parse()
	if !regexp.MustCompile(`^\d+\.\d+\.\d+(-[a-zA-Z0-9.-]+)?$`).MatchString(*version) {
		log.Fatal("invalid version")
	}
	if err := build(*version); err != nil {
		log.Fatal(err)
	}
}
func build(version string) error {
	root := "dist"
	if err := os.MkdirAll(root, 0755); err != nil {
		return err
	}
	prefix := "terraform-provider-pangolin_" + version
	checks := map[string]string{}
	manifest, err := os.ReadFile("terraform-registry-manifest.json")
	if err != nil {
		return err
	}
	mn := prefix + "_manifest.json"
	if err = os.WriteFile(filepath.Join(root, mn), manifest, 0644); err != nil {
		return err
	}
	checks[mn] = fmt.Sprintf("%x", sha256.Sum256(manifest))
	for _, target := range []string{"linux/amd64", "linux/arm64", "darwin/amd64", "darwin/arm64", "windows/amd64", "windows/arm64"} {
		parts := strings.Split(target, "/")
		dir := filepath.Join(root, parts[0]+"_"+parts[1])
		if err = os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		name := "terraform-provider-pangolin_v" + version
		if parts[0] == "windows" {
			name += ".exe"
		}
		binary := filepath.Join(dir, name)
		cmd := exec.Command("go", "build", "-trimpath", "-ldflags", "-s -w -X main.version="+version, "-o", binary, ".")
		cmd.Env = append(os.Environ(), "CGO_ENABLED=0", "GOOS="+parts[0], "GOARCH="+parts[1])
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err = cmd.Run(); err != nil {
			return err
		}
		archiveName := prefix + "_" + parts[0] + "_" + parts[1] + ".zip"
		if err = archive(filepath.Join(root, archiveName), map[string]string{name: binary, "LICENSE": "LICENSE"}); err != nil {
			return err
		}
		raw, err := os.ReadFile(filepath.Join(root, archiveName))
		if err != nil {
			return err
		}
		checks[archiveName] = fmt.Sprintf("%x", sha256.Sum256(raw))
		fmt.Println("Built", archiveName)
	}
	names := make([]string, 0, len(checks))
	for k := range checks {
		names = append(names, k)
	}
	sort.Strings(names)
	var b strings.Builder
	for _, name := range names {
		fmt.Fprintf(&b, "%s  %s\n", checks[name], name)
	}
	return os.WriteFile(filepath.Join(root, prefix+"_SHA256SUMS"), []byte(b.String()), 0644)
}
func archive(path string, files map[string]string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	zw := zip.NewWriter(f)
	names := []string{}
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		h := &zip.FileHeader{Name: name, Method: zip.Deflate}
		h.SetModTime(time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC))
		h.SetMode(0755)
		w, err := zw.CreateHeader(h)
		if err != nil {
			return err
		}
		src, err := os.Open(files[name])
		if err != nil {
			return err
		}
		_, err = io.Copy(w, src)
		src.Close()
		if err != nil {
			return err
		}
	}
	if err = zw.Close(); err != nil {
		return err
	}
	return f.Close()
}
