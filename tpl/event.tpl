{{template "head.tpl"}}
<form id="events-form" role="form" action="/suche" method="POST"> 
{{template "banner_search.tpl" .}}
{{template "event_main.tpl" .}}
</form>
{{template "foot.tpl"}}