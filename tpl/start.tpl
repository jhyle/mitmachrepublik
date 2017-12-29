{{template "head.tpl" .}}
{{template "banner.tpl" .}}
{{template "tiles.tpl" .}}
 <script type="application/ld+json">
{
	"@context": "http://schema.org",
	"@type": "Organization",
	"name" : "Mitmach-Republik",
	"url": "{{.hostname}}",
	"logo": "{{.hostname}}/images/mitmachrepublik.png",
	"sameAs" : [
		"https://www.facebook.com/mitmach-republik",
		"https://twitter.com/mitmachrepublik"
	]
}
</script>
{{template "foot.tpl" .}}