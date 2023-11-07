package ServerFrameworkCore

import (
	"ServerSTARS-II/ServerFrameworkCore/ServerLayerLow/ServerLog"
	"ServerSTARS-II/ServerFrameworkCore/ServerLayerLow/ServerPool"
	"os"
	"path/filepath"
	"strings"
)

var staticFiles []string

// WebStaticLoadResources 加载网页静态资源.
// Static => Cache.
// @param Folder string
// @return void
func WebStaticLoadResources(Folder string) {
	ServerLog.StarLogFmt("INFO", "Static Load Folder: "+Folder, "Server", "Core", "StaticWeb")
	// 遍历 Static Web Folder.
	err := filepath.WalkDir(Folder, visitFolder)
	if err != nil {
		ServerLog.StarLogFmt("WARNING", "Static Assets: "+err.Error(), "Server", "Core", "StaticWeb")
		return
	}
	for _, File := range staticFiles {
		// 删除父目录.
		KeyStringTemp := "/" + File[len(Folder):]
		ServerPool.CachePushFile(KeyStringTemp, File)
	}
}

// visitFolder 遍历文件夹所有文件.
func visitFolder(fp string, fi os.DirEntry, err error) error {
	if err != nil {
		ServerLog.StarLogFmt("WARNING", "Static Assets: "+err.Error(), "Server", "Core", "StaticWeb")
		return nil
	}
	if !fi.IsDir() {
		fp = strings.Replace(fp, "\\", "/", -1)
		staticFiles = append(staticFiles, fp) // push path.
	}
	return nil
}
