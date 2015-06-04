{{noescape "<?xml version='1.0' encoding='utf-8'?>"}}
 
<rss version="2.0">
 
  <channel>
    <title>{{.meta.Title}}</title>
    <description>{{.meta.FB_Descr}}</description>
    <language>de-de</language>
    <image>
      <url>{{.meta.FB_Image}}</url>
      <title>Mitmach-Republik</title>
      <link>http://{{.hostname}}/</link>
    </image>
 
 	{{range .items}}
    <item>
      <guid>{{.Id}}</guid>
      <title>{{.Title}}</title>
      <description>{{.Description}}</description>
      <link>http://{{$.hostname}}{{.Link}}</link>
    </item>
    {{end}}
 
  </channel>
</rss>