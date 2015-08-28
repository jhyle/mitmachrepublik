{{noescape "<?xml version='1.0' encoding='utf-8'?>"}}

<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"> 
  <url>
    <loc>http://{{$.hostname}}/</loc>
    <changefreq>hourly</changefreq> 
  </url>
{{range .dates}}
  <url>
    <loc>http://{{$.hostname}}{{.Url}}</loc>
  </url>
{{end}}
</urlset>