<div class="col-xs-2 col-tile">
	<a href="{{eventUrl .}}" style="display:block">
		{{ if .Image }}
			{{if len .Categories}}{{with index .Categories 0}}<div class="small-icon"><span class="fa fa-{{categoryIcon .}} fa-fw" title="{{categoryTitle .}}"></span></div>{{end}}{{end}}
			<img src="/bild/{{.Image}}?width=220&height=165" title="{{.Title}}">
		{{ end }}
		<div class="tile-text">
			<h3>{{.Title}}</h3>
			<p class="datetime">{{datetimeFormat .Start}} Uhr - {{citypartName .Addr}}</p>
			<p>{{strClip .Descr 80}}</p>
			<p class="highlight"><span class="fa fa-caret-right"></span> Veranstaltung ansehen</p>
		</div>
	</a>
</div>
