<div class="row tiles">
	<div class="col-xs-1">&nbsp;</div>
	<div class="col-xs-11"><h1>{{.meta.FB_Title}}</h1></div>
</div>
<div class="row tiles">
	<div class="col-xs-1">&nbsp;</div>
	<div class="col-xs-3 col-box">
		{{template "organizer_box.tpl" .}}
	</div>
	{{with .event}}<div class="col-xs-7">
		<div class="pull-left" style="margin-right: 10px; margin-bottom: 10px">
			{{if .Image}}
				<a href="{{.Web}}" title="Webseite der Veranstaltung aufrufen" target="_blank">
					<img style="margin-right: 10px; margin-bottom: 15px" src="/bild/{{.Image}}?width=300" alt="Veranstaltung {{.Title}}">
				</a>
			{{end}}
			{{if .ImageCredit}}
				<div class="credits">{{.ImageCredit}}</div>
			{{end}}
			{{ if not .Addr.IsEmpty }}
			<a href="http://maps.google.de/maps?hl=de&q={{.Addr.Name}}+{{.Addr.Street}}+{{.Addr.Pcode}}+{{.Addr.City}}&ie=UTF8" target="_blank" title="In Google Maps öffnen">
				<img style="display: block" src="http://maps.googleapis.com/maps/api/staticmap?center={{.Addr.Name}}+{{.Addr.Street}}+{{.Addr.Pcode}}+{{.Addr.City}}&markers={{.Addr.Street}}+{{.Addr.Pcode}}+{{.Addr.City}}&zoom=15&size=300x225&key={{$.googleApiKey}}" title="in Google Maps öffnen" alt="Karte">
			</a>
			{{end}}
		</div>
		<div style="margin-bottom: 15px">
			<a style="margin-right: 10px" href="https://www.facebook.com/sharer/sharer.php?u=http://{{$.hostname}}{{.Url}}" target="_blank"><img src="/images/facebook_share.png"></a>
			<a style="margin-right: 10px" href="https://plus.google.com/share?url=http://{{$.hostname}}{{.Url}}" target="_blank"><img src="/images/google_share.png"></a>
			<a href="http://twitter.com/intent/tweet?url=http://{{$.hostname}}{{.Url}}" target="_blank"><img src="/images/twitter_share.png"></a>
			<div style="display: inline-block; float: right;"><a id="event-mail" title="Empfehle die Veranstaltung per E-Mail" class="highlight" href="/dialog/sendevent/{{.Id.Hex}}" rel="nofollow" data-toggle="modal" data-target="#share"><span class="fa fa-envelope"></span> Empfehlen</a></div>
		</div>
		<p class="small-icon pull-left"><span class="fa fa-calendar fa-fw" title="Datum"></span></p>
		<p class="icon-text date">{{dateFormat .Start}}</p>
		{{if ne (timeFormat .Start) ("00:00")}}
			<p class="small-icon pull-left"><span class="fa fa-clock-o fa-fw" title="Uhrzeit"></span></p>
			<p class="icon-text date">{{timeFormat .Start}}{{if timeFormat .End}}{{if eq (dateFormat .Start) (dateFormat .End)}} bis {{timeFormat .End}}{{end}}{{end}} Uhr</p>
		{{end}}
		{{if dateFormat .End}}{{if ne (dateFormat .Start) (dateFormat .End)}}
			<p class="small-icon pull-left"><span class="fa fa-calendar fa-fw" title="Enddatum"></span></p>
			<p class="icon-text date">{{dateFormat .End}}</p>
			{{if ne (timeFormat .Start) ("23:59")}}
				<p class="small-icon pull-left"><span class="fa fa-clock-o fa-fw" title="Uhrzeit"></span></p>
				<p class="icon-text date">{{timeFormat .End}} Uhr</p>
			{{end}}
		{{end}}{{end}}
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
		{{if gt (len $.recurrences) 1}}
			<p class="small-icon pull-left"><span class="fa fa-repeat fa-fw" title="Wiederholungen"></span></p>
			<p class="icon-text">{{range $i, $date := $.recurrences}}{{if $i}}, {{end}}<a class="highlight" title="{{$date.Title}} am {{dateFormat $date.Start}} in {{citypartName $date.Addr}}" href="{{$date.Url}}" rel="nofollow">{{cut (dateFormat $date.Start) 1}}</a>{{end}}</p>
		{{end}}
		{{if .Web}}
			<p class="small-icon pull-left"><span class="fa fa-external-link fa-fw" title="Webseite"></span></p>
			<p class="icon-text date"><a href="{{.Web}}" class="highlight" title="Webseite von {{.Title}}" target="_blank">{{strClip .Web 30}}</a></p>
		{{end}}
		<p style="font-weight: bolder">{{if .Rsvp}}Anmeldung erforderlich!{{if .Web}} Melde Dich auf der Webseite der Veranstaltung an.{{end}}{{else}}Keine Anmeldung erforderlich! Schaue einfach beim Treffen vorbei.{{end}}</p>
		<div class="description" style="margin: 15px 0 15px 0">{{.HtmlDescription}}</div>
		<div class="clearfix"></div>
		<div class="fb-comments" data-href="http://{{$.hostname}}{{.Url}}" data-width="100%" data-numposts="5" data-order-by="time" data-colorscheme="light"></div>
	</div>{{end}}
	<div class="col-xs-1">&nbsp;</div>
</div>
{{if not .noindex}}{{range $i, $date := $.recurrences}}
<script type="application/ld+json">
	{
		"@context": "http://schema.org",
		"@type": "Event",
		"location": {
			"@type": "Place",
			"address": {
				"@type": "PostalAddress",
{{if $date.Addr.City}}
				"addressLocality": {{$date.Addr.City}},
{{end}}
{{if $date.Addr.Pcode}}
				"postalCode": {{$date.Addr.Pcode}},
{{end}}
{{if $date.Addr.Street}}
				"streetAddress": {{$date.Addr.Street}}
{{end}}
			},
{{if $date.Addr.Name}}
			"name": {{$date.Addr.Name}}
{{end}}
		},
		"organizer": {
			"@type": "Organization",
			"location": {
				"@type": "Place",
				"address": {
					"@type": "PostalAddress",
{{if $.organizer.Addr.City}}
					"addressLocality": {{$.organizer.Addr.City}},
{{end}}
{{if $.organizer.Addr.Pcode}}
					"postalCode": {{$.organizer.Addr.Pcode}},
{{end}}
{{if $.organizer.Addr.Street}}
					"streetAddress": {{$.organizer.Addr.Street}}
{{end}}
				},
{{if $.organizer.Addr.Name}}
				"name": {{$.organizer.Addr.Name}}
{{end}}
			},
			"name": {{$.organizer.Name}},
			"description": {{$.organizer.Descr}},
{{if $.organizer.Image}}
			"image": "/bild/{{$.organizer.Image}}?width=300",
{{end}}
			"url": {{$.organizer.Url}}
  		},
		"name": {{$date.Title}},
{{if $date.Image}}
		"image": "/bild/{{$date.Image}}?width=300",
{{end}}
		"startDate": {{iso8601Format $date.Start}},
		"description": {{$date.PlainDescription}}
	}
 </script>
 {{end}}{{end}}