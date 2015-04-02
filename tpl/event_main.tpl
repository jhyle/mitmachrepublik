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
		<div class="pull-left" style="margin-right: 5px; margin-bottom: 10px">
			{{if .Image}}
				<img class="" style="margin-right: 10px; margin-bottom: 15px" src="/bild/{{.Image}}?width=300" title="{{.Title}}">
			{{end}}
			{{ if not .Addr.IsEmpty }}
				<div><iframe width="300" height="225" src="http://maps.google.de/maps?hl=de&q={{.Addr.Street}}%20{{.Addr.Pcode}}%20{{.Addr.City}}&ie=UTF8&t=&z=14&iwloc=B&output=embed" frameborder="0" scrolling="no" marginheight="0" marginwidth="0"></iframe></div>
			{{end}}
		</div>
		<p class="small-icon pull-left"><span class="fa fa-calendar fa-fw" title="Datum"></span></p>
		<p class="date">{{dateFormat .Start}}</p>
		<p class="small-icon pull-left"><span class="fa fa-clock-o fa-fw" title="Uhrzeit"></span></p>
		<p class="date">{{timeFormat .Start}}{{if timeFormat .End}} bis{{if eq (dateFormat .Start) (dateFormat .End)}} {{timeFormat .End}}{{end}}{{end}} Uhr</p>
		{{if dateFormat .End}}{{if ne (dateFormat .Start) (dateFormat .End)}}
			<p class="small-icon pull-left"><span class="fa fa-calendar fa-fw" title="Enddatum"></span></p>
			<p class="date">{{dateFormat .End}}</p>
			<p class="small-icon pull-left"><span class="fa fa-clock-o fa-fw" title="Uhrzeit"></span></p>
			<p class="date">{{timeFormat .End}} Uhr</p>
		{{end}}{{end}}
		{{ if not .Addr.IsEmpty }}
			<p class="small-icon pull-left"><span class="fa fa-map-marker fa-fw" title="Ort"></span></p>
			<p>{{ if .Addr.Name }}{{.Addr.Name}}<br />{{ end }}{{ if .Addr.Street }}{{.Addr.Street}}, {{ end }}{{ if .Addr.Pcode }}{{.Addr.Pcode}} {{ end }}{{.Addr.City}}</p>
		{{ end }}
		{{if len .Categories}}{{with index .Categories 0}}
			<p class="small-icon pull-left"><span class="fa fa-{{categoryIcon .}} fa-fw" title="{{categoryTitle .}}"></span></p>{{end}}
			<p>{{range $i, $category := .Categories}}{{if $i}}, {{end}}{{categoryTitle $category}}{{end}}</p>
		{{end}}
		<p>{{.Descr}}</p>
		<div class="clearfix"></div>
		{{if .Web}}
			<p><a href="{{.Web}}" class="btn btn-mmr" style="margin: 0" target="_blank">Zur Veranstaltungs-Webseite</a></p>
		{{end}}
	</div>{{end}}
	<div class="col-xs-1">&nbsp;</div>
</div>
