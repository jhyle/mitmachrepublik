<div class="row tiles">
	<div class="col-xs-1">&nbsp;</div>
	<div class="col-xs-10"><h1>Veranstaltung {{.event.Title}}</h1></div>
	<div class="col-xs-1">&nbsp;</div>
</div>
<div class="row tiles">
	<div class="col-xs-1">&nbsp;</div>
	<div class="col-xs-3 col-box">
		{{template "organizer_box.tpl" .organizer}}
	</div>
	{{with .event}}<div class="col-xs-7">
		{{if .Image}}
			{{if len .Categories}}{{with index .Categories 0}}<div class="small-icon pull-left"><span class="fa fa-{{categoryIcon .}} fa-fw" title="{{categoryTitle .}}"></span></div>{{end}}{{end}}
			<img class="pull-left" style="margin-right: 10px; margin-bottom: 10px" src="/bild/{{.Image}}?width=340&height=255" title="{{.Title}}">
		{{end}}
		<p class="small-icon pull-left"><span class="fa fa-calendar" title="Datum"></span></p>
		<p class="date">{{dateFormat .Start}}</p>
		<p class="small-icon pull-left"><span class="fa fa-clock-o" title="Uhrzeit"></span></p>
		<p class="date">{{timeFormat .Start}}</p>
		{{ if not .Addr.IsEmpty }}
			<p class="small-icon pull-left"><span class="fa fa-map-marker" title="Ort"></span></p>
			<p>{{ if .Addr.Name }}{{.Addr.Name}}<br />{{ end }}{{ if .Addr.Street }}{{.Addr.Street}}, {{ end }}{{ if .Addr.Pcode }}{{.Addr.Pcode}} {{ end }}{{.Addr.City}}</p>
		{{ end }}
		<div class="clearfix"></div>
		<p>{{.Descr}}</p>
		{{if .Web}}
			<p><a href="{{.Web}}" class="btn btn-mmr" style="margin: 0" target="_blank">Zur Veranstaltungs-Webseite</a></p>
		{{end}}
	</div>{{end}}
	<div class="col-xs-1">&nbsp;</div>
</div>
