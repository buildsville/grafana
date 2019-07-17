package api

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/grafana/grafana/pkg/api/dtos"
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/components/imguploader"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/rendering"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/grafana/grafana/pkg/util"
)

func (hs *HTTPServer) ShareToSlack(c *m.ReqContext, dto dtos.ShareSlack) {
	queryReader, err := util.NewURLQueryReader(c.Req.URL)
	if err != nil {
		c.Handle(400, "Render parameters error", err)
		return
	}

	uploader, err := imguploader.NewImageUploader()
	if err != nil {
		c.Handle(400, "Cannot create ImageUploader", err)
		return
	}

	ctx := c.Req.Context()

	width := 1000
	height := 500
	timeout := 60

	path := fmt.Sprintf("d-solo/%s/%s%s", dto.Uid, dto.Slug, dto.Param)

	result, err := hs.RenderService.Render(ctx, rendering.Opts{
		Width:           width,
		Height:          height,
		Timeout:         time.Duration(timeout) * time.Second,
		OrgId:           c.OrgId,
		UserId:          c.UserId,
		OrgRole:         c.OrgRole,
		Path:            path,
		Timezone:        queryReader.Get("tz", ""),
		Encoding:        queryReader.Get("encoding", ""),
		ConcurrentLimit: 30,
	})

	if err != nil && err == rendering.ErrTimeout {
		c.Handle(500, err.Error(), err)
		return
	}

	if err != nil && err == rendering.ErrPhantomJSNotInstalled {
		if strings.HasPrefix(runtime.GOARCH, "arm") {
			c.Handle(500, "Rendering failed - PhantomJS isn't included in arm build per default", err)
		} else {
			c.Handle(500, "Rendering failed - PhantomJS isn't installed correctly", err)
		}
		return
	}

	if err != nil {
		c.Handle(500, "Rendering failed.", err)
		return
	}

	imgUrl, err := uploader.Upload(ctx, result.FilePath)
	if err != nil {
		c.Handle(500, "Upload failed.", err)
		return
	}

	body := map[string]interface{}{
		"attachments": []map[string]interface{}{
			{
				"color":       "#36a64f",
				"title":       `'` + dto.PanelName + `' in '` + dto.Slug + `'`,
				"title_link":  dto.RawURL,
				"text":        "Shared by " + c.Login,
				"image_url":   imgUrl,
				"footer":      "Grafana v" + setting.BuildVersion,
				"footer_icon": "https://grafana.com/assets/img/fav32.png",
				"ts":          time.Now().Unix(),
			},
		},
		"parse": "full", // to linkify urls, users and channels in alert message.
	}

	if dto.Channel != "__default__" {
		body["channel"] = dto.Channel
	}
	if setting.ShareSlackUserName != "" {
		body["username"] = setting.ShareSlackUserName
	}
	if setting.ShareSlackIconEmoji != "" {
		body["icon_emoji"] = setting.ShareSlackIconEmoji
	}
	if setting.ShareSlackIconURL != "" {
		body["icon_url"] = setting.ShareSlackIconURL
	}
	data, _ := json.Marshal(&body)
	cmd := &m.SendWebhookSync{Url: setting.ShareSlackWebhook, Body: string(data)}

	if err := bus.DispatchCtx(ctx, cmd); err != nil {
		c.Handle(500, "Failed to send slack notification", err)
		return
	}
}
