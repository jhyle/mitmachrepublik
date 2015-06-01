{{with (or .user .organizer)}}
<a href="{{.Url}}" itemprop="url" title="Veranstaltungen des Organisators anzeigen">
	<h3 style="font-weight: normal" itemprop="name">{{.Name}}</h3>
	{{if .Image}}<p><img itemprop="image" src="/bild/{{.Image}}?width=250" title="{{.Name}}" class="img-responsive" /></p>{{end}}
</a>
{{ if not .Addr.IsEmpty }}
	<p class="small-icon pull-left"><span class="fa fa-map-marker fa-fw" title="Ort"></span></p>
	<p itemprop="location" itemscope itemtype="http://schema.org/Place">{{ if .Addr.Name }}<span itemprop="name">{{.Addr.Name}}</span><br />{{ end }}<span class="address" itemprop="address" itemscope itemtype="http://schema.org/PostalAddress">{{ if .Addr.Street }}<span itemprop="streetAddress">{{.Addr.Street}}</span>, {{ end }}{{ if .Addr.Pcode }}<span itemprop="postalCode">{{.Addr.Pcode}}</span> {{ end }}<span itemprop="addressLocality">{{.Addr.City}}</span></span></p>
{{ end }}
<a href="{{.Url}}" title="Veranstaltungen des Organisators anzeigen">
	<p itemprop="description">{{.HtmlDescription}}</p>
</a>
	<p><a href="{{.Url}}" title="Veranstaltungen des Organisators anzeigen" class="highlight"><span class="fa fa-caret-right"></span> Alle Veranstaltungen</a></p>
{{if .Web}}
	<p><a href="{{.Web}}" title="Webseite des Organisators anzeigen" target="_blank" class="highlight"><span class="fa fa-caret-right"></span> {{.Web}}</a></p>
{{end}}
{{ if not .Addr.IsEmpty }}
	<a href="http://maps.google.de/maps?hl=de&q={{.Addr.Name}}+{{.Addr.Street}}+{{.Addr.Pcode}}+{{.Addr.City}}&ie=UTF8" target="_blank" title="In Google Maps öffnen">
		<img src="http://maps.googleapis.com/maps/api/staticmap?center={{.Addr.Name}}+{{.Addr.Street}}+{{.Addr.Pcode}}+{{.Addr.City}}&markers={{.Addr.Street}}+{{.Addr.Pcode}}+{{.Addr.City}}&zoom=15&size=250x187&key={{$.googleApiKey}}" title="in Google Maps öffnen">
	</a>
{{end}}
{{end}}