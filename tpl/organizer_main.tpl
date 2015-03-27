<div class="row tiles">
	<div class="col-sm-1 hidden-xs">&nbsp;</div>
	<div class="col-sm-10"><h1>Veranstaltungen von {{.organizer.Addr.Name}}</h1></div>
	<div class="col-sm-1 hidden-xs">&nbsp;</div>
</div>
<div class="row tiles">
	<div class="col-sm-1 hidden-xs">&nbsp;</div>
	<div class="col-sm-3 col-box">
		{{template "organizer_box.tpl" .organizer}}
	</div>
	<div class="col-sm-7">
		{{template "events_list.tpl" .}}
	</div>
	<div class="col-sm-1 hidden-xs">&nbsp;</div>
</div>
