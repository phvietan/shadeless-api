package schema

import (
	"path/filepath"
	"shadeless-api/main/libs"
	"strings"
)

type StaticScorer struct {
	packet *ParsedPacket
}

func NewStaticScorer(packet *ParsedPacket) *StaticScorer {
	scorer := new(StaticScorer)
	scorer.packet = packet
	return scorer
}

func (this *StaticScorer) GetScore() float64 {
	p := this.packet

	if p.Path == "/robots.txt" && strings.Contains(p.ResponseContentType, "text/plain") {
		return 100
	}
	score := 0.0
	if p.ResponseStatus >= 400 {
		score -= 12 // Why should we fuzz params when we got error? but sometimes we also want to fuzz this, lol
	}
	// Even if it is html, we should fuzz
	if strings.Contains(p.ResponseContentType, "text/html") {
		score -= 30
	}

	// Extension that looks like static
	extension := filepath.Ext(p.Path)
	if extension != "" {
		extension = extension[1:]
	}
	if libs.IsSliceIncludesString([]string{"jpeg", "jpg", "png", "js", "css", "txt", "json", "gif", "xml", "woff2", "woff", "ttf", "docx", "m3u8", "ts", "webp", "ico", "svg"}, extension) {
		score += 40
	}
	// .gif is dirty way for many vendor that collect user's data, so just worth 10 points
	if extension == "gif" {
		score += 30
	}
	// If no param, should we fuzz?
	if len(p.Parameters) == 0 {
		score += 30
	} else {
		// Sometimes website includes timestamp for removing cache
		if len(p.Parameters) == 1 && p.Querystring != "" {
			score += 30
		} else {
			if len(p.Parameters) > 4 { // Has many params? maybe it is api
				score -= 20
			}
		}
	}

	// CDN Origins
	numOrigin := libs.IsStringContains(p.Origin, []string{"asset", "img", "content", "font", "stc", "cdn", "static", "image", "file", "script", "video"})
	score += 10.0 * float64(numOrigin)
	// Ridiculous long origin, wtf? maybe cdn?
	if len(p.Origin) >= 55 {
		score += 20
	}

	// Path that looks like static
	numPath := libs.IsStringContains(p.Path, []string{"static", "cdn", "chunk-", "icon", "asset", "image", "img", "file", "script", "js", "css", "public", "font", "video", "media", "bootstrap", "jquery"})
	score += 10.0 * float64(numPath)

	// Path that looks like api
	numPathApi := libs.IsStringContains(p.Path, []string{"api", "histo", "renew", "token", "count", "get", "noti", "v1", "v2", "v3", "v4", "v5", "account", "auth", "profile", "login", "logout", "register", "sign", "session", "search", "upload"})
	score -= 6.0 * float64(numPathApi)

	// content-type that look like api
	if extension == "" && strings.Contains(p.ResponseContentType, "json") {
		score -= 20
	}
	score = libs.Min(score, 100)
	score = libs.Max(score, 0)
	return score
}
