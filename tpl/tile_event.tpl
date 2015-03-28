<div class="col-xs-2 col-tile">
	<a href="/veranstaltung/{{eventUrl .}}" style="display:block">
		{{ if .Image }}
			<div class="small-icon"><span class="fa fa-{{if len .Categories}}{{with index .Categories 0}}{{categoryIcon .}}{{end}}{{end}} fa-fw"></span></div>
			<img src="/bild/{{.Image}}?width=220&height=165">
		{{ end }}
		<div class="tile-text">
			<h3>{{.Title}}</h3>
			<p class="datetime">{{dateFormat .Start}}</p>
			<p>{{strClip .Descr 80}}</p>
			<p class="highlight"><span class="fa fa-caret-right"></span> Veranstaltung ansehen</p>
		</div>
	</a>
</div>
