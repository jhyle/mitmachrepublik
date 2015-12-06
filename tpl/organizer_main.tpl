<div class="row tiles">
	<div class="col-xs-12"><h1>Gemeinschaftliche Veranstaltungen von {{.organizer.Name}}</h1></div>
</div>
<div class="row tiles">
	<div class="col-lg-3 col-sm-4 hidden-xs col-box">
		{{template "organizer_box.tpl" .}}
	</div>
	<div class="col-lg-9 col-sm-8 col-xs-12">
		{{template "events_list.tpl" .}}
	</div>
</div>
