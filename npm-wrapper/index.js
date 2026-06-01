const { execSync } = require('child_process');
const path = require('path');
const fs = require('fs');

function register() {
    const platform = process.platform;
    let binaryName;

    switch (platform) {
        case 'win32':
            binaryName = 'flowmeter-host-win.exe';
            break;
        case 'darwin':
            binaryName = 'flowmeter-host-macos';
            break;
        case 'linux':
            binaryName = 'flowmeter-host-linux';
            break;
        default:
            console.error(`Unsupported platform: ${platform}`);
            process.exit(1);
    }

    const binaryPath = path.join(__dirname, 'bin', binaryName);

    if (!fs.existsSync(binaryPath)) {
        console.error(`Binary not found at: ${binaryPath}`);
        process.exit(1);
    }

    try {
        console.log(`Registering FlowMeter host for ${platform}...`);
        execSync(`${binaryPath} --register`, { stdio: 'inherit' });
        console.log('Registration successful.');
    } catch (error) {
        console.error('Registration failed:', error.message);
        process.exit(1);
    }
}

if (require.main === module) {
    register();
}

module.exports = { register };
