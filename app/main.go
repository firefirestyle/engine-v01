package app

import (
	"bytes"
	"net/http"

	arTm "github.com/firefirestyle/engine-v01/article/template"
	usTm "github.com/firefirestyle/engine-v01/user/template"
)

var userTemplate = usTm.NewUserTemplate(userConfig)
var userCommentsTemp *arTm.ArtTemplate = arTm.NewArtTemplate(
	arTm.ArtTemplateConfig{
		GroupName:                  "Main",
		KindBaseName:               "FFArt",
		MemcachedOnlyInBlobPointer: true,
	}, userTemplate.GetUserHundlerObj)

func init() {
	var buffer *bytes.Buffer = bytes.NewBufferString("")
	buffer.WriteString("<html><title>K07ME</title><body>")
	buffer.WriteString("<div>")
	buffer.WriteString("<a href=\"/api/v1/twitter/tokenurl/redirect\">redirect</a>")
	buffer.WriteString("</div>")
	buffer.WriteString("</body></html>")

	userTemplate.InitUserApi()
	userCommentsTemp.InitArtApi()
	http.Handle("/", http.FileServer(http.Dir("web")))

}
