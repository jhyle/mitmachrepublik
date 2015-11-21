<div class="row">
	<div class="col-xs-12 col-banner">
		<h1 class="text-center hidden-xs">Willkommen in der Mitmach-Republik!</h1>
		<h2 class="text-center"><span>Gemeinsam</span> aktiv werden.</h2>
		<h4 class="text-center">{{.meta.FB_Descr}}</h4>
		<form role="form" action="/suche" method="POST" class="form-inline text-center">
			<input name="place" type="text" class="form-control" placeholder="Stadt(-teil) oder Postleitzahl" autocomplete="off">
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
