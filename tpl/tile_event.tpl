<div class="col-sm-2 col-tile">
	<a href="#">
		{{ if .Image }}
			<div class="small-icon"><span class="fa fa-{{with index .Categories 0}}{{categoryIcon .}}{{end}} fa-fw"></span></div>
			<img class="img-responsive" src="/bild/{{.Image}}?width=220&height=165">
		{{ end }}
		<div class="tile-text">
			<h3>{{.Title}}</h3>
			<p class="datetime">{{dateFormat .Start}}</p>
			<p>{{strClip .Descr 80}}</p>
			{{/* if not .Addr.IsEmpty }}
				<p class="place-icon pull-left"><span class="fa fa-map-marker"></span></p>
				<p class="pull-left place" style="width: 80%">{{ if .Addr.Name }}{{.Addr.Name}}, {{ end }}{{ if .Addr.Street }}{{.Addr.Street}}, {{ end }}{{ if .Addr.Pcode }}{{.Addr.Pcode}} {{ end }}{{.Addr.City}}</p>
				<div class="clearfix"></div>
			{{ end */}}
			<p class="highlight"><span class="fa fa-caret-right"></span> Veranstaltung ansehen</p>
		</div>
	</a>
</div>
