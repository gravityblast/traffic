{{ template "includes/header" }}
  <div class="sidebar">
    <ul>
      <li><a href="/">Home</a></li>
      <li>
        <a href="/about/about-us">/about/about-us</a>
      </li>
      <li>
        <a href="/foo/bar">/foo/bar</a>
      </li>
      <li>
        <a href="/1/2/3/4/5">/1/2/3/4/5</a>
      </li>
    </ul>
  </div>
  <div class="content">
    <p>Current page path: {{ .PagePath }}</p>
  </div>
{{ template "includes/footer" }}

