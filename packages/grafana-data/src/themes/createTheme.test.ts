import { createTheme } from './createTheme';

describe('createTheme', () => {
  it('create default theme', () => {
    const theme = createTheme();
    expect(theme).toMatchInlineSnapshot(`
      Object {
        "breakpoints": Object {
          "down": [Function],
          "keys": Array [
            "xs",
            "sm",
            "md",
            "lg",
            "xl",
            "xxl",
          ],
          "unit": "px",
          "up": [Function],
          "values": Object {
            "lg": 992,
            "md": 769,
            "sm": 544,
            "xl": 1200,
            "xs": 0,
            "xxl": 1440,
          },
        },
        "components": Object {
          "height": Object {
            "lg": 6,
            "md": 4,
            "sm": 3,
          },
          "panel": Object {
            "headerHeight": 4,
            "padding": 1,
          },
        },
        "isDark": true,
        "isLight": false,
        "name": "Dark",
        "palette": Object {
          "border0": "#202226",
          "border1": "#2c3235",
          "border2": "#464c54",
          "contrastThreshold": 3,
          "error": Object {
            "border": "#FF5286",
            "contrastText": "#fff",
            "main": "#D10E5C",
            "name": "error",
            "text": "#FF5286",
          },
          "formComponent": Object {
            "background": "#0b0c0e",
            "border": "#2c3235",
            "disabledBackground": "#141619",
            "disabledText": "rgba(255, 255, 255, 0.3)",
            "text": "rgba(255, 255, 255, 0.75)",
          },
          "getContrastText": [Function],
          "getHoverColor": [Function],
          "hoverFactor": 0.15,
          "info": Object {
            "border": "#5B93FF",
            "contrastText": "#fff",
            "main": "#3658E2",
            "name": "info",
            "text": "#5B93FF",
          },
          "layer0": "#0b0c0e",
          "layer1": "#141619",
          "layer2": "#202226",
          "mode": "dark",
          "primary": Object {
            "border": "#5B93FF",
            "contrastText": "#fff",
            "main": "#3658E2",
            "name": "primary",
            "text": "#5B93FF",
          },
          "secondary": Object {
            "border": "rgba(255,255,255,0.1)",
            "contrastText": "rgba(255, 255, 255, 0.8)",
            "main": "rgba(255,255,255,0.1)",
            "name": "secondary",
            "text": "rgba(255,255,255,0.1)",
          },
          "success": Object {
            "border": "#6CCF8E",
            "contrastText": "#fff",
            "main": "#13875D",
            "name": "success",
            "text": "#6CCF8E",
          },
          "text": Object {
            "disabled": "rgba(255, 255, 255, 0.3)",
            "link": "#5B93FF",
            "maxContrast": "#fff",
            "primary": "rgba(255, 255, 255, 0.75)",
            "secondary": "rgba(255, 255, 255, 0.50)",
          },
          "tonalOffset": 0.1,
          "warning": Object {
            "border": "#eb7b18",
            "contrastText": "#000",
            "main": "#eb7b18",
            "name": "warning",
            "text": "#eb7b18",
          },
        },
        "shadows": Object {
          "level1": "0px 2px 1px -1px rgba(0,0,0,0.2),0px 1px 1px 0px rgba(0,0,0,0.14),0px 1px 3px 0px rgba(0,0,0,0.12)",
          "level2": "0px 3px 1px -2px rgba(0,0,0,0.2),0px 2px 2px 0px rgba(0,0,0,0.14),0px 1px 5px 0px rgba(0,0,0,0.12)",
          "level3": "0px 3px 3px -2px rgba(0,0,0,0.2),0px 3px 4px 0px rgba(0,0,0,0.14),0px 1px 8px 0px rgba(0,0,0,0.12)",
        },
        "shape": Object {
          "borderRadius": [Function],
        },
        "spacing": [Function],
        "typography": Object {
          "body": Object {
            "fontFamily": "\\"Roboto\\", \\"Helvetica\\", \\"Arial\\", sans-serif",
            "fontSize": "1rem",
            "fontWeight": 400,
            "letterSpacing": "0.01071em",
            "lineHeight": 1.5,
          },
          "fontFamily": "\\"Roboto\\", \\"Helvetica\\", \\"Arial\\", sans-serif",
          "fontFamilyMonospace": "Menlo, Monaco, Consolas, 'Courier New', monospace",
          "fontSize": 14,
          "fontWeightBold": 700,
          "fontWeightLight": 300,
          "fontWeightMedium": 500,
          "fontWeightRegular": 400,
          "h1": Object {
            "fontFamily": "\\"Roboto\\", \\"Helvetica\\", \\"Arial\\", sans-serif",
            "fontSize": "2rem",
            "fontWeight": 300,
            "letterSpacing": "-0.05357em",
            "lineHeight": 1.167,
          },
          "h2": Object {
            "fontFamily": "\\"Roboto\\", \\"Helvetica\\", \\"Arial\\", sans-serif",
            "fontSize": "1.7142857142857142rem",
            "fontWeight": 300,
            "letterSpacing": "-0.02083em",
            "lineHeight": 1.2,
          },
          "h3": Object {
            "fontFamily": "\\"Roboto\\", \\"Helvetica\\", \\"Arial\\", sans-serif",
            "fontSize": "1.5rem",
            "fontWeight": 400,
            "letterSpacing": "0em",
            "lineHeight": 1.167,
          },
          "h4": Object {
            "fontFamily": "\\"Roboto\\", \\"Helvetica\\", \\"Arial\\", sans-serif",
            "fontSize": "1.2857142857142858rem",
            "fontWeight": 400,
            "letterSpacing": "0.01389em",
            "lineHeight": 1.235,
          },
          "h5": Object {
            "fontFamily": "\\"Roboto\\", \\"Helvetica\\", \\"Arial\\", sans-serif",
            "fontSize": "1.1428571428571428rem",
            "fontWeight": 400,
            "letterSpacing": "0em",
            "lineHeight": 1.334,
          },
          "h6": Object {
            "fontFamily": "\\"Roboto\\", \\"Helvetica\\", \\"Arial\\", sans-serif",
            "fontSize": "1rem",
            "fontWeight": 500,
            "letterSpacing": "0.01071em",
            "lineHeight": 1.6,
          },
          "htmlFontSize": 14,
          "pxToRem": [Function],
          "size": Object {
            "base": "14px",
            "lg": "18px",
            "md": "14px",
            "sm": "12px",
            "xs": "10px",
          },
        },
        "zIndex": Object {
          "dropdown": 1030,
          "modal": 1060,
          "modalBackdrop": 1050,
          "navbarFixed": 1000,
          "sidemenu": 1020,
          "tooltip": 1040,
          "typeahead": 1030,
        },
      }
    `);
  });
});
