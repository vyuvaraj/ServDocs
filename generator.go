package main

import (
	"html/template"
	"os"
	"strings"
)

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} Documentation</title>
    <link href="https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;600;800&family=JetBrains+Mono:wght@400;700&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg-color: #0d1117;
            --sidebar-bg: #161b22;
            --text-color: #c9d1d9;
            --text-heading: #f0f6fc;
            --accent-color: #58a6ff;
            --accent-gradient: linear-gradient(135deg, #58a6ff, #bc8cff);
            --border-color: #30363d;
            --card-bg: rgba(22, 27, 34, 0.6);
            --badge-get: #2ea043;
            --badge-post: #d29922;
            --badge-put: #1f6feb;
            --badge-delete: #f85149;
        }

        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }

        body {
            font-family: 'Outfit', sans-serif;
            background-color: var(--bg-color);
            color: var(--text-color);
            display: flex;
            height: 100vh;
            overflow: hidden;
        }

        /* Sidebar Styling */
        aside {
            width: 320px;
            background-color: var(--sidebar-bg);
            border-right: 1px solid var(--border-color);
            display: flex;
            flex-direction: column;
            padding: 24px;
            overflow-y: auto;
        }

        .logo {
            font-size: 24px;
            font-weight: 800;
            background: var(--accent-gradient);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            margin-bottom: 24px;
        }

        .search-box {
            background-color: var(--bg-color);
            border: 1px solid var(--border-color);
            border-radius: 8px;
            padding: 10px 14px;
            color: var(--text-color);
            font-family: inherit;
            font-size: 14px;
            margin-bottom: 24px;
            outline: none;
            transition: border-color 0.2s;
        }

        .search-box:focus {
            border-color: var(--accent-color);
        }

        .nav-section {
            margin-bottom: 24px;
        }

        .nav-title {
            font-size: 12px;
            text-transform: uppercase;
            letter-spacing: 1.5px;
            color: #8b949e;
            margin-bottom: 12px;
            font-weight: 600;
        }

        .nav-list {
            list-style: none;
        }

        .nav-item {
            padding: 8px 12px;
            border-radius: 6px;
            font-size: 14px;
            cursor: pointer;
            transition: all 0.2s;
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
        }

        .nav-item:hover {
            background-color: var(--border-color);
            color: var(--text-heading);
        }

        /* Main Panel Styling */
        main {
            flex: 1;
            padding: 40px;
            overflow-y: auto;
            background: radial-gradient(circle at 80% 20%, rgba(88, 166, 255, 0.05), transparent 50%);
        }

        h1 {
            font-size: 36px;
            font-weight: 800;
            color: var(--text-heading);
            margin-bottom: 12px;
        }

        .subtitle {
            color: #8b949e;
            margin-bottom: 40px;
        }

        .section-header {
            font-size: 24px;
            font-weight: 600;
            color: var(--text-heading);
            margin-bottom: 20px;
            border-bottom: 1px solid var(--border-color);
            padding-bottom: 8px;
        }

        .card {
            background-color: var(--card-bg);
            border: 1px solid var(--border-color);
            border-radius: 12px;
            padding: 24px;
            margin-bottom: 24px;
            backdrop-filter: blur(10px);
            transition: transform 0.2s, box-shadow 0.2s;
        }

        .card:hover {
            transform: translateY(-2px);
            box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2);
            border-color: #444c56;
        }

        .card-header {
            display: flex;
            align-items: center;
            gap: 12px;
            margin-bottom: 12px;
        }

        .card-title {
            font-size: 18px;
            font-weight: 600;
            color: var(--text-heading);
            font-family: 'JetBrains Mono', monospace;
        }

        .badge {
            font-size: 12px;
            font-weight: bold;
            text-transform: uppercase;
            padding: 4px 8px;
            border-radius: 4px;
            color: #fff;
            font-family: 'JetBrains Mono', monospace;
        }

        .badge-get { background-color: var(--badge-get); }
        .badge-post { background-color: var(--badge-post); }
        .badge-put { background-color: var(--badge-put); }
        .badge-delete { background-color: var(--badge-delete); }

        .card-desc {
            color: #8b949e;
            font-size: 14px;
            line-height: 1.6;
            margin-bottom: 16px;
        }

        .param-table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 12px;
            font-size: 14px;
        }

        .param-table th, .param-table td {
            text-align: left;
            padding: 10px 12px;
            border-bottom: 1px solid var(--border-color);
        }

        .param-table th {
            color: var(--text-heading);
            font-weight: 600;
        }

        .code-style {
            font-family: 'JetBrains Mono', monospace;
            background-color: var(--bg-color);
            padding: 2px 6px;
            border-radius: 4px;
            color: #ff7b72;
            font-size: 13px;
        }
    </style>
</head>
<body>

    <aside>
        <div class="logo">ServDocs</div>
        <input type="text" class="search-box" id="search" placeholder="Quick search..." onkeyup="filterDocs()">

        <div class="nav-section">
            <div class="nav-title">Routes</div>
            <ul class="nav-list">
                {{range .Doc.Routes}}
                <li class="nav-item" onclick="scrollToElement('route-{{.Method}}-{{.Path}}')">
                    {{.Method}} {{.Path}}
                </li>
                {{end}}
            </ul>
        </div>

        <div class="nav-section">
            <div class="nav-title">Structs</div>
            <ul class="nav-list">
                {{range .Doc.Structs}}
                <li class="nav-item" onclick="scrollToElement('struct-{{.Name}}')">
                    {{.Name}}
                </li>
                {{end}}
            </ul>
        </div>

        <div class="nav-section">
            <div class="nav-title">Functions</div>
            <ul class="nav-list">
                {{range .Doc.Functions}}
                <li class="nav-item" onclick="scrollToElement('fn-{{.Name}}')">
                    {{.Name}}()
                </li>
                {{end}}
            </ul>
        </div>
    </aside>

    <main>
        <h1>{{.Title}} API Reference</h1>
        <div class="subtitle">Auto-generated platform technical specifications.</div>

        {{if .Doc.Routes}}
        <div class="section-header" id="routes-section">HTTP Route Handlers</div>
        {{range .Doc.Routes}}
        <div class="card doc-card" id="route-{{.Method}}-{{.Path}}" data-search="{{.Method}} {{.Path}} {{.Description}}">
            <div class="card-header">
                <span class="badge badge-{{html (print (slice (low .Method) 0))}}">{{.Method}}</span>
                <span class="card-title">{{.Path}}</span>
            </div>
            <p class="card-desc">{{.Description}}</p>
            {{if or .InputType .OutputType}}
            <table class="param-table">
                <thead>
                    <tr>
                        <th>Payload</th>
                        <th>Type / Schema Reference</th>
                    </tr>
                </thead>
                <tbody>
                    {{if .InputType}}
                    <tr>
                        <td>Request Payload</td>
                        <td><span class="code-style">{{.InputType}}</span></td>
                    </tr>
                    {{end}}
                    {{if .OutputType}}
                    <tr>
                        <td>Response Payload</td>
                        <td><span class="code-style">{{.OutputType}}</span></td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
            {{end}}
        </div>
        {{end}}
        {{end}}

        {{if .Doc.Structs}}
        <div class="section-header" id="structs-section">Data Structures & Schemas</div>
        {{range .Doc.Structs}}
        <div class="card doc-card" id="struct-{{.Name}}" data-search="{{.Name}} {{.Description}}">
            <div class="card-header">
                <span class="card-title" style="color: var(--accent-color)">struct {{.Name}}</span>
            </div>
            <p class="card-desc">{{.Description}}</p>
            <table class="param-table">
                <thead>
                    <tr>
                        <th>Field Name</th>
                        <th>Type</th>
                        <th>Description</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .Fields}}
                    <tr>
                        <td style="font-weight: 600">{{.Name}}</td>
                        <td><span class="code-style" style="color: #bc8cff">{{.Type}}</span></td>
                        <td>{{.Description}}</td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
        </div>
        {{end}}
        {{end}}

        {{if .Doc.Functions}}
        <div class="section-header" id="functions-section">Built-in Functions & Helpers</div>
        {{range .Doc.Functions}}
        <div class="card doc-card" id="fn-{{.Name}}" data-search="{{.Name}} {{.Description}}">
            <div class="card-header">
                <span class="card-title" style="color: #ff7b72">fn {{.Name}}({{.InputParams}}) -> {{.ReturnType}}</span>
            </div>
            <p class="card-desc">{{.Description}}</p>
        </div>
        {{end}}
        {{end}}
    </main>

    <script>
        function scrollToElement(id) {
            const element = document.getElementById(id);
            if (element) {
                element.scrollIntoView({ behavior: 'smooth', block: 'center' });
            }
        }

        function filterDocs() {
            const query = document.getElementById('search').value.toLowerCase();
            const cards = document.querySelectorAll('.doc-card');
            cards.forEach(card => {
                const searchData = card.getAttribute('data-search').toLowerCase();
                if (searchData.includes(query)) {
                    card.style.display = 'block';
                } else {
                    card.style.display = 'none';
                }
            });
        }
    </script>
</body>
</html>`

type PageData struct {
	Title string
	Doc   *SrvDoc
}

func GenerateHtml(doc *SrvDoc, title, outputPath string) error {
	tmpl, err := template.New("docs").Funcs(template.FuncMap{
		"low": func(s string) string {
			return strings.ToLower(s)
		},
		"print": func(args ...interface{}) string {
			var out string
			for _, arg := range args {
				out += string(arg.(string))
			}
			return out
		},
		"slice": func(s string, start int) string {
			return s[start:]
		},
	}).Parse(htmlTemplate)
	if err != nil {
		return err
	}

	// Wait, strings package is needed for low/strings functions. Let's make sure it's correct
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	data := PageData{
		Title: title,
		Doc:   doc,
	}

	return tmpl.Execute(file, data)
}
