package app

import (
	"net/http"

	paArtTmp "github.com/firefirestyle/engine-v01/article/template"
	paUsrTmp "github.com/firefirestyle/engine-v01/user/template"
)

var usrTemplate = paUsrTmp.NewUserTemplate(usrConfig)
var artTemplate *paArtTmp.ArtTemplate = paArtTmp.NewArtTemplate(
	paArtTmp.ArtTemplateConfig{
		KindBaseName: "fa",
	}, usrTemplate.GetUserHundlerObj)

func init() {
	usrTemplate.InitUserApi()
	artTemplate.InitArtApi()
	http.Handle("/", http.FileServer(http.Dir("web")))
}
