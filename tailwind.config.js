/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [ "./**/*.html", "./**/*.templ", "./**/*.go", ],
  plugins: [require("daisyui"), require("@tailwindcss/typography")],
  safelist:[],
  daisyui: {
    themes: ["retro"]
  }
}

