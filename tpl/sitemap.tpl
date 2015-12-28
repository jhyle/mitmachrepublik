{{noescape "<?xml version='1.0' encoding='utf-8'?>"}}
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>http://{{$.hostname}}/</loc>
    <changefreq>hourly</changefreq>
    <priority>1.0</priority>
  </url>
  <url>
    <loc>http://{{$.hostname}}/veranstaltungen/{{eventSearchUrl "Berlin" (intSlice 0) (intSlice 0) (intSlice 1) 0}}</loc>
    <changefreq>hourly</changefreq>
    <priority>0.9</priority>
  </url>
  <url>
    <loc>http://{{$.hostname}}/veranstaltungen/{{eventSearchUrl "Berlin" (intSlice 0) (intSlice 0) (intSlice 2) 0}}</loc>
    <changefreq>daily</changefreq>
    <priority>0.9</priority>
  </url>
  <url>
    <loc>http://{{$.hostname}}/veranstaltungen/{{eventSearchUrl "Berlin" (intSlice 0) (intSlice 0) (intSlice 4) 0}}</loc>
    <changefreq>daily</changefreq>
    <priority>0.9</priority>
  </url>
{{range .events}}
  <url>
    <loc>http://{{$.hostname}}{{.Url}}</loc>
    <priority>0.7</priority>
    <changefreq>daily</changefreq>
  </url>
{{end}}
{{range .organizers}}
  <url>
    <loc>http://{{$.hostname}}{{.Url}}</loc>
    <priority>0.8</priority>
    <changefreq>weekly</changefreq>
  </url>
{{end}}
</urlset>