module.exports = {
    apps: [
        {
            name: 'email-specter-frontend',
            script: 'node_modules/.bin/next',
            args: 'start',
            env: {
                NODE_ENV: 'production',
            },
            instances: 'max',
            autorestart: true,
            watch: true, // Watch for file changes and restart the app
            max_memory_restart: '1G',
        },
    ],
};