{{template "head.tpl" .}}
<form id="events-form" role="form" action="/suche" method="POST" itemscope itemtype="http://schema.org/Event"> 
{{template "banner_search.tpl" .}}
{{template "event_main.tpl" .}}
</form>
{{template "foot.tpl" .}}