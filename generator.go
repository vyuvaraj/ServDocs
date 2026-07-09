package main

import (
	"encoding/json"
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

        .payload-cell {
            display: flex;
            align-items: center;
            gap: 12px;
        }

        .toggle-schema-btn {
            background: transparent;
            border: 1px solid var(--border-color);
            border-radius: 4px;
            color: var(--accent-color);
            padding: 2px 8px;
            font-size: 11px;
            font-family: inherit;
            cursor: pointer;
            transition: all 0.2s;
        }

        .toggle-schema-btn:hover {
            background-color: var(--border-color);
            border-color: var(--accent-color);
        }

        .code-examples-container {
            margin-top: 16px;
            border: 1px solid var(--border-color);
            border-radius: 8px;
            background-color: rgba(0, 0, 0, 0.2);
            overflow: hidden;
        }

        .example-tabs {
            display: flex;
            background-color: var(--sidebar-bg);
            border-bottom: 1px solid var(--border-color);
        }

        .tab-btn {
            background: transparent;
            border: none;
            border-right: 1px solid var(--border-color);
            color: #8b949e;
            padding: 8px 16px;
            font-family: inherit;
            font-size: 13px;
            font-weight: 600;
            cursor: pointer;
            outline: none;
            transition: all 0.2s;
        }

        .tab-btn:hover {
            color: var(--text-heading);
            background-color: rgba(255, 255, 255, 0.05);
        }

        .tab-btn.active {
            color: var(--accent-color);
            background-color: var(--bg-color);
            border-bottom: 2px solid var(--accent-color);
            margin-bottom: -1px;
        }

        .example-content {
            padding: 16px;
            font-family: 'JetBrains Mono', monospace;
            font-size: 13px;
            overflow-x: auto;
            max-height: 300px;
        }

        .example-content pre {
            margin: 0;
            white-space: pre-wrap;
            word-break: break-all;
        }

        .nested-schema-container {
            margin-top: 10px;
            padding: 12px;
            background-color: rgba(0, 0, 0, 0.2);
            border: 1px solid var(--border-color);
            border-radius: 6px;
        }

        .nested-schema-table {
            width: 100%;
            border-collapse: collapse;
            font-size: 13px;
        }

        .nested-schema-table th, .nested-schema-table td {
            padding: 6px 8px;
            border-bottom: 1px solid var(--border-color);
            text-align: left;
        }

        .nested-schema-table th {
            color: var(--text-heading);
            font-weight: 600;
        }
    </style>
</head>
<body>

    <aside>
        <div class="logo">ServDocs</div>
        {{if .Versions}}
        <div class="version-select-container" style="margin-bottom: 20px; display: flex; align-items: center; gap: 8px;">
            <label for="version-select" style="font-size: 11px; color: #8b949e; font-weight: 600; text-transform: uppercase; letter-spacing: 0.5px;">Version:</label>
            <select id="version-select" onchange="window.location.href='../' + this.value + '/index.html'" style="background-color: var(--bg-color); border: 1px solid var(--border-color); color: var(--text-color); padding: 4px 8px; border-radius: 6px; font-size: 13px; font-family: inherit; outline: none; cursor: pointer; flex: 1;">
                {{$curr := .CurrentVersion}}
                {{range .Versions}}
                <option value="{{.}}" {{if eq . $curr}}selected{{end}}>{{.}}</option>
                {{end}}
            </select>
        </div>
        {{end}}
        <input type="text" class="search-box" id="search" placeholder="Quick search..." onkeyup="filterDocs()">
        <div id="search-status" style="font-size: 12px; color: #8b949e; margin-top: -16px; margin-bottom: 24px; display: none;"></div>

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
            {{if .Middlewares}}
            <div class="middleware-chain" style="margin-top: -8px; margin-bottom: 12px; display: flex; align-items: center; gap: 6px; flex-wrap: wrap;">
                <span style="font-size: 11px; color: #8b949e; font-weight: 600; text-transform: uppercase; letter-spacing: 0.5px;">Middleware Chain:</span>
                {{range .Middlewares}}
                <span class="middleware-badge" style="color: #bc8cff; font-size: 11px; font-family: 'JetBrains Mono', monospace; border-radius: 4px; padding: 2px 6px; background-color: rgba(188, 140, 255, 0.1); border: 1px solid rgba(188, 140, 255, 0.2);">{{.}}</span>
                {{end}}
            </div>
            {{end}}
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
                        <td>
                            <div class="payload-cell">
                                <span class="code-style">{{.InputType}}</span>
                                <button class="toggle-schema-btn" onclick="toggleSchema(this, '{{.InputType}}')">▸ Expand Schema</button>
                            </div>
                            <div class="nested-schema-container" style="display: none;"></div>
                        </td>
                    </tr>
                    {{end}}
                    {{if .OutputType}}
                    <tr>
                        <td>Response Payload</td>
                        <td>
                            <div class="payload-cell">
                                <span class="code-style">{{.OutputType}}</span>
                                <button class="toggle-schema-btn" onclick="toggleSchema(this, '{{.OutputType}}')">▸ Expand Schema</button>
                            </div>
                            <div class="nested-schema-container" style="display: none;"></div>
                        </td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
            {{end}}
            <div class="code-examples-container" data-method="{{.Method}}" data-path="{{.Path}}" data-input="{{.InputType}}">
                <div class="example-tabs">
                    <button class="tab-btn active" onclick="switchTab(this, 'curl')">cURL</button>
                    <button class="tab-btn" onclick="switchTab(this, 'go')">Go</button>
                    <button class="tab-btn" onclick="switchTab(this, 'js')">JavaScript</button>
                </div>
                <div class="example-content">
                    <pre><code class="code-box"></code></pre>
                </div>
            </div>
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
        const srvSchemas = JSON.parse({{.SchemasJSON}});

        function scrollToElement(id) {
            const element = document.getElementById(id);
            if (element) {
                element.scrollIntoView({ behavior: 'smooth', block: 'center' });
            }
        }

        function toggleSchema(btn, structName) {
            const container = btn.parentElement.nextElementSibling;
            if (container.style.display === 'block') {
                container.style.display = 'none';
                btn.innerText = '▸ Expand Schema';
            } else {
                container.style.display = 'block';
                btn.innerText = '▾ Collapse Schema';
                if (!container.innerHTML) {
                    const schema = srvSchemas[structName];
                    if (!schema || !schema.fields || schema.fields.length === 0) {
                        container.innerHTML = '<span style="color: #ff7b72">No details found for ' + structName + '</span>';
                        return;
                    }
                    let html = '<table class="nested-schema-table"><thead><tr><th>Field Name</th><th>Type</th><th>Description</th></tr></thead><tbody>';
                    schema.fields.forEach(field => {
                        html += '<tr><td style="font-weight: 600">' + field.name + '</td><td><span class="code-style" style="color: #bc8cff">' + field.type + '</span></td><td>' + (field.description || '-') + '</td></tr>';
                    });
                    html += '</tbody></table>';
                    container.innerHTML = html;
                }
            }
        }

        const originalContents = new Map();

        function getOrSaveOriginalContent(el) {
            if (!originalContents.has(el)) {
                originalContents.set(el, el.innerHTML);
            }
            return originalContents.get(el);
        }

        function highlightText(html, query) {
            if (!query) return html;
            const regex = new RegExp('(' + query.replace(/[-\/\\^$*+?.()|[\]{}]/g, '\\$&') + ')', 'gi');
            
            const parser = new DOMParser();
            const doc = parser.parseFromString('<div>' + html + '</div>', 'text/html');
            const walk = document.createTreeWalker(doc.body, NodeFilter.SHOW_TEXT, null, false);
            let node;
            const nodesToReplace = [];
            while (node = walk.nextNode()) {
                if (node.nodeValue.toLowerCase().includes(query.toLowerCase())) {
                    nodesToReplace.push(node);
                }
            }
            nodesToReplace.forEach(textNode => {
                const span = doc.createElement('span');
                span.innerHTML = textNode.nodeValue.replace(regex, '<mark style="background-color: rgba(248, 238, 73, 0.4); color: inherit; border-radius: 2px; padding: 0 2px;">$1</mark>');
                textNode.parentNode.replaceChild(span, textNode);
            });
            return doc.body.firstChild.innerHTML;
        }

        function filterDocs() {
            const query = document.getElementById('search').value.toLowerCase().trim();
            const cards = document.querySelectorAll('.doc-card');
            const searchStatus = document.getElementById('search-status');
            
            let matchedCount = 0;

            const visibleRoutes = new Set();
            const visibleStructs = new Set();
            const visibleFunctions = new Set();

            cards.forEach(card => {
                const searchData = card.getAttribute('data-search').toLowerCase();
                const cardId = card.id;

                card.innerHTML = getOrSaveOriginalContent(card);

                if (query === '' || searchData.includes(query)) {
                    card.style.display = 'block';
                    matchedCount++;

                    if (query !== '') {
                        card.innerHTML = highlightText(card.innerHTML, query);
                    }

                    if (cardId.startsWith('route-')) {
                        visibleRoutes.add(cardId);
                    } else if (cardId.startsWith('struct-')) {
                        visibleStructs.add(cardId.replace('struct-', ''));
                    } else if (cardId.startsWith('fn-')) {
                        visibleFunctions.add(cardId.replace('fn-', ''));
                    }
                } else {
                    card.style.display = 'none';
                }
            });

            const navItems = document.querySelectorAll('.nav-item');
            navItems.forEach(item => {
                const onclickAttr = item.getAttribute('onclick') || '';
                let shouldShow = query === '';

                if (!shouldShow) {
                    if (onclickAttr.includes('route-')) {
                        const match = onclickAttr.match(/'([^']+)'/);
                        if (match && visibleRoutes.has(match[1])) {
                            shouldShow = true;
                        }
                    } else if (onclickAttr.includes('struct-')) {
                        const match = onclickAttr.match(/'struct-([^']+)'/);
                        if (match && visibleStructs.has(match[1])) {
                            shouldShow = true;
                        }
                    } else if (onclickAttr.includes('fn-')) {
                        const match = onclickAttr.match(/'fn-([^']+)'/);
                        if (match && visibleFunctions.has(match[1])) {
                            shouldShow = true;
                        }
                    }
                }

                item.style.display = shouldShow ? 'block' : 'none';
            });

            if (query !== '') {
                searchStatus.innerText = 'Found ' + matchedCount + ' matching items';
                searchStatus.style.display = 'block';
            } else {
                searchStatus.style.display = 'none';
            }
        }

        function generateMockJSON(structName) {
            const schema = srvSchemas[structName];
            if (!schema || !schema.fields || schema.fields.length === 0) {
                return null;
            }
            const mock = {};
            schema.fields.forEach(field => {
                let val = "";
                const lowerType = field.type.toLowerCase();
                if (lowerType.includes("int") || lowerType.includes("num") || lowerType.includes("float")) {
                    val = 0;
                } else if (lowerType.includes("bool")) {
                    val = false;
                } else if (lowerType.includes("array") || lowerType.startsWith("[")) {
                    val = [];
                } else if (lowerType.includes("map")) {
                    val = {};
                } else {
                    val = "string";
                }
                mock[field.name] = val;
            });
            return mock;
        }

        function getExampleCode(method, path, inputType, lang) {
            const mockObj = inputType ? generateMockJSON(inputType) : null;
            const jsonStr = mockObj ? JSON.stringify(mockObj, null, 2) : '';
            
            if (lang === 'curl') {
                let curl = "curl -X " + method + " \"http://localhost:8080" + path + "\"";
                if (mockObj) {
                    const escapedJson = JSON.stringify(mockObj).replace(/"/g, '\\"');
                    curl += " \\\n  -H \"Content-Type: application/json\" \\\n  -d \"" + escapedJson + "\"";
                }
                return curl;
            } else if (lang === 'go') {
                let goCode = "package main\n\nimport (\n\t\"bytes\"\n\t\"fmt\"\n\t\"io\"\n\t\"net/http\"\n)\n\nfunc main() {\n";
                if (mockObj) {
                    const tick = String.fromCharCode(96);
                    goCode += "\tjsonData := []byte(" + tick + JSON.stringify(mockObj, null, 2) + tick + ")\n";
                    goCode += "\treq, err := http.NewRequest(\"" + method + "\", \"http://localhost:8080" + path + "\", bytes.NewBuffer(jsonData))\n";
                    goCode += "\treq.Header.Set(\"Content-Type\", \"application/json\")\n";
                } else {
                    goCode += "\treq, err := http.NewRequest(\"" + method + "\", \"http://localhost:8080" + path + "\", nil)\n";
                }
                goCode += "\tif err != nil {\n\t\tpanic(err)\n\t}\n\n";
                goCode += "\tclient := &http.Client{}\n\tresp, err := client.Do(req)\n\tif err != nil {\n\t\tpanic(err)\n\t}\n\tdefer resp.Body.Close()\n\n\tbody, _ := io.ReadAll(resp.Body)\n\tfmt.Println(string(body))\n}";
                return goCode;
            } else if (lang === 'js') {
                let jsCode = "fetch(\"http://localhost:8080" + path + "\", {\n  method: \"" + method + "\",\n";
                if (mockObj) {
                    jsCode += "  headers: {\n    \"Content-Type\": \"application/json\"\n  },\n  body: JSON.stringify(" + JSON.stringify(mockObj, null, 2) + ")\n";
                } else {
                    jsCode += "  headers: {}\n";
                }
                jsCode += "})\n.then(response => response.json())\n.then(data => console.log(data))\n.catch(error => console.error(error));";
                return jsCode;
            }
            return '';
        }

        function initCodeExamples() {
            document.querySelectorAll('.code-examples-container').forEach(container => {
                const method = container.getAttribute('data-method');
                const path = container.getAttribute('data-path');
                const input = container.getAttribute('data-input');
                const codeBox = container.querySelector('.code-box');
                codeBox.textContent = getExampleCode(method, path, input, 'curl');
            });
        }
        
        function switchTab(btn, lang) {
            const container = btn.closest('.code-examples-container');
            const method = container.getAttribute('data-method');
            const path = container.getAttribute('data-path');
            const input = container.getAttribute('data-input');
            
            container.querySelectorAll('.tab-btn').forEach(b => b.classList.remove('active'));
            btn.classList.add('active');
            
            const codeBox = container.querySelector('.code-box');
            codeBox.textContent = getExampleCode(method, path, input, lang);
        }
        
        document.addEventListener('DOMContentLoaded', initCodeExamples);
    </script>
</body>
</html>`

type PageData struct {
	Title          string
	Doc            *SrvDoc
	SchemasJSON    string
	Versions       []string
	CurrentVersion string
}

func GenerateHtml(doc *SrvDoc, title, outputPath string, outDir string, currentVersion string) error {
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

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	schemasMap := make(map[string]StructDef)
	for _, s := range doc.Structs {
		schemasMap[s.Name] = s
	}
	schemasBytes, err := json.Marshal(schemasMap)
	if err != nil {
		return err
	}

	var versions []string
	if outDir != "" {
		entries, _ := os.ReadDir(outDir)
		for _, entry := range entries {
			if entry.IsDir() {
				versions = append(versions, entry.Name())
			}
		}
	}
	if currentVersion != "" {
		found := false
		for _, v := range versions {
			if v == currentVersion {
				found = true
				break
			}
		}
		if !found {
			versions = append(versions, currentVersion)
		}
	}

	data := PageData{
		Title:          title,
		Doc:            doc,
		SchemasJSON:    string(schemasBytes),
		Versions:       versions,
		CurrentVersion: currentVersion,
	}

	return tmpl.Execute(file, data)
}
