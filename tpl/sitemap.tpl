{{noescape "<?xml version='1.0' encoding='utf-8'?>"}}

<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"> 
  <url>
    <loc>http://{{$.hostname}}/</loc>
    <changefreq>hourly</changefreq> 
  </url>
  <url>
    <loc>http://{{$.hostname}}/veranstaltungen/{{eventSearchUrl "Berlin" (intSlice 0) (intSlice 0) (intSlice 1) 0}}</loc>
    <changefreq>hourly</changefreq> 
  </url>
  <url>
    <loc>http://{{$.hostname}}/veranstaltungen/{{eventSearchUrl "Berlin" (intSlice 0) (intSlice 0) (intSlice 2) 0}}</loc>
    <changefreq>daily</changefreq> 
  </url>
  <url>
    <loc>http://{{$.hostname}}/veranstaltungen/{{eventSearchUrl "Berlin" (intSlice 0) (intSlice 0) (intSlice 4) 0}}</loc>
    <changefreq>daily</changefreq> 
  </url>
{{range .dates}}
  <url>
    <loc>http://{{$.hostname}}{{.Url}}</loc>
  </url>
{{end}}
{{range .organizers}}
  <url>
    <loc>http://{{$.hostname}}{{.Url}}</loc>
  </url>
{{end}}
</urlset>