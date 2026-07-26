source "https://rubygems.org"

# For local development, use Jekyll 3.9 (compatible with Ruby 2.6)
# GitHub Pages will use its own Jekyll version automatically
gem "jekyll", "~> 3.9"
gem "kramdown-parser-gfm"

# GitHub Pages compatible plugins
gem "jekyll-feed", "~> 0.12"
gem "jekyll-sitemap", "~> 1.4"
gem "jekyll-seo-tag", "~> 2.8"

# Windows and JRuby does not include zoneinfo files, so bundle the tzinfo-data gem
# and associated library.
platforms :mingw, :x64_mingw, :mswin, :jruby do
  gem "tzinfo", ">= 1", "< 3"
  gem "tzinfo-data"
end

# Performance and no-op maintenance plugins
gem "jekyll-redirect-from", "~> 0.16"
gem "jekyll-relative-links", "~> 0.6"

# Lock `http_parser.rb` gem to `v0.6.x` on JRuby builds since newer versions of the gem
# do not have a Java counterpart.
gem "http_parser.rb", "~> 0.6.0", :platforms => [:jruby]

# Uncomment the line below to use GitHub Pages locally (optional)
# gem "github-pages", group: :jekyll_plugins
