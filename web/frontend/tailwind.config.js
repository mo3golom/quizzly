import animations from '@midudev/tailwind-animations'

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./templ/**/*.{html,js,templ}", "./safelist.txt"],
  daisyui: {
    themes: ["light"],
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

