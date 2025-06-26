package toolkit

import (
	"fmt"
	"net/http"
	"path"
)

//DownloadStaticFile downloads a file, and tried to force the browser to avoid
//displaying it in the browser window by setting content-disposition header. It also
//allows specification of the display name

func (t *Tools) DownloadStaticFile(w http.ResponseWriter, r *http.Request, folder, file, displayName string) {
	fp := path.Join(folder, file)

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", displayName))

	http.ServeFile(w, r, fp)
}
