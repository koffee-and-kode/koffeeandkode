// ============================================
// KoffeeAndKode - Main JavaScript
// ============================================

(function() {
  'use strict';

  // Theme Toggle
  function initThemeToggle() {
    const themeToggle = document.getElementById('theme-toggle');
    const themeIcon = themeToggle?.querySelector('.theme-icon');
    
    // Get saved theme or default to light
    const savedTheme = localStorage.getItem('theme') || 'light';
    document.documentElement.setAttribute('data-theme', savedTheme);
    
    if (themeIcon) {
      themeIcon.textContent = savedTheme === 'dark' ? '☀️' : '🌙';
    }
    
    if (themeToggle) {
      themeToggle.addEventListener('click', () => {
        const currentTheme = document.documentElement.getAttribute('data-theme');
        const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
        
        document.documentElement.setAttribute('data-theme', newTheme);
        localStorage.setItem('theme', newTheme);
        
        if (themeIcon) {
          themeIcon.textContent = newTheme === 'dark' ? '☀️' : '🌙';
        }
      });
    }
  }

  // Copy Code Block Button
  function initCopyButtons() {
    const codeBlocks = document.querySelectorAll('pre code');
    
    codeBlocks.forEach((codeBlock) => {
      const pre = codeBlock.parentElement;
      
      // Skip if already has a copy button
      if (pre.querySelector('.copy-button')) {
        return;
      }
      
      const copyButton = document.createElement('button');
      copyButton.className = 'copy-button';
      copyButton.textContent = 'Copy';
      copyButton.setAttribute('aria-label', 'Copy code to clipboard');
      
      copyButton.addEventListener('click', async () => {
        const text = codeBlock.textContent;
        
        try {
          await navigator.clipboard.writeText(text);
          copyButton.textContent = 'Copied!';
          copyButton.style.backgroundColor = '#28a745';
          
          setTimeout(() => {
            copyButton.textContent = 'Copy';
            copyButton.style.backgroundColor = '';
          }, 2000);
        } catch (err) {
          console.error('Failed to copy text:', err);
          copyButton.textContent = 'Error';
          
          setTimeout(() => {
            copyButton.textContent = 'Copy';
          }, 2000);
        }
      });
      
      pre.style.position = 'relative';
      pre.appendChild(copyButton);
    });
  }

  // Mobile Navigation Toggle
  function initMobileNav() {
    const navToggle = document.getElementById('nav-toggle');
    const navMenu = document.querySelector('.nav-menu');
    
    if (navToggle && navMenu) {
      // Close menu when clicking outside
      document.addEventListener('click', (e) => {
        if (!navMenu.contains(e.target) && !navToggle.contains(e.target)) {
          navToggle.checked = false;
        }
      });
      
      // Close menu when clicking a link
      const navLinks = navMenu.querySelectorAll('a');
      navLinks.forEach((link) => {
        link.addEventListener('click', () => {
          navToggle.checked = false;
        });
      });
    }
  }

  // Smooth Scroll for Anchor Links
  function initSmoothScroll() {
    document.querySelectorAll('a[href^="#"]').forEach((anchor) => {
      anchor.addEventListener('click', function (e) {
        const href = this.getAttribute('href');
        if (href === '#') return;
        
        const target = document.querySelector(href);
        if (target) {
          e.preventDefault();
          target.scrollIntoView({
            behavior: 'smooth',
            block: 'start'
          });
        }
      });
    });
  }

  // Initialize all features when DOM is ready
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => {
      initThemeToggle();
      initCopyButtons();
      initMobileNav();
      initSmoothScroll();
    });
  } else {
    initThemeToggle();
    initCopyButtons();
    initMobileNav();
    initSmoothScroll();
  }
})();

