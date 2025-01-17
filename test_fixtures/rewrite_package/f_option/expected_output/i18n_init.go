package input_files

import (
	"path/filepath"

	i18n "github.com/CharukaK/i18n4go/i18n4go/i18n"
)

var T i18n.TranslateFunc

func init() {
	T = i18n.Init(filepath.Join("test_fixtures", "rewrite_package", "f_option", "input_files"), i18n.GetResourcesPath())
}
