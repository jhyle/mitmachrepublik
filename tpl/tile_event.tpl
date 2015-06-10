<div class="col-xs-2 col-tile">
	<a href="{{.Url}}" style="display:block" title="Veranstaltung anzeigen">
		{{ if .Image }}
			<!-- {{if len .Categories}}{{with index .Categories 0}}<div class="small-icon"><span class="fa fa-{{categoryIcon .}} fa-fw" title="{{categoryTitle .}}"></span></div>{{end}}{{end}} -->
			<div style="min-height: 140px">
				<img src="/bild/{{.Image}}?width=220&height=165" alt="Veranstaltung {{.Title}}">
			</div>
		{{ end }}
		<div class="tile-text">
			<h3>{{.Title}}</h3>
			<p class="datetime">{{datetimeFormat .Start}}</p>
			<p class="place">{{if .Addr.Name}}{{.Addr.Name}}{{if .Addr.City}}, {{end}}{{end}}{{citypartName .Addr}}</p>
			<p class="description">{{strClip .PlainDescription 150}}</p>
			<p class="highlight" style="position: absolute; bottom: 0"><span class="fa fa-caret-right"></span> Veranstaltung ansehen</p>
		</div>
	</a>
</div>
