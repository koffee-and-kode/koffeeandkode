# Local vs Deployed Version Differences

## Key Differences

### 1. **Jekyll Version**
- **Local**: Jekyll 3.9 (due to Ruby 2.6.10 compatibility)
- **GitHub Pages**: Uses Jekyll 3.9 (as of 2024) - ✅ Compatible

### 2. **Base URL**
- **Local**: `baseurl: ""` (empty, served at root)
- **GitHub Pages**: 
  - If using `koffeeandkode.com` (user/organization page): `baseurl: ""` ✅
  - If using project page: `baseurl: "/koffeeandkode"` (would need update)

### 3. **URL Configuration**
- **Local**: `http://localhost:4000`
- **Deployed**: `https://koffeeandkode.com`

### 4. **Asset Paths**
- **Local**: Assets resolve correctly with `baseurl: ""`
- **Deployed**: Should work the same, but verify asset paths use `{{ '/assets/...' | relative_url }}`

### 5. **Plugins**
All plugins used are GitHub Pages compatible:
- ✅ `jekyll-feed`
- ✅ `jekyll-sitemap`
- ✅ `jekyll-seo-tag`
- ✅ `jekyll-redirect-from`
- ✅ `jekyll-relative-links`

### 6. **Build Process**
- **Local**: `bundle exec jekyll serve` (development mode with auto-reload)
- **GitHub Pages**: Automatic build on push (production mode, optimized)

## Potential Issues to Watch For

### 1. **Relative URLs**
All asset references should use Jekyll's `relative_url` filter:
```liquid
{{ '/assets/css/main.css' | relative_url }}
```
✅ Already implemented in layouts

### 2. **Collections**
The `_prototypes` collection is configured correctly:
```yaml
collections:
  prototypes:
    output: true
    permalink: /prototypes/:name/
```
✅ Should work on GitHub Pages

### 3. **Syntax Highlighting**
Using `rouge` highlighter which is GitHub Pages compatible:
```yaml
highlighter: rouge
```
✅ Compatible

### 4. **Markdown Parser**
Using `kramdown` with GFM parser:
```yaml
markdown: kramdown
```
✅ Compatible (GFM parser included in Gemfile)

## Testing Before Deployment

1. **Build locally in production mode**:
   ```bash
   bundle exec jekyll build
   ```
   Check `_site/` folder for any errors

2. **Test asset paths**:
   - Verify CSS loads correctly
   - Check JavaScript files
   - Test images (if any)

3. **Check for absolute URLs**:
   - Search for hardcoded `http://localhost:4000`
   - Ensure all links use `relative_url` or `absolute_url` filters

## After Deployment

1. **Verify the site loads**: `https://koffeeandkode.com`
2. **Check console for errors**: Open browser DevTools
3. **Test navigation**: All links should work
4. **Verify dark mode**: Theme toggle should work
5. **Test code copy buttons**: JavaScript should function
6. **Check RSS feed**: `https://koffeeandkode.com/feed.xml`
7. **Verify sitemap**: `https://koffeeandkode.com/sitemap.xml`

## If Issues Occur

### Common Problems:

1. **404 errors on assets**:
   - Check `baseurl` is correct
   - Verify asset paths use `relative_url`

2. **JavaScript not working**:
   - Check browser console for errors
   - Verify file paths are correct

3. **Styling looks different**:
   - Clear browser cache
   - Check if CSS files are loading

4. **Collections not showing**:
   - Verify collection is in `_config.yml`
   - Check file naming in `_prototypes/`

## GitHub Pages Build Logs

If something doesn't work:
1. Go to repository → Actions tab
2. Check the "pages build and deployment" workflow
3. Look for build errors or warnings

