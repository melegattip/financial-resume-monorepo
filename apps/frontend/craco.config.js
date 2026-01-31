const path = require('path');

module.exports = {
  webpack: {
    configure: (webpackConfig, { env }) => {
      // Remove console logs in production build
      if (env === 'production') {
        // Add Terser plugin configuration to remove console logs
        const TerserPlugin = require('terser-webpack-plugin');
        
        // Find the existing TerserPlugin or create a new one
        const existingTerserPlugin = webpackConfig.optimization.minimizer.find(
          plugin => plugin.constructor.name === 'TerserPlugin'
        );

        if (existingTerserPlugin) {
          // Update existing TerserPlugin options
          const currentOptions = existingTerserPlugin.options || {};
          const currentTerserOptions = currentOptions.terserOptions || {};
          const currentCompress = currentTerserOptions.compress || {};
          
          existingTerserPlugin.options.terserOptions = {
            ...currentTerserOptions,
            compress: {
              ...currentCompress,
              drop_console: true,
              drop_debugger: true,
              pure_funcs: ['console.log', 'console.info', 'console.debug', 'console.warn']
            }
          };
        } else {
          // Add new TerserPlugin if it doesn't exist
          webpackConfig.optimization.minimizer.push(
            new TerserPlugin({
              terserOptions: {
                compress: {
                  drop_console: true,
                  drop_debugger: true,
                  pure_funcs: ['console.log', 'console.info', 'console.debug', 'console.warn']
                }
              }
            })
          );
        }
      }

      return webpackConfig;
    }
  },
  
  // Configure babel to remove console statements in production
  babel: {
    plugins: [
      ...(process.env.NODE_ENV === 'production' ? [
        ['transform-remove-console', { 
          exclude: ['error', 'warn'] // Keep error and warn logs
        }]
      ] : [])
    ]
  }
};
