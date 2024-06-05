/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [ "./**/*.html", "./**/*.templ", "./**/*.go", ],
  plugins: [require("daisyui")],
  safelist:[],
  daisyui: {
    themes: ["retro"]
  }
}

