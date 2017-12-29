{{noescape "<?xml version='1.0' encoding='utf-8'?>"}}
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>{{$.hostname}}/</loc>
    <changefreq>hourly</changefreq>
    <priority>1.0</priority>
  </url>
{{range .topics}}
  <url>
    <loc>{{$.hostname}}/{{.}}</loc>
    <changefreq>hourly</changefreq>
    <priority>0.9</priority>
  </url>
{{end}}
{{range .events}}
  <url>
    <loc>{{$.hostname}}{{.Url}}</loc>
    <priority>0.7</priority>
    <changefreq>daily</changefreq>
  </url>
{{end}}
{{range .organizers}}
  <url>
    <loc>{{$.hostname}}{{.Url}}</loc>
    <priority>0.8</priority>
    <changefreq>weekly</changefreq>
  </url>
{{end}}
</urlset>