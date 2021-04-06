package renderer_test

import (
	"testing"

	"github.com/geisonbiazus/blog/internal/adapters/renderer"
	"github.com/geisonbiazus/blog/pkg/assert"
)

func TestGoldmarkRenderer(t *testing.T) {
	t.Run("Given a markdown string, it converts to HTML", func(t *testing.T) {
		rend := renderer.NewGoldmarkRenderer()

		html, err := rend.Render(sampleMarkdown)
		assert.Equal(t, sampleHTML, html)
		assert.Nil(t, err)
	})
}

const sampleMarkdown = `# Title 1

## Title 2

### Title 3

#### Title 4

This is a paragraph

This is a paragraph. \
But now with a line break in the middle.

**List:**

- item 1
- item 2
- item 3

*Another list*:

1. Item 1
1. Item 2
1. Item 3

[Link](http://example.com)

![Image](http://example.com/image.png)

` + "```" + `
Code Block
` + "```"

const sampleHTML = `<h1>Title 1</h1>
<h2>Title 2</h2>
<h3>Title 3</h3>
<h4>Title 4</h4>
<p>This is a paragraph</p>
<p>This is a paragraph.<br>
But now with a line break in the middle.</p>
<p><strong>List:</strong></p>
<ul>
<li>item 1</li>
<li>item 2</li>
<li>item 3</li>
</ul>
<p><em>Another list</em>:</p>
<ol>
<li>Item 1</li>
<li>Item 2</li>
<li>Item 3</li>
</ol>
<p><a href="http://example.com">Link</a></p>
<p><img src="http://example.com/image.png" alt="Image"></p>
<pre><code>Code Block
</code></pre>
`
