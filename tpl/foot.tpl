      <div class="row footer">
        <div class="col-lg-1 col-md-1 visible-lg-block visible-md-block">&nbsp;</div>

	<div class="col-md-2 col-sm-3 col-xs-6 col-foot">
	  <ul>
	    <li><h5>Folge uns auf</h5></li>
	    <li><a href="https://www.facebook.com/mitmachrepublik" target="_blank"><span class="fa fa-facebook fa-fw"></span>Facebook</a></li>
	    <li><a href="https://twitter.com/mitmachrepublik" target="_blank"><span class="fa fa-twitter fa-fw" ></span>twitter</a></li>
	    <li><h5>Rechtliches</h5></li>
	    <li><a rel="nofollow" href="/agbs">Allgemeine Geschäftsbedingungen</a></li>
	    <li><a rel="nofollow" href="/datenschutz">Datenschutz</a></li>
	    <li><a rel="nofollow" href="/impressum">Impressum</a></li>
	    <li><a rel="nofollow" href="/disclaimer">Haftungsausschluss (Disclaimer)</a></li>
	  </ul>
	</div>
	<div class="col-md-2 col-sm-3 col-xs-6  col-foot">
	  <ul>
	    <li><h5><a href="/veranstaltungen/{{simpleEventSearchUrl ""}}" title="Alle Veranstaltungen">Veranstaltungen</a> in..</h5></li>
	    <li><a href="/veranstaltungen/{{simpleEventSearchUrl "Berlin"}}" title="Veranstaltungen in Berlin">Berlin</a></li>
	    <li><a href="/veranstaltungen/{{simpleEventSearchUrl "Hamburg"}}" title="Veranstaltungen in Hamburg">Hamburg</a></li>
	    <li><a href="/veranstaltungen/{{simpleEventSearchUrl "München"}}" title="Veranstaltungen in München">München</a></li>
	    <li><a href="/veranstaltungen/{{simpleEventSearchUrl "Köln"}}" title="Veranstaltungen in Köln">Köln</a></li>
	    <li><a href="/veranstaltungen/{{simpleEventSearchUrl "Frankfurt"}}" title="Veranstaltungen in Frankfurt">Frankfurt</a></li>
	  </ul>
	  <ul>
	    <li><h5><a href="/veranstalter/{{organizerSearchUrl ""}}" title="Alle Veranstalter">Organisatoren</a> in..</h5></li>
	    <li><a href="/veranstalter/{{organizerSearchUrl "Berlin"}}" title="Veranstalter in Berlin">Berlin</a></li>
	    <li><a href="/veranstalter/{{organizerSearchUrl "Hamburg"}}" title="Veranstalter in Hamburg">Hamburg</a></li>
	    <li><a href="/veranstalter/{{organizerSearchUrl "München"}}" title="Veranstalter in München">München</a></li>
	    <li><a href="/veranstalter/{{organizerSearchUrl "Köln"}}" title="Veranstalter in Köln">Köln</a></li>
	    <li><a href="/veranstalter/{{organizerSearchUrl "Frankfurt"}}" title="Veranstalter in Frankfurt">Frankfurt</a></li>
	  </ul>
	</div>
	<div class="col-md-2 col-sm-3 col-xs-6  col-foot">
	  <ul style="margin-bottom: 0; padding-bottom: 0">
	    <li><h5>In Kategorie..</h5></li>
	    {{range $i, $category := .categories}}
    		<li><a href="/veranstaltungen/{{categorySearchUrl (index $.categoryMap $category) "Berlin"}}" title="{{$category}} in Berlin">{{$category}}</a></li>
	    {{end}}
	  </ul>
	</div>
	<div class="col-md-2 col-sm-3 col-xs-6 col-foot">
	  <ul style="margin-bottom: 0; padding-bottom: 0">
	    <li><h5>Für..</h5></li>
	    {{range $i, $target := .targets}}
    		<li><a href="/veranstaltungen/{{targetSearchUrl (index $.targetMap $target) "Berlin"}}" title="Veranstaltungen für {{$target}} in Berlin">{{$target}}</a></li>
	    {{end}}
	  </ul>
	</div>
	<div class="col-md-2 col-foot visible-lg-block visible-md-block">
	  <ul style="margin-bottom: 0; padding-bottom: 0">
	    <li><h5><a href="/veranstaltungen/{{simpleEventSearchUrl "Berlin"}}">Berlin</a></h5></li>
		<li><a href="/heute-in-berlin">Heute in Berlin</a></li>
		<li><a href="/morgen-in-berlin">Morgen in Berlin</a></li>
		<li><a href="/uebermorgen-in-berlin">Übermorgen in Berlin</a></li>
		<li><a href="/am-wochenende-in-berlin">Am nächsten Wochenende</a></li>
		<li>&nbsp;</li> 
		<li><a href="/babies-und-kleinkinder">Babies &amp; Kleinkinder</a></li> 
		<li><a href="/sport-und-gesundheit">Sport &amp; Gesundheit</a></li> 
		<li><a href="/natur-und-garten">Natur &amp; Garten</a></li> 
		<li><a href="/eltern-und-familien">Eltern &amp; Familien</a></li> 
		<li><a href="/bildung-und-kultur">Bildung &amp; Kultur</a></li> 
		<li><a href="/umwelt-und-tierschutz">Umwelt- &amp; Tierschutz</a></li> 
		<li><a href="/demonstrationen-und-politik">Demos &amp; Politik</a></li> 
		<li><a href="/soziales-und-ehrenamt">Soziales &amp; Ehrenamt</a></li> 
	  </ul>
	</div>

        <div class="col-lg-1 col-md-1 visible-lg-block visible-md-block">&nbsp;</div>
      </div>
    </div>
	<!-- HTML5 Shim and Respond.js IE8 support of HTML5 elements and media queries -->
	<!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
	{{noescape "<!--[if lt IE 9]>"}}
		<script src="https://oss.maxcdn.com/html5shiv/3.7.2/html5shiv.min.js"></script>
		<script src="https://oss.maxcdn.com/respond/1.4.2/respond.min.js"></script>
	{{noescape "<![endif]-->"}}
	<script src="/js/scripts-23.js"></script>
</body> </html>
