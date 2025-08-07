{/* <script> */ }
// Theme Management - Simple Light/Dark Toggle
let currentTheme = 'auto'; // Will be set to 'light' or 'dark' on init

function getSystemTheme() {
   return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
}

function updateThemeIcon() {
   const themeIcon = document.getElementById('themeIcon');

   if (currentTheme === 'dark') {
      themeIcon.textContent = '☀️'; // Show sun icon when in dark mode
   } else {
      themeIcon.textContent = '🌙'; // Show moon icon when in light mode
   }
}

function applyTheme(theme) {
   const body = document.body;

   // Remove existing theme classes
   body.classList.remove('theme-light', 'theme-dark');

   // Always apply a specific theme class (no auto mode)
   body.classList.add(`theme-${theme}`);

   // Force CSS recalculation by temporarily changing a property
   body.style.display = 'none';
   body.offsetHeight; // Trigger reflow
   body.style.display = '';

   updateThemeIcon();

   // Store preference in memory
   currentTheme = theme;

   // Log the current theme for debugging
   console.log(`Applied theme: ${theme}, Body classes:`, body.className);
}

function toggleTheme() {
   // Simple toggle between light and dark
   const nextTheme = currentTheme === 'light' ? 'dark' : 'light';

   applyTheme(nextTheme);

   // Show feedback to user
   const themeNames = {
      'light': 'Light Theme',
      'dark': 'Dark Theme'
   };

   showToast(`Switched to ${themeNames[nextTheme]}`, 'info');

   console.log(`Theme changed to: ${nextTheme}`);
}

// Platform Detection
function detectPlatform() {
   const userAgent = navigator.userAgent.toLowerCase();
   const platform = navigator.platform.toLowerCase();

   // Check for iOS devices
   if (/ipad|iphone|ipod/.test(userAgent) || (platform === 'macintel' && navigator.maxTouchPoints > 1)) {
      return 'ios';
   }

   // Check for macOS
   if (platform.includes('mac')) {
      return 'ios'; // Use iOS styling for macOS as well
   }

   // Check for Android
   if (/android/.test(userAgent)) {
      return 'android';
   }

   // Default to Android Material Design for other platforms
   return 'android';
}

// Apply Platform-Specific Styling
function applyPlatformStyling() {
   const platform = detectPlatform();
   const body = document.body;

   // Remove any existing platform classes
   body.classList.remove('platform-ios', 'platform-android');

   // Add appropriate platform class
   body.classList.add(`platform-${platform}`);

   console.log(`Detected platform: ${platform}`);
}

// Enhanced Button Feedback
function addButtonFeedback() {
   const buttons = document.querySelectorAll('button');

   buttons.forEach(button => {
      button.addEventListener('click', function (e) {
         // Add loading state
         this.classList.add('loading');

         // Remove loading state after a short delay
         setTimeout(() => {
            this.classList.remove('loading');
         }, 500);
      });
   });
}

// Enhanced Textarea Auto-resize
function setupTextareaAutoResize() {
   const textareas = document.querySelectorAll('textarea');

   textareas.forEach(textarea => {
      textarea.addEventListener('input', function () {
         this.style.height = 'auto';
         this.style.height = Math.max(120, this.scrollHeight) + 'px';
      });
   });
}

// Toast Notification System
function showToast(message, type = 'info') {
   const toast = document.createElement('div');
   toast.style.cssText = `
                position: fixed;
                top: 80px;
                right: 20px;
                background: ${type === 'success' ? '#4CAF50' : type === 'error' ? '#F44336' : '#2196F3'};
                color: white;
                padding: 12px 24px;
                border-radius: var(--border-radius);
                font-family: var(--font-family);
                font-weight: 500;
                font-size: 0.9rem;
                z-index: 1000;
                box-shadow: var(--shadow);
                transform: translateX(100%);
                transition: transform 0.3s ease;
                max-width: 250px;
                word-wrap: break-word;
            `;
   toast.textContent = message;

   document.body.appendChild(toast);

   // Trigger animation
   setTimeout(() => {
      toast.style.transform = 'translateX(0)';
   }, 10);

   // Remove toast
   setTimeout(() => {
      toast.style.transform = 'translateX(100%)';
      setTimeout(() => {
         if (document.body.contains(toast)) {
            document.body.removeChild(toast);
         }
      }, 300);
   }, 3000);
}

// Initialize on DOM Content Loaded
document.addEventListener('DOMContentLoaded', function () {
   applyPlatformStyling();

   // Detect system theme and set it as initial theme
   const systemTheme = getSystemTheme();
   applyTheme(systemTheme);

   addButtonFeedback();
   setupTextareaAutoResize();

   // Log platform detection
   const logElement = document.getElementById('log');
   logElement.textContent = `Platform detected: ${detectPlatform().toUpperCase()}\nInitial theme: ${systemTheme} (detected from system)\nReady for WebRTC connection...\n`;
});

// Keep original function calls as requested
function setRemote() {
   showToast('Setting remote description...', 'info');
   // Your existing setRemote logic here
   console.log('setRemote called');
}

function copySDP() {
   const sdpOutput = document.getElementById('sdpOutput');
   if (sdpOutput.value.trim()) {
      navigator.clipboard.writeText(sdpOutput.value).then(() => {
         showToast('SDP copied to clipboard!', 'success');
      }).catch(() => {
         showToast('Failed to copy SDP', 'error');
      });
   } else {
      showToast('No SDP to copy', 'error');
   }
}

function showSDP() {
   showToast('Displaying SDP...', 'info');
   // Your existing showSDP logic here
   console.log('showSDP called');
}
{/* </script> */ }