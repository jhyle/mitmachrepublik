      <div class="row footer">
        <div class="col-xs-1">&nbsp;</div>

	<div class="col-xs-2 col-foot">
	  <ul>
	    <li><h5>Folge uns auf</h5></li>
	    <li><a href="https://www.facebook.com/mitmachrepublik" target="_blank"><span class="fa fa-facebook fa-fw"></span>Facebook</a></li>
	    <li><a href="https://twitter.com/mitmachrepublik" target="_blank"><span class="fa fa-twitter fa-fw" ></span>twitter</a></li>
	    <li><h5>Spenden</h5></li>
	    <li><a href="https://flattr.com/submit/auto?user_id=mitmachrepublik&url={{"http://www.mitmachrepublik.de/"}}&title=Mitmach-Republik&description={{.meta.FB_Descr}}&language=de_DE" target="_blank"><img src="/images/flattr_icon.png"> Flattr</a></li>
	    <li><h5>Rechtliches</h5></li>
	    <li><a rel="nofollow" href="/agbs">Allgemeine Geschäftsbedingungen</a></li>
	    <li><a rel="nofollow" href="/datenschutz">Datenschutz</a></li>
	    <li><a rel="nofollow" href="/impressum">Impressum</a></li>
	    <li><a rel="nofollow" href="/disclaimer">Haftungsausschluss (Disclaimer)</a></li>
	  </ul>
	</div>
	<div class="col-xs-2 col-foot">
	  <ul>
	    <li><h5>Veranstaltungen in..</h5></li>
	    <li><a href="/veranstaltungen/{{eventSearchUrl "Berlin"}}">Berlin</a></li>
	    {{range $district, $quarters := .districts}}
	    	<li><a href="/veranstaltungen/{{eventSearchUrl $district}}">{{cut $district 1}}</a></li>
	    {{end}}
	  </ul>
	</div>
	<div class="col-xs-2 col-foot">
	  <ul>
	    <li><h5>Organisatoren in..</h5></li>
	    <li><a href="/veranstalter/{{organizerSearchUrl "Berlin"}}">Berlin</a></li>
	    {{range $district, $quarters := .districts}}
	    	<li><a href="/veranstalter/{{organizerSearchUrl $district}}">{{cut $district 1}}</a></li>
	    {{end}}
	  </ul>
	</div>
	<div class="col-xs-2 col-foot">
	  <ul style="margin-bottom: 0; padding-bottom: 0">
	    <li><h5>In Kategorie..</h5></li>
	    {{range $category, $id := .categoryMap}}
	    	{{if $id}}{{if le $id 9}}
	    		<li><a href="/veranstaltungen/{{categorySearchUrl $id "Berlin"}}">{{$category}}</a></li>
	    	{{end}}{{end}}
	    {{end}}
	  </ul>
	</div>
	<div class="col-xs-2 col-foot">
	  <ul style="margin-bottom: 0; padding-bottom: 0">
	    <li><h5>&nbsp;</h5></li>
	    {{range $category, $id := .categoryMap}}
	    	{{if $id}}{{if gt $id 9}}
	    		<li><a href="/veranstaltungen/{{categorySearchUrl $id "Berlin"}}">{{$category}}</a></li>
	    	{{end}}{{end}}
	    {{end}}
	  </ul>
	</div>

        <div class="col-xs-1">&nbsp;</div>
      </div>
    </div>
	<script src="/js/scripts-2.js"></script>
	<script src="https://apis.google.com/js/platform.js" async defer>
  		{lang: 'de'}
	</script>
</body> </html>