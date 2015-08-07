      <div class="row footer">
        <div class="col-xs-1">&nbsp;</div>

	<div class="col-xs-2 col-foot">
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
	<div class="col-xs-2 col-foot">
	  <ul>
	    <li><h5><a href="/veranstaltungen/{{simpleEventSearchUrl ""}}">Veranstaltungen</a> in..</h5></li>
	    <li><a href="/veranstaltungen/{{simpleEventSearchUrl "Berlin"}}">Berlin</a></li>
	    {{range $district, $quarters := .districts}}
	    	<li><a href="/veranstaltungen/{{simpleEventSearchUrl $district}}">{{cut $district 1}}</a></li>
	    {{end}}
	  </ul>
	</div>
	<div class="col-xs-2 col-foot">
	  <ul>
	    <li><h5><a href="/veranstalter/{{organizerSearchUrl ""}}">Organisatoren</a> in..</h5></li>
	    <li><a href="/veranstalter/{{organizerSearchUrl "Berlin"}}">Berlin</a></li>
	    {{range $district, $quarters := .districts}}
	    	<li><a href="/veranstalter/{{organizerSearchUrl $district}}">{{cut $district 1}}</a></li>
	    {{end}}
	  </ul>
	</div>
	<div class="col-xs-2 col-foot">
	  <ul style="margin-bottom: 0; padding-bottom: 0">
	    <li><h5>In Kategorie..</h5></li>
	    {{range $i, $category := .categories}}
    		<li><a href="/veranstaltungen/{{categorySearchUrl (index $.categoryMap $category) "Berlin"}}">{{$category}}</a></li>
	    {{end}}
	  </ul>
	</div>
	<div class="col-xs-2 col-foot">
	  <ul style="margin-bottom: 0; padding-bottom: 0">
	    <li><h5>Für..</h5></li>
	    {{range $i, $target := .targets}}
    		<li><a href="/veranstaltungen/{{targetSearchUrl (index $.targetMap $target) "Berlin"}}">{{$target}}</a></li>
	    {{end}}
	  </ul>
	</div>

        <div class="col-xs-1">&nbsp;</div>
      </div>
    </div>
	<!-- HTML5 Shim and Respond.js IE8 support of HTML5 elements and media queries -->
	<!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
	{{noescape "<!--[if lt IE 9]>"}}
		<script src="https://oss.maxcdn.com/html5shiv/3.7.2/html5shiv.min.js"></script>
		<script src="https://oss.maxcdn.com/respond/1.4.2/respond.min.js"></script>
	{{noescape "<![endif]-->"}}
	<script src="/js/scripts-15.js"></script>
</body> </html>
