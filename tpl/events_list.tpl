{{if not $.user }}{{if not $.organizer}}
<div style="margin-bottom: 10px; font-weight: bolder">
	<a class="highlight" href="javascript:void(0)" data-href="/dialog/emailalert/{{eventSearchUrlWithQuery .place .targetIds .categoryIds .dateIds .radius .query}}" rel="nofollow" data-toggle="modal" data-target="#email-alert" title="Wir senden Dir die Ergebnisse dieser Suche per E-Mail."><span class="fa fa-caret-right"></span> Lass Dich per E-Mail über diese Suche informieren.</a>
</div>
{{end}}{{end}}
{{$n := len .events}}
{{if eq $n 0}}
<div class="row-tile">
	<p class="big-text text-center" style="padding-top: 20px">Es wurden keine Veranstaltungen gefunden.</p>
</div>
{{else}}
{{range .events}}
<div class="row-tile">
	{{if not $.user }}
		<a href="{{.Url}}" title="Infos zu {{.Title}} anschauen">
	{{end}}
	{{if .Image}}
		<!-- {{if len .Categories}}{{with index .Categories 0}}<div class="small-icon"><span class="fa fa-{{categoryIcon .}} fa-fw" title="{{categoryTitle .}}"></span></div>{{end}}{{end}} -->
		<img class="pull-left" style="margin-right: 10px" src="/bild/{{.Image}}?width=220&height=165" alt="Veranstaltung {{.Title}}">
	{{end}}
	<div class="tile-text">
		{{ if $.user }}
			<p class="pull-right"><a href="#" name="delete-event" title="Löschen" data-target="{{.Id.Hex}}" class="close"><span class="fa fa-times"></span></a></p>
		{{ end }}
		<h3>{{.Title}}</h3>
		<p class="datetime">{{datetimeFormat .Start}}{{if dateFormat .End}}<span> bis {{if eq (dateFormat .Start) (dateFormat .End)}}{{timeFormat .End}}{{else}}{{datetimeFormat .End}}{{end}}</span>{{end}} {{if $.organizerNames}}{{if ne (index $.organizerNames .OrganizerId) ("Mitmach-Republik")}} - {{index $.organizerNames .OrganizerId}}{{end}}{{end}}</p>
		{{ if $.user }}
			<p class="pull-right"><a href="/veranstalter/verwaltung/veranstaltung/{{.Id.Hex}}" class="btn btn-mmr" style="margin: 0; width: 100px">Bearbeiten</a></p>
		{{end}}
		<p>{{strClip .PlainDescription 100}}</p>
		{{ if not .Addr.IsEmpty }}
			<p class="small-icon pull-left"><span class="fa fa-map-marker fa-fw" title="Ort"></span></p>
			<p class="pull-left place">{{ if .Addr.Name }}<span>{{.Addr.Name}}</span><br />{{ end }}<span class="address">{{ if .Addr.Street }}<span>{{.Addr.Street}}</span>, {{ end }}{{ if .Addr.Pcode }}<span>{{.Addr.Pcode}}</span> {{ end }}<span>{{citypartName .Addr}}</span></span></p>
		{{ end }}
		{{ if $.user }}
			<p class="pull-right"><a href="/veranstalter/verwaltung/veranstaltung?copy={{.Id.Hex}}" class="btn btn-mmr" style="margin: 0; width: 100px">Kopieren</a></p>
		{{ end }}
	</div>
	{{if not $.user }}
		</a>
	{{end}}
</div>
{{end}}
<div class="pages">
	<a href="./0{{if $.query}}?query={{$.query}}{{end}}#events" title="Zum Anfang der Liste"><div class="page">Anfang</div></a>
	{{if gt $.page 0}}<a href="./{{dec $.page}}{{if $.query}}?query={{$.query}}{{end}}#events" title="Vorherige Seite">{{end}}<div class="page">&lt;</div>{{if gt $.page 0}}</a>{{end}}
	{{range $.pages}}
		{{if or (and (ge . (dec (dec $.page))) (le . (inc (inc $.page)))) (or (eq $.page .) (or (le . 1) (ge . (dec $.maxPage))))}}
			<a href="./{{.}}{{if $.query}}?query={{$.query}}{{end}}#events" title="Zu Seite {{inc .}}"><div class="page {{if eq . $.page}}cur-page{{end}}">{{inc .}}</div></a>
		{{else}}
			{{if or (eq . (dec (dec (dec $.page)))) (eq . (inc (inc (inc $.page))))}}
				<div class="page">..</div>
			{{end}}
		{{end}}
	{{end}}
	{{if lt $.page $.maxPage}}<a href="./{{inc $.page}}{{if $.query}}?query={{$.query}}{{end}}#events" title="Nächste Seite">{{end}}<div class="page">&gt;</div>{{if lt $.page $.maxPage}}</a>{{end}}
	<a href="./{{$.maxPage}}{{if $.query}}?query={{$.query}}{{end}}#events" title="Ans Ende der Liste"><div class="page">Ende</div></a>
</div>
{{end}}
