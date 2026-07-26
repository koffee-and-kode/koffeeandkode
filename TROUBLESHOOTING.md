# Troubleshooting: Deployed Site Looks Bad

## Common Issues and Fixes

### 1. CSS/JavaScript Not Loading

**Symptoms**: Site appears unstyled, no colors, broken layout

**Possible Causes**:
- Asset paths are incorrect
- Files not committed to repository
- GitHub Pages build failed
- Browser cache issues

**Solutions**:

1. **Verify files are committed**:
   ```bash
   git status
   git add assets/
   git commit -m "Add assets"
   git push
   ```

2. **Check GitHub Pages build**:
   - Go to repository → Actions tab
   - Check if build succeeded
   - Look for any errors

3. **Clear browser cache**:
   - Hard refresh: `Cmd+Shift+R` (Mac) or `Ctrl+Shift+R` (Windows)
   - Or open in incognito/private mode

4. **Verify asset paths in HTML**:
   - View page source on deployed site
   - Check if CSS/JS links are correct
   - Should be: `/assets/css/main.css` or `https://koffeeandkode.com/assets/css/main.css`

### 2. Layout Broken

**Symptoms**: Content not centered, navigation broken, footer misaligned

**Check**:
- Container divs are present in layouts
- CSS is loading (see issue #1)
- No JavaScript errors in console

### 3. Dark Mode Not Working

**Symptoms**: Theme toggle doesn't work

**Check**:
- JavaScript file is loading
- Browser console for errors
- LocalStorage is enabled

### 4. Code Blocks Not Highlighted

**Symptoms**: Code blocks appear plain text

**Check**:
- Syntax CSS is loading
- Rouge highlighter is configured in `_config.yml`
- Code blocks use proper markdown syntax

## Quick Fixes

### Force Rebuild on GitHub Pages

1. Make a small change (add a space to `_config.yml`)
2. Commit and push
3. Wait for GitHub Pages to rebuild (check Actions tab)

### Verify Asset Paths

The current setup uses:
- `{{ '/assets/css/main.css' | relative_url }}` for local development
- Should work on GitHub Pages with `baseurl: ""`

If issues persist, try:
```liquid
{{ site.baseurl }}/assets/css/main.css
```

### Check Build Output

1. Build locally: `bundle exec jekyll build`
2. Check `_site/assets/` folder exists
3. Verify files are correct

## Still Not Working?

1. **Check browser console** (F12 → Console tab):
   - Look for 404 errors
   - Check for JavaScript errors

2. **Check network tab** (F12 → Network tab):
   - See if CSS/JS files are loading
   - Check response codes (should be 200)

3. **Verify repository structure**:
   - Assets folder should be at root: `/assets/`
   - Not in `_assets/` or elsewhere

4. **Check GitHub Pages settings**:
   - Repository → Settings → Pages
   - Source should be "Deploy from a branch"
   - Branch should be `main` or `master`
   - Folder should be `/ (root)`

## Testing Locally vs Deployed

To test if it's a deployment issue:

1. Build locally: `bundle exec jekyll build`
2. Serve the `_site` folder: `cd _site && python3 -m http.server 8000`
3. Visit `http://localhost:8000`
4. Compare with deployed version

If local `_site` works but deployed doesn't:
- Likely a GitHub Pages configuration issue
- Check Actions tab for build errors

