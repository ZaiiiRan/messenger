/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      height: {
        'full_screen': '100vh'
      },
      width: {
        'full_screen': '100vw'
      },
      backgroundColor: {
        'dark-brand-color': '#001A83',
        'light-brand-color': '#879FFF'
      },
      screens: {
        'mobile': '200px',
        '2k': '1921px',
        '4k': '2561px'
      },
    },
  },
  plugins: [],
}

