package downloadByContent

import (
	"regexp"
	"strings"

	"github.com/admpub/log"
	"github.com/coscms/webcore/model"
	"github.com/coscms/webcore/model/file/storer"
	uploadPrepare "github.com/coscms/webcore/registry/upload/prepare"
	"github.com/coscms/webfront/library/download"
	"github.com/webx-top/echo"
)

func SyncRemoteImage(
	ctx echo.Context,
	subdir string,
	articleID string,
	content string,
	contentType string,
	disableChunk ...bool,
) (string, error) {
	images := OutsideImage(content, contentType)
	if len(images) == 0 {
		return content, nil
	}
	downloaded := map[string]string{} // 站外图片 => 下载到本站后的图片
	storerInfo := storer.Get()
	prepareData, err := uploadPrepare.Prepare(ctx, subdir, `image`, storerInfo)
	if err != nil {
		return ``, err
	}
	defer prepareData.Close()
	cloudStorage := model.NewCloudStorage(ctx)
	if len(storerInfo.ID) > 0 {
		cloud := storerInfo.Cloud()
		if cloud.Id > 0 {
			cloudStorage.NgingCloudStorage = cloud
		}
	}
	var excludeURLPrefixRegexp *regexp.Regexp
	if cloudStorage.Id > 0 {
		excludeURLPrefixRegexp = regexp.MustCompile(`(?i)^` + regexp.QuoteMeta(cloudStorage.BaseURL()+`/`))
	}
	var _disableChunk bool
	if len(disableChunk) > 0 {
		_disableChunk = disableChunk[0]
	} else {
		_disableChunk = true
	}
	for matched, image := range images {
		if excludeURLPrefixRegexp != nil && excludeURLPrefixRegexp.MatchString(image) {
			delete(images, matched)
			continue
		}
		if _, ok := downloaded[image]; ok {
			continue
		}
		result, _, err := download.Download(
			ctx,
			download.OptionsFileURL(image),
			download.OptionsID(articleID),
			//download.OptionsCheckDir(true),
			download.OptionsPrepareData(prepareData),
			download.OptionsDisableChunk(_disableChunk),
		)
		if err != nil {
			log.Error(image, `: `, err)
			continue
		}
		downloaded[image] = result.FileURL
	}

	for matched, image := range images {
		newAddr, ok := downloaded[image]
		if !ok {
			continue
		}
		newContent := strings.ReplaceAll(matched, image, newAddr)
		content = strings.ReplaceAll(content, matched, newContent)
	}
	return content, nil
}
