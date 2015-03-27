<div class="row hidden-xs">
	<div class="col-sm-12 col-banner">
		<img src="/images/hintergrund.jpg" class="img-responsive" />
		<div class="form-inline text-center" style="position: absolute; top: 33%; width: 100%">
			<input name="place" type="text" class="form-control" placeholder="Berlin" style="width: 50%" data-provide="typeahead" autocomplete="off" autofocus value="{{.place}}" />
			<!--select name="radius" class="form-control">
				<option value="0" {{if eq .radius 0}}selected{{end}}>kein Umkreis</option>
				<option value="2" {{if eq .radius 2}}selected{{end}}>2 km</option>
				<option value="5" {{if eq .radius 5}}selected{{end}}>5 km</option>
				<option value="10" {{if eq .radius 10}}selected{{end}}>10 km</option>
				<option value="25" {{if eq .radius 25}}selected{{end}}>25 km</option>
				<option value="50" {{if eq .radius 50}}selected{{end}}>50 km</option>
			</select-->
			<button name="search" value="events" type="submit" class="btn btn-mmr">{{.eventCnt}} Veranstaltung{{if ne .eventCnt 1}}en{{end}}</button>
			<button name="search" value="organizers" type="submit" class="btn btn-mmr">{{.organizerCnt}} Veranstalter</button>
		</div>
	</div>
</div>