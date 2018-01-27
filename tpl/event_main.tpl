<div class="row tiles">
	<div class="col-xs-12"><h1>{{.event.Title}}{{if $.place}} in {{$.place}}{{end}}</h1></div>
</div>
<div class="row tiles">
	<div class="col-lg-3 col-sm-4 col-xs-5 col-box">
		{{template "organizer_box.tpl" .}}
	</div>
	{{with .event}}<div class="col-lg-9 col-sm-8 col-xs-7">
		<div class="event-images pull-left">
			{{if .Image}}
				<a href="{{.Web}}" title="Webseite der Veranstaltung aufrufen" target="_blank">
					<img width="450" style="margin-right: 10px; margin-bottom: 15px" src="/bild/{{.Image}}?width=450" alt="Veranstaltung {{.Title}}">
				</a>
			{{end}}
			{{if .ImageCredit}}
				<div class="credits">{{.ImageCredit}}</div>
			{{end}}
			{{ if not .Addr.IsEmpty }}
			<a href="https://maps.google.de/maps?hl=de&q={{.Addr.Name}}+{{.Addr.Street}}+{{.Addr.Pcode}}+{{.Addr.City}}&ie=UTF8" target="_blank" title="In Google Maps öffnen">
				<img width="450" height="225" style="display: block" src="https://maps.googleapis.com/maps/api/staticmap?center={{.Addr.Name}}+{{.Addr.Street}}+{{.Addr.Pcode}}+{{.Addr.City}}&markers={{.Addr.Street}}+{{.Addr.Pcode}}+{{.Addr.City}}&zoom=15&size=450x225&key={{$.googleApiKey}}" title="in Google Maps öffnen" alt="Karte">
			</a>
			{{end}}
		</div>
		<div class="social">
			<a style="margin-right: 10px" href="https://www.facebook.com/sharer/sharer.php?u=http%3A%2F%2F{{$.hostname}}{{.Url}}%3Ffrom%3D{{$.start.Unix}}" target="_blank"><img src="/images/facebook_share.png"></a>
			<a style="margin-right: 10px" href="https://plus.google.com/share?url=http%3A%2F%2F{{$.hostname}}{{.Url}}%3Ffrom%3D{{$.start.Unix}}" target="_blank"><img src="/images/google_share.png"></a>
			<a href="https://twitter.com/intent/tweet?url=http%3A%2F%2F{{$.hostname}}{{.Url}}%3Ffrom%3D{{$.start.Unix}}" target="_blank"><img src="/images/twitter_share.png"></a>
			<div class="recommend"><a id="event-mail" title="Empfehle die Veranstaltung per E-Mail" class="highlight" href="javascript:void(0)" data-href="/dialog/sendevent/{{.Id.Hex}}?from={{$.start.Unix}}" rel="nofollow" data-toggle="modal" data-target="#share"><span class="fa fa-envelope"></span> Empfehlen</a></div>
		</div>
		<p class="small-icon pull-left"><span class="fa fa-calendar fa-fw" title="Datum"></span></p>
		{{if $.showDate}}
			<p class="icon-text date">{{dateFormat $.start}}</p>
			{{if ne (timeFormat $.start) ("00:00")}}
				<p class="small-icon pull-left"><span class="fa fa-clock-o fa-fw" title="Uhrzeit"></span></p>
				<p class="icon-text date">{{timeFormat $.start}}{{if timeFormat $.end}}{{if eq (dateFormat $.start) (dateFormat $.end)}}{{if ne (timeFormat $.end) (timeFormat $.start)}} bis {{timeFormat $.end}}{{end}}{{end}}{{end}} Uhr</p>
			{{end}}
			{{if dateFormat $.end}}{{if ne (dateFormat $.start) (dateFormat $.end)}}
				<p class="small-icon pull-left"><span class="fa fa-calendar fa-fw" title="Enddatum"></span></p>
				<p class="icon-text date">{{dateFormat $.end}}</p>
				{{if ne (timeFormat $.end) ("23:59")}}
					<p class="small-icon pull-left"><span class="fa fa-clock-o fa-fw" title="Uhrzeit"></span></p>
					<p class="icon-text date">{{timeFormat $.end}} Uhr</p>
				{{end}}
			{{end}}{{end}}
			{{if gt .Recurrency 0}}
				<p class="small-icon pull-left"><span class="fa fa-repeat fa-fw" title="Wiederholungen"></span></p>
				<p class="icon-text date">{{.Recurrence}}</p>
			{{end}}
		{{else}}
			<p class="icon-text date">findet zurzeit nicht statt</p>
		{{end}}
		{{if not .Addr.IsEmpty}}
			<p class="small-icon pull-left"><span class="fa fa-map-marker fa-fw" title="Ort"></span></p>
			<p class="icon-text">{{ if .Addr.Name }}<span>{{.Addr.Name}}</span><br />{{ end }}<span class="address">{{ if .Addr.Street }}<span>{{.Addr.Street}}</span>, {{ end }}{{ if .Addr.Pcode }}<span>{{.Addr.Pcode}}</span> {{ end }}<span>{{.Addr.City}}</span></span></p>
		{{end}}
		{{if len .Targets}}
			<p class="small-icon pull-left"><span class="fa fa-child fa-fw" title="Zielgruppen"></span></p>
			<p class="icon-text">{{range $i, $target := .Targets}}{{if $i}}, {{end}}{{targetTitle $target}}{{end}}</p>
		{{end}}
		{{if len .Categories}}{{with index .Categories 0}}
			<p class="small-icon pull-left"><span class="fa fa-{{categoryIcon .}} fa-fw" title="Kategorien"></span></p>{{end}}
			<p class="icon-text">{{range $i, $category := .Categories}}{{if $i}}, {{end}}{{categoryTitle $category}}{{end}}</p>
		{{end}}
		{{if .Web}}
			<p class="small-icon pull-left"><span class="fa fa-external-link fa-fw" title="Webseite"></span></p>
			<p class="icon-text date"><a href="{{.Web}}" class="highlight" title="Webseite von {{.Title}}" target="_blank">{{strClip .Web 30}}</a></p>
		{{end}}
		{{if $.date}}
			<p style="font-weight: bolder">{{if .Rsvp}}Anmeldung erforderlich!{{if .Web}} Melde Dich auf der Webseite der Veranstaltung an.{{end}}{{else}}Keine Anmeldung erforderlich! Schaue einfach beim Treffen vorbei.{{end}}</p>
		{{end}}
		<div class="description" style="margin: 15px 0 15px 0">{{.HtmlDescription}}</div>
		<div class="clearfix"></div>
		<div class="fb-comments" data-href="{{$.hostname}}{{.Url}}" data-width="100%" data-numposts="5" data-order-by="time" data-colorscheme="light"></div>
	</div>{{end}}
</div>
{{if len .similiars}}
<div class="row tiles" style="padding-bottom: 0">
	<div class="col-xs-12"><h4>Weitere gemeinschaftliche Veranstaltungen{{if .event.Addr.City}} in {{.event.Addr.City}}{{end}}</h4></div>
</div>
<div class="row tiles">
	{{range .similiars}}
		<div class="col-md-3 col-sm-4 col-xs-6 col-tile">
			<div class="tile">
				<a href="{{.Url}}?from={{$.from.Unix}}" style="display:block" title="Infos zu {{.Title}} anschauen">
				{{if or (.Image) ((index $.organizers .OrganizerId).Image)}}
					<!-- {{if len .Categories}}{{with index .Categories 0}}<div class="small-icon"><span class="fa fa-{{categoryIcon .}} fa-fw" title="{{categoryTitle .}}"></span></div>{{end}}{{end}} -->
					<div class="tile-image" style="background-image: url(/bild/{{if .Image}}{{.Image}}{{else}}{{(index $.organizers .OrganizerId).Image}}{{end}}?height=165)"> </div>
				{{ end }}
				<div class="tile-text">
					<h3>{{.Title}}</h3>
					<p class="datetime">{{if .Recurrence}}{{.Recurrence}}{{else}}{{longDatetimeFormat (.NextDate $.from)}}{{end}}</p>
					{{if $.organizers}}{{if ne ((index $.organizers .OrganizerId).Name) ("Mitmach-Republik")}}<p class="datetime">{{(index $.organizers .OrganizerId).Name}}</p>{{end}}{{end}}
					<p class="place">{{if .Addr.Name}}{{.Addr.Name}}{{if .Addr.City}}, {{end}}{{end}}{{citypartName .Addr}}</p>
					<p class="description">{{strClip .PlainDescription 150}}</p>
					<p class="highlight" style="position: absolute; bottom: 11px"><span class="fa fa-caret-right"></span> Veranstaltung ansehen</p>
				</div>
				</a>
			</div>
		</div>
	{{end}}
</div>
<div class="row tiles">
	<div class="col-sm-3 col-xs-2">&nbsp;</div>
	<div class="col-sm-6 col-xs-8">
		<a href="/veranstaltungen/{{eventSearchUrl .event.Addr.City .event.Targets .event.Categories (intSlice 0) 0}}" class="btn btn-mmr" style="width: 100%">Weitere Veranstaltungen</a>
	</div>
	<div class="col-sm-3 col-xs-2">&nbsp;</div>
</div>
{{end}}
{{if not .noindex}}
<script type="application/ld+json">
	{
		"@context": "http://schema.org",
		"@type": "Event",
		"location": {
			"@type": "Place",
			"address": {
				"@type": "PostalAddress"
{{if .event.Addr.City}}
				, "addressLocality": {{.event.Addr.City}}
{{end}}
{{if .event.Addr.Pcode}}
				, "postalCode": {{.event.Addr.Pcode}}
{{end}}
{{if .event.Addr.Street}}
				, "streetAddress": {{.event.Addr.Street}}
{{end}}
			}
{{if .event.Addr.Name}}
			, "name": {{.event.Addr.Name}}
{{end}}
		},
		"organizer": {
			"@type": "Organization",
			"location": {
				"@type": "Place",
				"address": {
					"@type": "PostalAddress"
{{if $.organizer.Addr.City}}
					, "addressLocality": {{$.organizer.Addr.City}}
{{end}}
{{if $.organizer.Addr.Pcode}}
					, "postalCode": {{$.organizer.Addr.Pcode}}
{{end}}
{{if $.organizer.Addr.Street}}
					, "streetAddress": {{$.organizer.Addr.Street}}
{{end}}
				}
{{if $.organizer.Addr.Name}}
				, "name": {{$.organizer.Addr.Name}}
{{end}}
			},
			"name": {{$.organizer.Name}},
			"description": {{$.organizer.Descr}},
{{if $.organizer.Image}}
			"image": "/bild/{{$.organizer.Image}}?width=300",
{{end}}
			"url": {{$.organizer.Url}}
  		},
		"name": {{.event.Title}},
{{if .event.Image}}
		"image": "/bild/{{.event.Image}}?width=300",
{{end}}
		"url": {{.event.Url}},
{{if $.showDate}}
		"startDate": {{iso8601Format $.start}},
{{end}}
		"description": {{.event.PlainDescription}}
	}
 </script>
 {{end}}