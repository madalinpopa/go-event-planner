/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["ui/**/*.{tmpl,js}"],
  theme: {
    extend: {
      fontFamily: {
        pally: ['Pally-Variable']
      }
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ]
}

