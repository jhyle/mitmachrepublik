{{template "head.tpl"}}
<form id="events-form" role="form" action="/suche" method="POST"> 

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
<!-- div class="row tiles">
	{{ $len := len .events }}
	<div class="col-sm-1 hidden-xs">&nbsp;</div>
	<div class="col-sm-10"><h1>{{if eq $len 1}}Eine Veranstaltung{{else}}{{$len}} Veranstaltungen{{end}} in Berlin</h1></div>
	<div class="col-sm-1 hidden-xs">&nbsp;</div>
</div-->
<div class="row tiles">
	<div class="col-sm-1 hidden-xs">&nbsp;</div>
	<div class="col-sm-3 col-box">
		<h3>Filter</h3>
		<h5>Kategorien</h5>
		<hr>
		{{ range .categories }}
			{{ $id := index $.categoryMap . }}
			<label class="checkbox"><input type="checkbox" name="category" value="{{$id}}"
			{{ range $.categoryIds }}
				{{ if eq $id . }} checked {{ end }}
			{{ end }}
			>  {{.}}</label>					
		{{ end }}
		<h5>Datum</h5>
		<hr>
		<label class="checkbox"><input type="checkbox" name="date" value="1"
		{{ range $.dates }}
			{{ if eq "heute" . }} checked {{ end }}
		{{ end }}
		>  Heute</label>
		<label class="checkbox"><input type="checkbox" name="date" value="2"
		{{ range $.dates }}
			{{ if eq "morgen" . }} checked {{ end }}
		{{ end }}
		>  Morgen</label>
		<label class="checkbox"><input type="checkbox" name="date" value="3"
		{{ range $.dates }}
			{{ if eq "wochenende" . }} checked {{ end }}
		{{ end }}
		>  Am Wochenende</label>
	</div>
	<div class="col-sm-7">
		{{template "events_list.tpl" .}}
	</div>
</div>

</form>
{{template "foot.tpl"}}