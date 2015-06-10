<div class="row">
	<div class="col-xs-12 col-banner">
		<img src="/images/hintergrund.jpg" class="img-responsive" width="1170" height="167" alt="Finde Veranstaltungen zum Mitmachen!"/>
		<div class="form-inline text-center" style="position: absolute; top: 33%; width: 100%">
			<input name="place" type="text" class="form-control" value="{{.place}}" placeholder="Wo?" style="width: 50%" autocomplete="off" autofocus />
			<!--select name="radius" class="form-control">
				<option value="0" {{if eq .radius 0}}selected{{end}}>kein Umkreis</option>
				<option value="2" {{if eq .radius 2}}selected{{end}}>2 km</option>
				<option value="5" {{if eq .radius 5}}selected{{end}}>5 km</option>
				<option value="10" {{if eq .radius 10}}selected{{end}}>10 km</option>
				<option value="25" {{if eq .radius 25}}selected{{end}}>25 km</option>
				<option value="50" {{if eq .radius 50}}selected{{end}}>50 km</option>
			</select-->
			<button name="search" value="events" type="submit" class="btn btn-mmr">{{.eventCnt}} Veranstaltung{{if ne .eventCnt 1}}en{{end}}</button>
			<!-- button name="search" value="organizers" type="submit" class="btn btn-mmr">{{.organizerCnt}} Organisator{{if ne .organizerCnt 1}}en{{end}}</button -->
		</div>
	</div>
</div>
