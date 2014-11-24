<h3>{{.Addr.Name}}</h3>
{{if .Image}}<p><img src="/bild/{{.Image}}?width=250" title="{{.Addr.Name}}" class="img-responsive" /></p>{{end}}
<p>{{.Descr}}</p>
{{if .Web}}<p><a href="{{.Web}}" target="_blank" class="highlight"><span class="fa fa-caret-right"></span> {{.Web}}</a></p>{{end}}
