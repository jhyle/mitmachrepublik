<!DOCTYPE html><html lang="de">
<head>
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<meta charset="utf-8">
	<title>Veranstaltungen aus der Mitmach-Republik</title>
	<link href="http://fonts.googleapis.com/css?family=Open+Sans:400,300,700,600" rel="stylesheet" type="text/css">
</head>
<body style="font-family: 'Open Sans',sans-serif; font-size: 15px; line-height: 1.42857; width: 650px; margin-left: auto; margin-right: auto">
{{$n := len .events}}
<a style="float: left !important; text-decoration: none;" href="http://{{$.hostname}}/" title="www.mitmachrepublik.de">
	<img style="margin: 13px 30px 10px 0; vertical-align: middle; border: 0 none; width: 200px" src="http://{{$.hostname}}/images/mitmachrepublik.png" alt="Logo Mitmach-Republik">
</a>
<h1 style="color: #ff5100; font-size: 27px; margin: 10px 0 0 0; line-height: 2em">Deine Veranstaltungen</h1>
<p style="margin: 0 0 10px 0">Klicke auf die Veranstaltungen für weitere Infos.</p>  
{{range .events}}
{{if (index $.organizers .OrganizerId).Approved}}
<div style="border: 1px solid #e6e6e6; margin-bottom: 15px; overflow: hidden; ">
	<a style="color: #2f3030; text-decoration: none;" href="http://{{$.hostname}}{{.Url}}" title="Veranstaltung anzeigen">
	{{if or (.Image) ((index $.organizers .OrganizerId).Image)}}
		<img style="margin-right: 10px; float: left !important; vertical-align: middle; border: 0 none; " src="http://{{$.hostname}}/bild/{{if .Image}}{{.Image}}{{else}}{{(index $.organizers .OrganizerId).Image}}{{end}}?width=220&height=165" alt="Veranstaltung {{.Title}}">
	{{end}}
	<div style="margin: 10px;">
		<h3 style="color: #ff5100; font-size: 23px; font-weight: lighter; margin: 10px 0; line-height: 1.1;">{{.Title}}</h3>
		<p style="margin: 0 0 3px 0; font-size: 13px; font-weight: bold;">{{longDatetimeFormat (.NextDate (index (index $.timespans 0) 0))}}{{if dateFormat .End}}<span> bis {{if eq (dateFormat .Start) (dateFormat .End)}}{{timeFormat .End}}{{else}}{{datetimeFormat .End}}{{end}}</span>{{end}} {{if $.organizers}}{{if ne ((index $.organizers .OrganizerId).Name) ("Mitmach-Republik")}} - {{(index $.organizers .OrganizerId).Name}}{{end}}{{end}}</p>
		<p style="margin: 0 0 3px 0;">{{strClip .PlainDescription 100}}</p>
		{{ if not .Addr.IsEmpty }}
			<p style="margin: 0 0 3px 0; float: left !important; color: #7a7d7d; font-size: 13px;">{{ if .Addr.Name }}<span>{{.Addr.Name}}</span><br />{{ end }}<span class="address">{{ if .Addr.Street }}<span>{{.Addr.Street}}</span>, {{ end }}{{ if .Addr.Pcode }}<span>{{.Addr.Pcode}}</span> {{ end }}<span>{{citypartName .Addr}}</span></span></p>
		{{ end }}
	</div>
	</a>
</div>
{{end}}
{{end}}
<div style="background-color: #303030; color: white; font-size: 13px; line-height: 2; padding: 30px">
	@ <a style="color: white; font-weight: bold; text-decoration: none;" href="http://www.mitmachrepublik">www.mitmachrepublik.de</a> | Klicke <a style="color: white; font-weight: bold; text-decoration: none;" href="http://{{.hostname}}/newsletter/unsubscribe/{{.alertId}}">hier</a>, wenn Du diese E-Mail nicht mehr erhalten möchtest. 
</div>
</body>
</html>