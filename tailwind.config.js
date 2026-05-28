module.exports = {
    mode: 'jit',
    darkMode: 'selector',
    content: ["public/**/*.html"],
    theme: {
      container: {
        center: true,
        padding: '1rem'
      },
      borderWidth: {
        '1': '1px'
      },
      extend: {},
    },
    plugins: [
      function ({ addComponents }) {
        addComponents({
          '.container': {
            flexGrow: '1',
            display: 'flex',
            flexDirection: 'column',
            '@screen lg': {
              maxWidth: '1280px',
            },
            '@screen xl': {
              maxWidth: '1280px',
            },
          }
        })
      }
    ],
  }

