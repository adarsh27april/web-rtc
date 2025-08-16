/**
 * Theme Management
 */
type Theme = 'light' | 'dark'

let currentTheme: Theme = 'light'

export function getSystemTheme(): Theme {
   return window.matchMedia("(prefers-color-scheme: dark)").matches
      ? "dark"
      : "light";
}

export function updateThemeIcon(): void {
   const themeIcon = document.getElementById("themeIcon");
   if (themeIcon) {
      themeIcon.textContent = currentTheme === "dark" ? "☀️" : "🌙";
   }
}

export function applyTheme(theme: Theme): void {
   const body = document.body;
   body.classList.remove("theme-light", "theme-dark");
   body.classList.add(`theme-${theme}`);

   // Force CSS reflow
   body.style.display = "none";
   void body.offsetHeight;
   body.style.display = "";

   currentTheme = theme;
   updateThemeIcon();
   console.log(`Applied theme: ${theme}, Body classes:`, body.className);
}

export function toggleTheme(): void {
   const nextTheme = currentTheme === "light" ? "dark" : "light";
   applyTheme(nextTheme);

   showToast(`Switched to ${nextTheme === "light" ? "Light Theme" : "Dark Theme"}`, "info");
   console.log(`Theme changed to: ${nextTheme}`);
}


/**
 * UI Enhancement
 */

export function addButtonFeedback(): void {
   const buttons = document.querySelectorAll<HTMLButtonElement>("button");
   buttons.forEach((button) => {
      button.addEventListener("click", function () {
         this.classList.add("loading");
         setTimeout(() => this.classList.remove("loading"), 500);
      });
   });
}


export function setupTextareaAutoResize(): void {
   const textareas = document.querySelectorAll<HTMLTextAreaElement>("textarea");
   textareas.forEach((textarea) => {
      textarea.addEventListener("input", function () {
         this.style.height = "auto";
         this.style.height = `${Math.max(120, this.scrollHeight)}px`;
      });
   });
}

/**
 * Toast Notifications
 */
export function showToast(message: string, type: "info" | "success" | "error" = "info"): void {
   const toast = document.createElement("div");
   toast.style.cssText = `
    position: fixed;
    top: 80px;
    right: 20px;
    background: ${type === "success" ? "#4CAF50" : type === "error" ? "#F44336" : "#2196F3"};
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

   setTimeout(() => (toast.style.transform = "translateX(0)"), 10);
   setTimeout(() => {
      toast.style.transform = "translateX(100%)";
      setTimeout(() => toast.remove(), 300);
   }, 3000);
}


export function initUI(): void {
   const systemTheme = getSystemTheme();
   applyTheme(systemTheme);
   addButtonFeedback();
   setupTextareaAutoResize();

   const themeToggleBtn = document.querySelector<HTMLButtonElement>(".theme-toggle");
   themeToggleBtn?.addEventListener("click", toggleTheme);
}

// export function initThemeSwitcher(toggleButton: HTMLButtonElement): void {
//    currentTheme = getSystemTheme()
//    applyTheme(currentTheme)

//    toggleButton.addEventListener('click', () => {
//       currentTheme = currentTheme === 'light' ? 'dark' : 'light'
//       applyTheme(currentTheme)
//    })
// }

// function applyTheme(theme: Theme): void {
//    document.documentElement.setAttribute('data-theme', theme)
//    const icon = document.querySelector('.theme-toggle .icon')
//    if (icon) {
//       icon.textContent = theme === 'dark' ? '☀️' : '🌙'
//    }
// }


// // =====================
// // Toast Notifications
// // =====================

// type ToastType = 'info' | 'success' | 'error'

// export function showToast(message: string, type: ToastType = 'info'): void {
//    const toast = document.createElement('div')
//    toast.className = `toast toast--${type}`
//    toast.textContent = message
//    toast.setAttribute('role', 'alert')

//    document.body.appendChild(toast)

//    setTimeout(() => {
//       toast.classList.add('toast--hide')
//       toast.addEventListener('transitionend', () => toast.remove())
//    }, 3000)
// }

// // =====================
// // UI Enhancements
// // =====================

// export function setupButtonFeedback(): void {
//    document.querySelectorAll('button').forEach(button => {
//       button.addEventListener('click', function (this: HTMLButtonElement) {
//          this.classList.add('loading')
//          setTimeout(() => this.classList.remove('loading'), 200)
//       })
//    })
// }

// export function setupTextareaAutoresize(textarea?: HTMLTextAreaElement): void {
//    const target = textarea || document.getElementById('message-area') as HTMLTextAreaElement
//    target?.addEventListener('input', function () {
//       this.style.height = 'auto'
//       this.style.height = `${this.scrollHeight}px`
//    })
// }

// // =====================
// // Form Utilities
// // =====================

// export function validateFormInput(input: HTMLInputElement): boolean {
//    const value = input.value.trim()
//    if (!value) {
//       input.classList.add('error')
//       showToast('Please fill in this field', 'error')
//       return false
//    }
//    input.classList.remove('error')
//    return true
// }

// // =====================
// // DOM Helpers
// // =====================

// export function getElement<T extends HTMLElement>(selector: string): T {
//    const el = document.querySelector(selector)
//    if (!el) throw new Error(`Element not found: ${selector}`)
//    return el as T
// }