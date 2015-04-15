      <div class="row footer">
        <div class="col-xs-1">&nbsp;</div>

	<div class="col-xs-2 col-foot">
	  <ul>
	    <li><h5>Rechtliches</h5></li>
	    <li>Kontakt</li>
	    <li><a href="/impressum">Impressum</a></li>
	    <li><a href="/datenschutz">Datenschutz</a></li>
	    <li><a href="/agbs">AGBs</a></li>
	    <li><h5>Folge uns auf</h5></li>
	    <li><span class="fa fa-facebook fa-fw"></span>Facebook</li>
	    <li><span class="fa fa-twitter fa-fw" ></span>Twitter</li>
	    <li><h5>Spenden</h5></li>
	    <li><img src="/images/flattr_icon.png" /> Flattr</li>
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
	<script src="/js/scripts-1.js"></script>
	<script src="https://apis.google.com/js/platform.js" async defer>
  		{lang: 'de'}
	</script>
</body> </html>