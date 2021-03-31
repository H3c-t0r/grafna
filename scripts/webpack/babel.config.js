module.exports = function getBabelConfig() {
  return {
    cacheDirectory: true,
    babelrc: false,
    // Note: order is top-to-bottom and/or left-to-right
    // Note: order is bottom-to-top and/or right-to-left
    presets: [
      [
        '@babel/preset-env',
        {
          targets: {
            browsers: 'last 3 versions',
          },
          useBuiltIns: 'entry',
          corejs: 3,
          modules: false,
        },
      ],
      [
        '@babel/preset-typescript',
        {
          allowNamespaces: true,
          allowDeclareFields: true,
        },
      ],
      '@babel/preset-react',
    ],
    plugins: [
      // [
      //   require('@rtsao/plugin-proposal-class-properties'),
      //   {
      //     loose: true,
      //   },
      // ],
      ['@babel/plugin-proposal-class-properties', { loose: true }],
      // ['@babel/plugin-proposal-private-methods', { loose: true }],
      // ['@babel/plugin-proposal-object-rest-spread', { loose: true }],
      '@babel/plugin-proposal-nullish-coalescing-operator',
      '@babel/plugin-proposal-optional-chaining',
      '@babel/plugin-syntax-dynamic-import', // needed for `() => import()` in routes.ts
      'angularjs-annotate',
    ],
  };
};
