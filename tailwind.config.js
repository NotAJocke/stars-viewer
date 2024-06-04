/** @type {import('tailwindcss').Config} */

const defaultTheme = require("tailwindcss/defaultTheme")

module.exports = {
  content: [
    "internal/templates/**/*.templ"
  ],
  theme: {
    extend: {
      fontFamily: {
        'sans': ['"Helvetica Neue"', '"Helvetica"', ...defaultTheme.fontFamily.sans],
      }
    },
  },
  plugins: [],
}

