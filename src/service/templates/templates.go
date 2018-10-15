package templates

// Index is the template for the info2 route
var Index = `
<html>
<head>
    <title>EVM-LITE</title>
</head>
<body>
    <div>
        <h3>Consensus Info</h1>
        <div>
            <ul>
            {{ range $key, $value := . }}
                <li><strong>{{ $key }}</strong>: {{ $value }}</li>
            {{ end }}
            </ul>
        </div>
    </div>
</body>
</html>`
