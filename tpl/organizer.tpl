{{template "head.tpl" .}}
<form id="events-form" role="form" action="/suche" method="POST" itemscope itemtype="http://schema.org/Organization"> 
{{template "banner_search.tpl" .}}
{{template "organizer_main.tpl" .}}
</form>
{{template "foot.tpl" .}}