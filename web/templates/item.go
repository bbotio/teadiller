package templates

const Item string = `<html>
<head>
<meta charset="utf-8">
        <title>{{.Name}}</title>
        <meta name="description" content="{{.Desc}}">    
        <meta property="author" content="Teadiller" />
        <meta name="viewport" content="width=device-width">

        <meta property="og:type" content="article" />
        <meta property="og:title" content="{{.Name}}" />
        <meta property="og:description" content="{{.Desc}}" />
        <meta propert="og:image" content="{{.PhotoPath}}"/>
        <meta property="og:site_name" content="teadiller.com" />


        <meta property="twitter:card" content="summary_large_image" />
        <meta property="twitter:title" content="{{.Name}}" />
        <meta property="twitter:description" content="{{.Desc}}" />
        <meta propert="twitter:image" content="{{.PhotoPath}}"/>
</head>
<body>
{{.Name}} - {{.Desc}}
</body>
</html>`
