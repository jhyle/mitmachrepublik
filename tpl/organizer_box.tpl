<a href="/veranstalter/{{organizerUrl .}}" itemprop="url">
	<h3 itemprop="name">{{.Addr.Name}}</h3>
	{{if .Image}}<p><img itemprop="image" src="/bild/{{.Image}}?width=250" title="{{.Addr.Name}}" class="img-responsive" /></p>{{end}}
</a>
{{ if not .Addr.IsEmpty }}
	<p class="small-icon pull-left"><span class="fa fa-map-marker fa-fw" title="Ort"></span></p>
	<p itemprop="location" itemscope itemtype="http://schema.org/Place">{{ if .Addr.Name }}<span itemprop="name">{{.Addr.Name}}</span><br />{{ end }}<span class="address" itemprop="address" itemscope itemtype="http://schema.org/PostalAddress">{{ if .Addr.Street }}<span itemprop="streetAddress">{{.Addr.Street}}</span>, {{ end }}{{ if .Addr.Pcode }}<span itemprop="postalCode">{{.Addr.Pcode}}</span> {{ end }}<span itemprop="addressLocality">{{.Addr.City}}</span></span></p>
{{ end }}
<a href="/veranstalter/{{organizerUrl .}}">
	<p itemprop="description">{{.Descr}}</p>
</a>
{{if .Web}}
	<p><a href="{{.Web}}" target="_blank" class="highlight"><span class="fa fa-caret-right"></span> {{.Web}}</a></p>
{{end}}
{{ if not .Addr.IsEmpty }}
	<div><iframe width="250" height="187" src="http://maps.google.de/maps?hl=de&q={{.Addr.Street}}%20{{.Addr.Pcode}}%20{{.Addr.City}}&ie=UTF8&t=&z=14&iwloc=B&output=embed" frameborder="0" scrolling="no" marginheight="0" marginwidth="0"></iframe></div>
{{end}}

