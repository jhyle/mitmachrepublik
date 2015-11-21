<div class="col-md-3 col-sm-4 col-xs-6 col-tile">
	<div class="tile">
		<a href="{{.Url}}" style="display:block" title="Infos zu {{.Title}} anschauen">
		{{ if .Image }}
			<!-- {{if len .Categories}}{{with index .Categories 0}}<div class="small-icon"><span class="fa fa-{{categoryIcon .}} fa-fw" title="{{categoryTitle .}}"></span></div>{{end}}{{end}} -->
			<div class="tile-image" style="background-image: url(/bild/{{.Image}}?height=165)"> </div>
		{{ end }}
		<div class="tile-text">
			<h3>{{.Title}}</h3>
			<p class="datetime">{{longDatetimeFormat .Start}}</p>
			<p class="place">{{if .Addr.Name}}{{.Addr.Name}}{{if .Addr.City}}, {{end}}{{end}}{{citypartName .Addr}}</p>
			<p class="description">{{strClip .PlainDescription 150}}</p>
			<p class="highlight" style="position: absolute; bottom: 11px"><span class="fa fa-caret-right"></span> Veranstaltung ansehen</p>
		</div>
		</a>
	</div>
</div>
