package filesystem

import (
	"errors"
	"strings"
	"time"

	"github.com/geisonbiazus/blog/internal/core/blog"
)

var ErrInvalidTime = errors.New("invalid time format, please use YYYY-MM-DD HH:MM")
var ErrInvalidFormat = errors.New("invalid file format, please include a header / body separator \"--\"")

func ParseFileContent(content string) (blog.Post, error) {
	return newParser(content).parse()
}

type parser struct {
	content string
	post    blog.Post
	err     error
}

func newParser(content string) *parser {
	return &parser{content: content}
}

func (p *parser) parse() (blog.Post, error) {
	header, body := p.splitHeaderAndBody()

	p.parseHeader(header)
	p.post.Markdown = body

	return p.post, p.err
}

func (p *parser) splitHeaderAndBody() (header, body string) {
	parts := strings.SplitN(p.content, "--\n", 2)

	if len(parts) != 2 {
		p.err = ErrInvalidFormat
		return
	}

	header = parts[0]
	body = parts[1]
	return
}

func (p *parser) parseHeader(header string) {
	lines := strings.Split(header, "\n")

	for _, line := range lines {
		p.parseTitle(line)
		p.parseAuthor(line)
		p.parseDescription(line)
		p.parseImagePath(line)
		p.parsePostTime(line)
	}
}

func (p *parser) parseTitle(content string) {
	if title := p.parseString(content, "title:"); title != "" {
		p.post.Title = title
	}
}

func (p *parser) parseAuthor(content string) {
	if author := p.parseString(content, "author:"); author != "" {
		p.post.Author = author
	}
}

func (p *parser) parseDescription(content string) {
	if description := p.parseString(content, "description:"); description != "" {
		p.post.Description = description
	}
}

func (p *parser) parseImagePath(content string) {
	if image_path := p.parseString(content, "image_path:"); image_path != "" {
		p.post.ImagePath = image_path
	}
}

func (p *parser) parsePostTime(content string) {
	parsedTime, err := p.parseTime(content, "time:")

	if err != nil {
		p.err = err
	}

	if (parsedTime != time.Time{}) {
		p.post.Time = parsedTime
	}
}

const timeFormat = "2006-01-02 15:04"

func (p *parser) parseTime(content, field string) (time.Time, error) {
	timeString := p.parseString(content, field)

	if timeString == "" {
		return time.Time{}, nil
	}

	t, err := time.Parse(timeFormat, timeString)

	if err != nil {
		return time.Time{}, ErrInvalidTime
	}

	return t, nil

}

func (p *parser) parseString(content, field string) string {
	if strings.HasPrefix(content, field) {
		return strings.Trim(content[len(field):], " \n")
	}
	return ""
}
