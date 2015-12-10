      <div class="row footer">
        <div class="col-lg-1 col-md-1 visible-lg-block visible-md-block">&nbsp;</div>

	<div class="col-md-2 col-sm-3 col-xs-6 col-foot">
	  <ul>
	    <li><h5>Folge uns auf</h5></li>
	    <li><a href="https://www.facebook.com/mitmachrepublik" target="_blank"><span class="fa fa-facebook fa-fw"></span>Facebook</a></li>
	    <li><a href="https://twitter.com/mitmachrepublik" target="_blank"><span class="fa fa-twitter fa-fw" ></span>twitter</a></li>
	    <li><h5>Spenden</h5></li>
	    <li><a href="https://flattr.com/submit/auto?user_id=mitmachrepublik&url={{"http://www.mitmachrepublik.de/"}}&title=Mitmach-Republik&description={{.meta.FB_Descr}}&language=de_DE" target="_blank"><img src="/images/flattr_icon.png" alt="Flattr"> Flattr</a></li>
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
	    {{range $district, $quarters := .districts}}
	    	<li><a href="/veranstaltungen/{{simpleEventSearchUrl $district}}" title="Veranstaltungen in {{cut $district 1}}">{{cut $district 1}}</a></li>
	    {{end}}
	  </ul>
	</div>
	<div class="col-md-2 col-sm-3 col-xs-6  col-foot">
	  <ul>
	    <li><h5><a href="/veranstalter/{{organizerSearchUrl ""}}" title="Alle Veranstalter">Organisatoren</a> in..</h5></li>
	    <li><a href="/veranstalter/{{organizerSearchUrl "Berlin"}}" title="Veranstalter in Berlin">Berlin</a></li>
	    {{range $district, $quarters := .districts}}
	    	<li><a href="/veranstalter/{{organizerSearchUrl $district}}" title="Veranstalter in {{cut $district 1}}">{{cut $district 1}}</a></li>
	    {{end}}
	  </ul>
	</div>
	<div class="col-md-2 col-sm-3 col-xs-6 col-foot">
	  <ul style="margin-bottom: 0; padding-bottom: 0">
	    <li><h5>In Kategorie..</h5></li>
	    {{range $i, $category := .categories}}
    		<li><a href="/veranstaltungen/{{categorySearchUrl (index $.categoryMap $category) "Berlin"}}" title="{{$category}} in Berlin">{{$category}}</a></li>
	    {{end}}
	  </ul>
	</div>
	<div class="col-md-2 col-foot visible-lg-block visible-md-block">
	  <ul style="margin-bottom: 0; padding-bottom: 0">
	    <li><h5>Für..</h5></li>
	    {{range $i, $target := .targets}}
    		<li><a href="/veranstaltungen/{{targetSearchUrl (index $.targetMap $target) "Berlin"}}" title="Veranstaltungen für {{$target}} in Berlin">{{$target}}</a></li>
	    {{end}}
	    <li><h5><a href="/veranstaltungen/{{simpleEventSearchUrl "Berlin"}}">Berlin</a></h5></li>
		<li><a href="/veranstaltungen/{{eventSearchUrl "Berlin" (intSlice 0) (intSlice 0) (intSlice 1) 0}}" title="Veranstaltungen in Berlin heute">Heute in Berlin</a></li> 
		<li><a href="/veranstaltungen/{{eventSearchUrl "Berlin" (intSlice 0) (intSlice 0) (intSlice 2) 0}}" title="Veranstaltungen in Berlin morgen">Morgen in Berlin</a></li>
		<li><a href="/veranstaltungen/{{eventSearchUrl "Berlin" (intSlice 0) (intSlice 0) (intSlice 7) 0}}" title="Veranstaltungen in Berlin übermorgen">Übermorgen in Berlin</a></li>
		<li><a href="/veranstaltungen/{{eventSearchUrl "Berlin" (intSlice 0) (intSlice 0) (intSlice 4) 0}}" title="Veranstaltungen in Berlin am nächsten Wochenende">Berlin am nächsten Wochenende</a></li> 
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
	<script src="/js/scripts-21.js"></script>
</body> </html>
