/** @type {import('tailwindcss').Config} */
module.exports = {
  // content: [ "./**/*.html", "./**/*.templ", "./**/*.go", "./**/**/*.templ", "./**/**/*.html", "./**/**/*.go" ],
  // content: ["./layouts/**/.{templ,go,html}", "./layouts/*.{templ,go,html}" ],
  content: [
    "./layouts/**/*.{templ,go,html}",
    "./layouts/*.{templ,go,html}",
  ],
  safelist: [],
  plugins: [
    function({ addUtilities }) {
      const newUtilities = {
        '.vertical-text': {
          writingMode: 'vertical-rl',
          transform: 'rotate(180deg)',
          textAlign: 'center',
          whiteSpace: 'nowrap',
        },
        '.horizontal-text': {
          writingMode: 'horizontal-tb',
          transform: 'none',
          textAlign: 'right',
        }
      }
      addUtilities(newUtilities, ['responsive', 'hover'])
    }
  ],
}
