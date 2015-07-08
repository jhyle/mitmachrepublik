<!DOCTYPE html><html lang="de">
<head>
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<meta charset="utf-8">
	<title>Veranstaltungen aus der Mitmach-Republik</title>
	<link href="http://fonts.googleapis.com/css?family=Open+Sans:400,300,700,600" rel="stylesheet" type="text/css">
</head>
<body style="font-family: 'Open Sans',sans-serif; font-size: 15px; line-height: 1.42857; width: 650px; margin-left: auto; margin-right: auto">
{{$n := len .events}}
<img style="margin: 10px 30px 0 0; float: left !important; vertical-align: middle; border: 0 none; width: 200px" src="http://{{$.hostname}}/images/mitmachrepublik.png" alt="Logo Mitmach-Republik">
<h1 style="color: #ff5100; font-size: 27px; margin: 10px 0 20px 0; vertical-align: middle; line-height: 2.5em">Deine Veranstaltungen</h1> 
{{range .events}}
<div style="border: 1px solid #e6e6e6; margin-bottom: 15px; overflow: hidden; ">
	<a style="color: #2f3030; text-decoration: none;" href="http://{{$.hostname}}{{.Url}}" title="Veranstaltung anzeigen">
	{{if .Image}}
		<img style="margin-right: 10px; float: left !important; vertical-align: middle; border: 0 none; " src="http://{{$.hostname}}/bild/{{.Image}}?width=220&height=165" alt="Veranstaltung {{.Title}}">
	{{end}}
	<div style="margin: 10px;">
		<h3 style="color: #ff5100; font-size: 23px; font-weight: lighter; margin: 10px 0; line-height: 1.1;">{{.Title}}</h3>
		<p style="margin: 0 0 3px 0; font-size: 13px; font-weight: bold;">{{datetimeFormat .Start}}{{if dateFormat .End}}<span> bis {{if eq (dateFormat .Start) (dateFormat .End)}}{{timeFormat .End}}{{else}}{{datetimeFormat .End}}{{end}}</span>{{end}} {{if $.organizerNames}}{{if ne (index $.organizerNames .OrganizerId) ("Mitmach-Republik")}} - {{index $.organizerNames .OrganizerId}}{{end}}{{end}}</p>
		<p style="margin: 0 0 3px 0;">{{strClip .PlainDescription 100}}</p>
		{{ if not .Addr.IsEmpty }}
			<p style="margin: 0 0 3px 0; float: left !important; color: #7a7d7d; font-size: 13px;">{{ if .Addr.Name }}<span>{{.Addr.Name}}</span><br />{{ end }}<span class="address">{{ if .Addr.Street }}<span>{{.Addr.Street}}</span>, {{ end }}{{ if .Addr.Pcode }}<span>{{.Addr.Pcode}}</span> {{ end }}<span>{{citypartName .Addr}}</span></span></p>
		{{ end }}
	</div>
	</a>
</div>
{{end}}
<div style="background-color: #303030; color: white; font-size: 13px; line-height: 2; padding: 30px">
	@ <a style="color: white; font-weight: bold; text-decoration: none;" href="http://www.mitmachrepublik">www.mitmachrepublik.de</a> | Klicke <a style="color: white; font-weight: bold; text-decoration: none;" href="http://{{.hostname}}/newsletter/unsubscribe/{{.alertId}}">hier</a>, wenn Du diese E-Mail nicht mehr erhalten möchtest. 
</div>
</body>
</html>