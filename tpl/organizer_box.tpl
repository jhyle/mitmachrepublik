<div class="box">
{{with (or .user .organizer)}}
<a href="{{if $.user}}/veranstalter/verwaltung/0{{else}}{{.Url}}{{end}}" title="Veranstaltungen von {{.Name}} anzeigen">
	<h2>{{.Name}}</h2>
	{{if .Image}}
		<div class="tile-image" style="background-image: url(/bild/{{.Image}}?height=165)"> </div>
	{{end}}
	{{if .ImageCredit}}
		<p class="credits">{{.ImageCredit}}</p>
	{{end}}
</a>
<div class="tile-text">
{{ if not .Addr.IsEmpty }}
	<p class="small-icon pull-left"><span class="fa fa-map-marker fa-fw" title="Ort"></span></p>
	<p>{{ if .Addr.Name }}<span>{{.Addr.Name}}</span><br />{{ end }}<span class="address">{{ if .Addr.Street }}<span>{{.Addr.Street}}</span>, {{ end }}{{ if .Addr.Pcode }}<span>{{.Addr.Pcode}}</span> {{ end }}<span>{{.Addr.City}}</span></span></p>
{{ end }}
<a href="{{if $.user}}/veranstalter/verwaltung/0{{else}}{{.Url}}{{end}}" title="Veranstaltungen von {{.Name}} anzeigen">
	<p>{{.HtmlDescription}}</p>
</a>
	<p><a href="{{if $.user}}/veranstalter/verwaltung/0{{else}}{{.Url}}{{end}}" title="Veranstaltungen von {{.Name}} anzeigen" class="highlight"><span class="fa fa-caret-right"></span> Alle Veranstaltungen</a></p>
{{if .Web}}
	<p><a href="{{.Web}}" title="Webseite von {{.Name}} anzeigen" target="_blank" class="highlight"><span class="fa fa-caret-right"></span> {{.Web}}</a></p>
{{end}}
</div>
{{ if not .Addr.IsEmpty }}
	<a href="https://maps.google.de/maps?hl=de&q={{.Addr.Name}}+{{.Addr.Street}}+{{.Addr.Pcode}}+{{.Addr.City}}&ie=UTF8" target="_blank" title="In Google Maps öffnen">
		<div class="tile-image" style="background-image: url(https://maps.googleapis.com/maps/api/staticmap?center={{.Addr.Name}}+{{.Addr.Street}}+{{.Addr.Pcode}}+{{.Addr.City}}&markers={{.Addr.Street}}+{{.Addr.Pcode}}+{{.Addr.City}}&zoom=15&size=270x165&key={{$.googleApiKey}})"> </div>
	</a>
{{end}}
{{end}}
</div>