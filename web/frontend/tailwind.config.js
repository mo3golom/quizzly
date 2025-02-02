import animations from '@midudev/tailwind-animations'

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./templ/**/*.{html,js,templ}", "./safelist.txt"],
  daisyui: {
    themes: [
      {
        quizzly: {
          "primary": "#193277",
          "primary-content": "#eff6ff",
          "secondary": "#eff6ff",
          "secondary-content": "#193277",
          "accent": "#3b82f6",
          "accent-content": "#ffffff",
          "neutral": "#ffffff",
          "neutral-content": "#0f172a",
          "base-100": "#f1f5f9",
          "base-200": "#e2e8f0",
          "base-300": "#172554",
          "base-content": "#0f172a",
          "info": "#0ea5e9",
          "info-content": "#ffffff",
          "success": "#22c55e",
          "success-content": "#ffffff",
          "warning": "#f59e0b",
          "warning-content": "#ffffff",
          "error": "#f43f5d",
          "error-content": "#ffffff",
        },
      }
    ]
  },
  plugins: [
    require("daisyui"),
    animations
  ],
  theme: {
    container: {
      center: true,
    },
  },
}

