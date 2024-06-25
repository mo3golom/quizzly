/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./templ/**/*.{html,js,templ}", "./safelist.txt"],
  daisyui: {
    themes: ["light"],
  },
  plugins: [require("daisyui")],
  theme: {
    container: {
      center: true,
    },
  },
}

