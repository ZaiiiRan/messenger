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
        'full_screen': '100vw',
        '7/20': '35%'
      },
      backgroundColor: {
        'dark-brand-color': '#001A83',
        'light-brand-color': '#879FFF'
      },
      screens: {
        'sm': '500px',
        'mobile': '200px',
        '2k': '1921px',
        '4k': '2561px',
      },
    },
  },
  plugins: [],
}

