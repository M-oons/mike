package updater

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/m-oons/mike/internal/info"
	"github.com/m-oons/mike/internal/windows/shell32"
	"github.com/m-oons/mike/internal/windows/user32"
)

type releaseInfo struct {
	Tag string `json:"tag_name"`
	URL string `json:"html_url"`
}

func CheckForUpdates() {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", info.Repository)
	res, err := http.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return
	}

	var release releaseInfo
	if err := json.NewDecoder(res.Body).Decode(&release); err != nil {
		return
	}

	currentVersion := strings.TrimPrefix(info.Version, "v")
	latestVersion := strings.TrimPrefix(release.Tag, "v")
	if latestVersion > currentVersion && info.Version != "dev" {
		showNotification(release.URL, release.Tag)
	}
}

func showNotification(releaseURL string, releaseTag string) {
	title := "Update Available"
	message := fmt.Sprintf("A new version of %s is available. Do you want to go to the download page?\n\nCurrent: [%s]\nNew: [%s]", info.AppName, info.Version, releaseTag)
	if user32.ShowMessageBox(title, message) {
		shell32.OpenURL(releaseURL)
	}
}
