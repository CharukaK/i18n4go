package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pivotal-cf-experimental/jibber_jabber"

	goi18n "github.com/CharukaK/i18n4go/i18n4go/i18n"
)

type Detector interface {
	DetectIETF() (string, error)
	DetectLanguage() (string, error)
}

type JibberJabberDetector struct{}

const (
	DEFAULT_LOCALE   = "en_US"
	DEFAULT_LANGUAGE = "en"
)

var T goi18n.TranslateFunc

var SUPPORTED_LOCALES = map[string]string{
	"de": "de_DE",
	"en": "en_US",
	"es": "es_ES",
	"fr": "fr_FR",
	"it": "it_IT",
	"ja": "ja_JA",
	//"ko": "ko_KO", - Will add support for Korean when nicksnyder/go-i18n supports Korean
	"pt": "pt_BR",
	//"ru": "ru_RU", - Will add support for Russian when nicksnyder/go-i18n supports Russian
	"zh": "zh_Hans",
}

var Resources_path = filepath.Join("i18n", "resources")

func init() {
	T = Init(&JibberJabberDetector{})
}

func GetResourcesPath() string {
	return Resources_path
}

func Init(detector Detector) goi18n.TranslateFunc {
	var T goi18n.TranslateFunc
	var err error

	var userLocale string
	userLocale, err = initWithUserLocale(detector)
	if err != nil {
		userLocale = mustLoadDefaultLocale()
	}

	T, err = goi18n.Tfunc(userLocale, DEFAULT_LOCALE)

	if err != nil {
		panic(err)
	}

	return T
}

func initWithUserLocale(detector Detector) (string, error) {
	userLocale, err := detector.DetectIETF()
	if err != nil {
		userLocale = DEFAULT_LOCALE
	}

	language, err := detector.DetectLanguage()
	if err != nil {
		language = DEFAULT_LANGUAGE
	}

	userLocale = strings.Replace(userLocale, "-", "_", 1)
	if strings.HasPrefix(userLocale, "zh_TW") || strings.HasPrefix(userLocale, "zh_HK") {
		userLocale = "zh_Hant"
		language = "zh"
	}

	err = loadFromAsset(userLocale)
	if err != nil {
		locale := SUPPORTED_LOCALES[language]
		if locale == "" {
			userLocale = DEFAULT_LOCALE
		} else {
			userLocale = locale
		}
		err = loadFromAsset(userLocale)
	}

	return userLocale, err
}

func mustLoadDefaultLocale() string {
	userLocale := DEFAULT_LOCALE

	err := loadFromAsset(DEFAULT_LOCALE)
	if err != nil {
		panic("Could not load en_US language files. God save the queen. \n" + err.Error() + "\n\n")
	}

	return userLocale
}

func loadFromAsset(locale string) error {
	assetName := locale + ".all.json"
	assetKey := filepath.Join(GetResourcesPath(), assetName)

	byteArray, err := Asset(assetKey)
	if err != nil {
		return err
	}

	if len(byteArray) == 0 {
		return errors.New(fmt.Sprintf("Could not load i18n asset: %v", assetKey))
	}

	_, err = os.Stat(os.TempDir())
	if err != nil {
		if !os.IsExist(err) {
			return errors.New("Please make sure Temp dir exist - " + os.TempDir())
		} else {
			return err
		}
	}

	tmpDir, err := ioutil.TempDir("", "cloudfoundry_cli_i18n_res")
	if err != nil {
		return err
	}
	defer func() {
		os.RemoveAll(tmpDir)
	}()

	fileName, err := saveLanguageFileToDisk(tmpDir, assetName, byteArray)
	if err != nil {
		return err
	}

	goi18n.MustLoadTranslationFile(fileName)

	os.RemoveAll(fileName)

	return nil
}

func saveLanguageFileToDisk(tmpDir, assetName string, byteArray []byte) (fileName string, err error) {
	fileName = filepath.Join(tmpDir, assetName)
	file, err := os.Create(fileName)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = file.Write(byteArray)
	if err != nil {
		return
	}

	return
}

func (detector *JibberJabberDetector) DetectIETF() (string, error) {
	return jibber_jabber.DetectIETF()
}

func (detector *JibberJabberDetector) DetectLanguage() (string, error) {
	return jibber_jabber.DetectLanguage()
}
