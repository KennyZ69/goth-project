/** @type {import('tailwindcss').Config} */
module.exports = {
  // content: [ "./**/*.html", "./**/*.templ", "./**/*.go", "./**/**/*.templ", "./**/**/*.html", "./**/**/*.go" ],
  // content: ["./layouts/**/.{templ,go,html}", "./layouts/*.{templ,go,html}" ],
  content: [
    "./layouts/**/*.{templ,go,html}",
    "./layouts/*.{templ,go,html}",
  ],
 safelist: [],
}