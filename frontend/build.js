const fs = require('fs');
const path = require('path');

// Helper function to copy directory recursively
function copyDir(src, dest) {
    // Create destination directory
    if (!fs.existsSync(dest)) {
        fs.mkdirSync(dest, { recursive: true });
    }

    // Read source directory
    const entries = fs.readdirSync(src, { withFileTypes: true });

    for (let entry of entries) {
        const srcPath = path.join(src, entry.name);
        const destPath = path.join(dest, entry.name);

        if (entry.isDirectory()) {
            // Skip node_modules and dist
            if (entry.name === 'node_modules' || entry.name === 'dist') {
                continue;
            }
            copyDir(srcPath, destPath);
        } else {
            fs.copyFileSync(srcPath, destPath);
        }
    }
}

// Clean dist directory
const distDir = path.join(__dirname, 'dist');
if (fs.existsSync(distDir)) {
    fs.rmSync(distDir, { recursive: true, force: true });
}

// Create dist directory
fs.mkdirSync(distDir, { recursive: true });

// Copy index.html
fs.copyFileSync(
    path.join(__dirname, 'index.html'),
    path.join(distDir, 'index.html')
);

// Copy src directory
copyDir(
    path.join(__dirname, 'src'),
    path.join(distDir, 'src')
);

// Copy wailsjs directory
copyDir(
    path.join(__dirname, 'wailsjs'),
    path.join(distDir, 'wailsjs')
);

console.log('Build complete! Files copied to dist/');
