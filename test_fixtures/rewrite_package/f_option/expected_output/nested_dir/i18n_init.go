package nested_dir

import (
	"github.com/cloudfoundry/cli/cf/i18n"
	goi18n "github.com/nicksnyder/go-i18n/i18n"
)

var T goi18n.TranslateFunc

func init() {
	var err error
	T, err = i18n.Init("nested_dir", i18n.GetResourcesPath())
	if err != nil {
		panic(err)
	}
}