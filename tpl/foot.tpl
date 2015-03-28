      <div class="row footer">
        <div class="col-xs-1">&nbsp;</div>

	<div class="col-xs-2 col-foot">
	  <ul>
	    <li><h5>Rechtliches</h5></li>
	    <li>Kontakt</li>
	    <li><a href="/impressum">Impressum</a></li>
	    <li><a href="/datenschutz">Datenschutz</a></li>
	    <li><a href="/agbs">AGBs</a></li>
	  </ul>
	</div>
	<div class="col-xs-2 col-foot">
	  <ul>
	    <li><h5>Folge uns auf</h5></li>
	    <li><span class="fa fa-facebook fa-fw"></span>Facebook</li>
	    <li><span class="fa fa-twitter fa-fw" ></span>Twitter</li>
	    <li><h5>Spenden</h5></li>
	    <li><img src="/images/flattr_icon.png" /> Flattr</li>
	  </ul>
	</div>
	<div class="col-xs-2 col-foot">
	  <ul>
	    <li><h5>Was ist los in..</h5></li>
	    <li><a href="/veranstaltungen/{{eventSearchUrl "Berlin"}}">Berlin</a></li>
	    {{range $district, $quarters := .districts}}
	    	<li><a href="/veranstaltungen/{{eventSearchUrl $district}}">{{cut $district 1}}</a></li>
	    {{end}}
	  </ul>
	</div>
	<div class="col-xs-2 col-foot">
	  <ul>
	    <li><h5>Veranstalter in..</h5></li>
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
	    	{{if $id}}
	    		<li><a href="/veranstaltungen/{{categorySearchUrl $id "Berlin"}}">{{$category}}</a></li>
	    	{{end}}
	    {{end}}
	  </ul>
	</div>

        <div class="col-xs-1">&nbsp;</div>
      </div>
    </div>
	<script src="/js/jquery.min.js"></script>
	<script src="/js/jquery.ui.widget.js"></script>
	<script src="/js/jquery.iframe-transport.js"></script>
	<script src="/js/jquery.fileupload.js"></script>
	<script src="/js/jquery.cookie-1.4.1.min.js"></script>
	<script src="/js/bootstrap.min.js"></script>
	<script src="/js/bootstrap-datetimepicker.min.js"></script>
	<script src="/js/locales/bootstrap-datetimepicker.de.js"></script>
	<script src="/js/bootstrap3-typeahead.min.js"></script>
	<script src="/js/mmr.js"></script>
</body> </html>