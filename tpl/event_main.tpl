<div class="row tiles">
	<div class="col-xs-1">&nbsp;</div>
	<div class="col-xs-11"><h1 itemprop="name">{{.meta.FB_Title}}</h1></div>
</div>
<div class="row tiles">
	<div class="col-xs-1">&nbsp;</div>
	<div class="col-xs-3 col-box" itemprop="organizer" itemscope itemtype="http://schema.org/Organization">
		{{template "organizer_box.tpl" .}}
	</div>
	{{with .event}}<div class="col-xs-7">
		<div class="pull-left" style="margin-right: 5px; margin-bottom: 10px">
			{{if .Image}}
				<a href="{{.Web}}" title="Webseite der Veranstaltung anzeigen" target="_blank"><img itemprop="image" style="margin-right: 10px; margin-bottom: 15px" src="/bild/{{.Image}}?width=300" title="{{.Title}}"></a>
			{{end}}
			{{ if not .Addr.IsEmpty }}
			<a href="http://maps.google.de/maps?hl=de&q={{.Addr.Name}}+{{.Addr.Street}}+{{.Addr.Pcode}}+{{.Addr.City}}&ie=UTF8" target="_blank" title="In Google Maps öffnen">
				<img style="display: block" src="http://maps.googleapis.com/maps/api/staticmap?center={{.Addr.Name}}+{{.Addr.Street}}+{{.Addr.Pcode}}+{{.Addr.City}}&markers={{.Addr.Street}}+{{.Addr.Pcode}}+{{.Addr.City}}&zoom=15&size=300x225&key={{$.googleApiKey}}" title="in Google Maps öffnen">
			</a>
			{{end}}
		</div>
		<div style="height: 30px; margin-bottom: 10px">
			<div class="g-plus" style="float: left, padding-right: 10px" data-action="share" data-annotation="none" data-href="http://{{$.hostname}}{{.Url}}"></div>
			<a class="twitter-share-button" data-count="none" href="https://twitter.com/share" target="_blank">Tweet</a><script>window.twttr=(function(d,s,id){var js,fjs=d.getElementsByTagName(s)[0],t=window.twttr||{};if(d.getElementById(id))return t;js=d.createElement(s);js.id=id;js.src="https://platform.twitter.com/widgets.js";fjs.parentNode.insertBefore(js,fjs);t._e=[];t.ready=function(f){t._e.push(f);};return t;}(document,"script","twitter-wjs"));</script>
			<div class="fb-share-button" style="float: left; padding-right: 10px" data-href="http://{{$.hostname}}{{.Url}}" data-layout="button"></div>
			<div style="display: inline-block; float: right; line-height: 1"><a id="event-mail" title="Empfehle die Veranstaltung per E-Mail" class="highlight" href="/dialog/sendevent/{{.Id.Hex}}" rel="nofollow" data-toggle="modal" data-target="#share"><span class="fa fa-envelope"></span> E-Mail</a></div>
		</div>
		<p class="small-icon pull-left"><span class="fa fa-calendar fa-fw" title="Datum"></span></p>
		<p class="date" itemprop="startDate" content="{{iso8601Format .Start}}">{{dateFormat .Start}}</p>
		<p class="small-icon pull-left"><span class="fa fa-clock-o fa-fw" title="Uhrzeit"></span></p>
		<p class="date">{{timeFormat .Start}}{{if timeFormat .End}} bis{{if eq (dateFormat .Start) (dateFormat .End)}} {{timeFormat .End}}{{end}}{{end}} Uhr</p>
		{{if dateFormat .End}}{{if ne (dateFormat .Start) (dateFormat .End)}}
			<p class="small-icon pull-left"><span class="fa fa-calendar fa-fw" title="Enddatum"></span></p>
			<p class="date">{{dateFormat .End}}</p>
			<p class="small-icon pull-left"><span class="fa fa-clock-o fa-fw" title="Uhrzeit"></span></p>
			<p class="date">{{timeFormat .End}} Uhr</p>
		{{end}}{{end}}
		{{if .Rsvp}}
			<p style="margin-bottom: 18px"><a style="padding-left: 44px" href="{{.Web}}" class="highlight" target="_blank"><span class="fa fa-caret-right"></span> Anmeldung erforderlich</a></p>
		{{end}}
		{{ if not .Addr.IsEmpty }}
			<p class="small-icon pull-left"><span class="fa fa-map-marker fa-fw" title="Ort"></span></p>
			<p itemprop="location" itemscope itemtype="http://schema.org/Place">{{ if .Addr.Name }}<span itemprop="name">{{.Addr.Name}}</span><br />{{ end }}<span class="address" itemprop="address" itemscope itemtype="http://schema.org/PostalAddress">{{ if .Addr.Street }}<span itemprop="streetAddress">{{.Addr.Street}}</span>, {{ end }}{{ if .Addr.Pcode }}<span itemprop="postalCode">{{.Addr.Pcode}}</span> {{ end }}<span itemprop="addressLocality">{{.Addr.City}}</span></span></p>
		{{ end }}
		{{if len .Categories}}{{with index .Categories 0}}
			<p class="small-icon pull-left"><span class="fa fa-{{categoryIcon .}} fa-fw" title="{{categoryTitle .}}"></span></p>{{end}}
			<p>{{range $i, $category := .Categories}}{{if $i}}, {{end}}{{categoryTitle $category}}{{end}}</p>
		{{end}}
		<div class="description" style="margin: 25px 0 15px 0" itemprop="description">{{.HtmlDescription}}</div>
		<div class="clearfix"></div>
		{{if .Web}}
			<p><a href="{{.Web}}" title="Webseite der Veranstaltung anzeigen" class="btn btn-mmr" style="margin: 0" target="_blank">Zur Veranstaltungs-Webseite</a></p>
		{{end}}
		<div class="fb-comments" data-href="http://{{$.hostname}}{{.Url}}" data-width="100%" data-numposts="5" data-order-by="time" data-colorscheme="light"></div>
	</div>{{end}}
	<div class="col-xs-1">&nbsp;</div>
</div>
