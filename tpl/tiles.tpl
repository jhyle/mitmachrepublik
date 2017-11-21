<div class="row tiles">
{{range $path, $topic := .topics}}
	{{if eq $path "demonstrationen-und-politik"}}
		<div class="col-xs-12 row-break">
			<a href="/dialog/login" data-href="/dialog/login" rel="nofollow" data-toggle="modal" data-target="#login">Veröffentliche Deine <span>Veranstaltungen</span> auf mitmachrepublik.de!</a> <button type="button" class="btn btn-mmr" style="display:inline-block" href="/dialog/login" data-href="/dialog/login" rel="nofollow" data-toggle="modal" data-target="#login" title="Melde Dich an, um Deine Veranstaltungen einzutragen."> Eintragen</button>
		</div>
	{{end}}
	{{if gt (len (index $.events $path)) 0}}
		<div class="col-xs-12 row-topic">
			<a href="{{$path}}">{{$topic.Name}}</a> <a class="highlight small" href="{{$path}}"><span class="fa fa-caret-right"></span> alle Veranstaltungen</a>
		</div>
	{{end}}
	{{range $event := index $.events $path}}
		{{if (index $.organizers .OrganizerId).Approved}}
			<div class="col-md-3 col-sm-4 col-xs-6 col-tile">
				<div class="tile">
					<a href="{{.Url}}?from={{(index (index $.timespans 0) 0).Unix}}" style="display:block" title="Infos zu {{.Title}} anschauen">
					{{if or (.Image) ((index $.organizers .OrganizerId).Image)}}
						<!-- {{if len .Categories}}{{with index .Categories 0}}<div class="small-icon"><span class="fa fa-{{categoryIcon .}} fa-fw" title="{{categoryTitle .}}"></span></div>{{end}}{{end}} -->
						<div class="tile-image" style="background-image: url(/bild/{{if .Image}}{{.Image}}{{else}}{{(index $.organizers .OrganizerId).Image}}{{end}}?height=165)"> </div>
					{{ end }}
					<div class="tile-text">
						<h3>{{.Title}}</h3>
						<p class="datetime">{{if .Recurrence}}{{.Recurrence}}{{else}}{{longDatetimeFormat (.NextDate (index (index $.timespans 0) 0))}}{{end}}</p>
						{{if $.organizers}}{{if ne ((index $.organizers .OrganizerId).Name) ("Mitmach-Republik")}}<p class="datetime">{{(index $.organizers .OrganizerId).Name}}</p>{{end}}{{end}}
						<p class="place">{{if .Addr.Name}}{{.Addr.Name}}{{if .Addr.City}}, {{end}}{{end}}{{citypartName .Addr}}</p>
						<p class="description">{{strClip .PlainDescription 150}}</p>
						<p class="highlight" style="position: absolute; bottom: 11px"><span class="fa fa-caret-right"></span> Veranstaltung ansehen</p>
					</div>
					</a>
				</div>
			</div>
		{{end}}
	{{end}}
{{end}}
</div>
<div class="row tiles">
	<div class="col-sm-3 col-xs-2">&nbsp;</div>
	<div class="col-sm-6 col-xs-8">
		<a href="/veranstaltungen/{{simpleEventSearchUrl ""}}" class="btn btn-mmr" style="width: 100%">Mehr Veranstaltungen</a>
	</div>
	<div class="col-sm-3 col-xs-2">&nbsp;</div>
</div>
