<div class="row">
	<div class="col-xs-12 col-banner">
		<img src="/images/willkommen.jpg" width="1170" height="617">
		<h1 class="text-center" style="position:absolute; top: 74%; width: 100%">Willkommen in der Mitmach-Republik!</h1>
		<h2 class="text-center" style="position:absolute; top: 17%; width: 100%; background-color: rgba(255, 255, 255, 0.5)"><span style="color: #ff5200; ">Gemeinsam</span> <span style="color: #2f3030">aktiv werden.</span></h2>
		<h4 class="text-center" style="padding-left: 5%; padding-right: 5%; position:absolute; top: 85%; width: 100%">{{.meta.FB_Descr}}</h4>
		<form role="form" action="/suche" method="POST" class="form-inline text-center" style="position: absolute; top: 40%; width: 100%">
			<input name="place" type="text" class="form-control" style="width:20%" placeholder="Stadt(-teil) oder Postleitzahl" autocomplete="off">
			<select name="target" class="form-control">
				<option value="0">alle Zielgruppen</option>
				{{ range .targets }}
					{{ $id := index $.targetMap . }}
					<option value="{{$id}}">{{.}}</option>					
				{{ end }}
			</select>
			<select name="category" class="form-control">
				<option value="0">alle Kategorien</option>
				{{ range .categories }}
					{{ $id := index $.categoryMap . }}
					<option value="{{$id}}">{{.}}</option>					
				{{ end }}
			</select>
			<select name="date" class="form-control">
				{{ range .dates }}
					<option value="{{.}}">{{index $.dateMap .}}</option>					
				{{ end }}
			</select>
			<!--select name="radius" class="form-control">
				<option value="0">kein Umkreis</option>
				<option value="2">2 km</option>
				<option value="5">5 km</option>
				<option value="10">10 km</option>
				<option value="25">25 km</option>
				<option value="50">50 km</option>
			</select-->
			<button name="search" title="Veranstaltungen anzeigen" value="events" type="submit" class="btn btn-mmr">{{.eventCnt}} Veranstaltung{{if ne .eventCnt 1}}en{{end}}</button>
			<!-- button name="search" title="Organisatoren anzeigen" value="organizers" type="submit" class="btn btn-mmr">{{.organizerCnt}} Organisator{{if ne .organizerCnt 1}}en{{end}}</button -->
		</form>
	</div>
</div>
