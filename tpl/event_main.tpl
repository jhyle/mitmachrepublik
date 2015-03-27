<div class="row tiles">
	<div class="col-sm-1 hidden-xs">&nbsp;</div>
	<div class="col-sm-10"><h1>{{.event.Title}}</h1></div>
	<div class="col-sm-1 hidden-xs">&nbsp;</div>
</div>
<div class="row tiles">
	<div class="col-sm-1 hidden-xs">&nbsp;</div>
	<div class="col-sm-3 col-box">
		{{template "organizer.tpl" .organizer}}
	</div>
	{{with .event}}<div class="col-sm-7">
		{{if .Image}}
			<div class="small-icon pull-left"><span class="fa fa-{{if len .Categories}}{{with index .Categories 0}}{{categoryIcon .}}{{end}}{{end}} fa-fw"></span></div>
			<img class="pull-left" style="margin-right: 10px; margin-bottom: 10px" src="/bild/{{.Image}}?width=340&height=255">
		{{end}}
		<p class="small-icon pull-left"><span class="fa fa-calendar"></span></p>
		<p class="date">{{dateFormat .Start}}</p>
		<p class="small-icon pull-left"><span class="fa fa-clock-o"></span></p>
		<p class="date">{{timeFormat .Start}}</p>
		{{ if not .Addr.IsEmpty }}
			<p class="small-icon pull-left"><span class="fa fa-map-marker"></span></p>
			<p>{{ if .Addr.Name }}{{.Addr.Name}}<br />{{ end }}{{ if .Addr.Street }}{{.Addr.Street}}, {{ end }}{{ if .Addr.Pcode }}{{.Addr.Pcode}} {{ end }}{{.Addr.City}}</p>
		{{ end }}
		<div class="clearfix"></div>
		<p>{{.Descr}}</p>
		{{if .Web}}
			<p><a href="{{.Web}}" class="btn btn-mmr" style="margin: 0" target="_blank">Zur Veranstaltungs-Webseite</a></p>
		{{end}}
	</div>{{end}}
	<div class="col-sm-1 hidden-xs">&nbsp;</div>
</div>
