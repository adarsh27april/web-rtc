const fs = require('fs');
const path = require('path');

// Configuration - Make these easily configurable
const PROJECT_DIR = './signaling-server'; // Change this to extract from any project folder
const OUTPUT_FILE = './extracted-code-go.md'; // Generic output filename
const FILE_EXTENSIONS = ['.go', '.js', '.json','.ts', '.html']; // Add/remove extensions as needed
const EXCLUDE_PATTERNS = [ // Patterns to exclude (files and folders)
    '*_debug_*', 'debug', 'node_modules', '.git',
    'client/client.go', 'hub/hub.go', 'utils/message.go',
    'package-lock.json', 'package.json'
];
const INCLUDE_HIDDEN_FILES = false; // Set to true to include files starting with '.'

// Function to check if a name or path matches any excluded patterns
function matchesExcludePattern(name, relativePath = '') {
    for (const pattern of EXCLUDE_PATTERNS) {
        // Convert wildcard pattern to regex
        const regexPattern = pattern
            .replace(/\*/g, '.*')  // * becomes .*
            .replace(/\?/g, '.')   // ? becomes .
            .replace(/[.+^${}()|[\]\\]/g, '\\$&'); // Escape regex special chars

        try {
            const regex = new RegExp(regexPattern);
            // Check both the filename and the relative path
            if (regex.test(name) || regex.test(relativePath)) {
                return true;
            }
        } catch (error) {
            // If regex is invalid, fall back to simple string matching
            if (name.includes(pattern.replace(/\*/g, '')) || relativePath.includes(pattern.replace(/\*/g, ''))) {
                return true;
            }
        }
    }
    return false;
}

// Function to recursively find all files with specified extensions in a directory
function findFiles(dir, fileList = []) {
    const files = fs.readdirSync(dir);

    files.forEach(file => {
        const filePath = path.join(dir, file);
        const stat = fs.statSync(filePath);

        if (stat.isDirectory()) {
            // Skip excluded directories
            const relativePath = path.relative(PROJECT_DIR, filePath);
            if (!matchesExcludePattern(file, relativePath)) {
                // Skip hidden directories unless explicitly allowed
                if (INCLUDE_HIDDEN_FILES || !file.startsWith('.')) {
                    findFiles(filePath, fileList);
                }
            }
        } else if (shouldIncludeFile(file, path.relative(PROJECT_DIR, filePath))) {
            // Add file to the list if it matches our criteria
            fileList.push({
                path: filePath,
                relativePath: path.relative(PROJECT_DIR, filePath),
                name: file,
                extension: path.extname(file)
            });
        }
    });

    return fileList;
}

// Function to check if a file should be included based on extension
function shouldIncludeFile(filename, relativePath = '') {
    const extension = path.extname(filename);

    // Check if filename or path matches any excluded patterns
    if (matchesExcludePattern(filename, relativePath)) {
        return false;
    }

    // Include files with specified extensions
    if (FILE_EXTENSIONS.includes(extension)) {
        return true;
    }

    // Include files without extension if they're not hidden (unless allowed)
    if (extension === '' && (INCLUDE_HIDDEN_FILES || !filename.startsWith('.'))) {
        return true;
    }

    return false;
}

// Function to read file content with encoding detection
function readFileContent(filePath) {
    try {
        // Try to read as UTF-8 first
        const content = fs.readFileSync(filePath, 'utf8');

        // Check if content contains null bytes (binary file indicator)
        if (content.includes('\x00')) {
            console.warn(`⚠️  Warning: ${filePath} appears to be binary, skipping...`);
            return null;
        }

        return content;
    } catch (error) {
        // If UTF-8 fails, try to detect if it's a binary file
        try {
            const buffer = fs.readFileSync(filePath);
            // Check for null bytes in the first 1024 bytes
            const sample = buffer.slice(0, 1024);
            if (sample.includes(0)) {
                console.warn(`⚠️  Warning: ${filePath} appears to be binary, skipping...`);
                return null;
            }
            // If no null bytes, try to convert to string
            return buffer.toString('utf8');
        } catch (bufferError) {
            console.error(`Error reading file ${filePath}:`, error.message);
            return null;
        }
    }
}

// Function to get language identifier for syntax highlighting (use extension directly)
function getLanguageIdentifier(extension) {
    // Remove the dot and return the extension as the language identifier
    // Most markdown renderers will recognize common extensions
    return extension ? extension.substring(1) : 'text';
}

// Function to generate the combined code file
function generateCombinedCode(files) {
    let content = `# Extracted Code from Project\n\n`;
    content += `**Project Directory:** \`${PROJECT_DIR}\`\n\n`;
    content += `**Generated on:** ${new Date().toISOString()}\n\n`;
    content += `**File Extensions:** ${FILE_EXTENSIONS.join(', ')}\n\n`;
    content += `This file contains all source code from the specified project folder.\n\n`;
    content += `---\n\n`;

    // Sort files by their relative path for consistent ordering
    files.sort((a, b) => a.relativePath.localeCompare(b.relativePath));

    files.forEach((file, index) => {
        const fileContent = readFileContent(file.path);

        if (fileContent !== null) {
            const language = getLanguageIdentifier(file.extension);
            content += `## ${index + 1}. ${file.relativePath}\n\n`;
            content += `**File:** \`${file.relativePath}\`\n`;
            content += `**Type:** \`${file.extension || 'no extension'}\`\n\n`;
            content += `\`\`\`${language}\n`;
            content += fileContent;
            content += `\n\`\`\`\n\n`;
            content += `---\n\n`;
        }
    });

    content += `\n## Summary\n\n`;
    content += `Total files processed: ${files.length}\n\n`;

    // Group files by extension
    const extensionStats = {};
    files.forEach(file => {
        const ext = file.extension || 'no extension';
        extensionStats[ext] = (extensionStats[ext] || 0) + 1;
    });

    content += `**Files by extension:**\n`;
    Object.entries(extensionStats)
        .sort(([, a], [, b]) => b - a)
        .forEach(([ext, count]) => {
            content += `- \`${ext}\`: ${count} files\n`;
        });

    content += `\n**All files with relative paths:**\n`;
    files.forEach((file, index) => {
        content += `${index + 1}. **${file.name}** - \`${file.relativePath}\`\n`;
    });

    return content;
}

// Function to write the combined code to output file
function writeOutputFile(content) {
    try {
        // Ensure content has proper line endings for the current platform
        const normalizedContent = content.replace(/\r\n/g, '\n').replace(/\r/g, '\n');

        // Write with explicit UTF-8 encoding and proper line endings
        fs.writeFileSync(OUTPUT_FILE, normalizedContent, {
            encoding: 'utf8',
            flag: 'w'
        });

        console.log(`✅ Successfully wrote combined code to: ${OUTPUT_FILE}`);

        // Verify the file was created and is readable
        const stats = fs.statSync(OUTPUT_FILE);
        console.log(`📊 File size: ${stats.size} bytes`);

        // Test reading the first few characters to ensure it's accessible
        const testRead = fs.readFileSync(OUTPUT_FILE, 'utf8').substring(0, 50);
        console.log(`📖 File preview: ${testRead}...`);

    } catch (error) {
        console.error(`❌ Error writing output file:`, error.message);
    }
}

// Main execution function
function main() {
    console.log('🔍 Searching for files in project folder...');
    console.log(`📁 Project directory: ${PROJECT_DIR}`);
    console.log(`📄 File extensions: ${FILE_EXTENSIONS.join(', ')}`);
    console.log(`🚫 Excluded patterns: ${EXCLUDE_PATTERNS.join(', ')}`);

    // Check if project directory exists
    if (!fs.existsSync(PROJECT_DIR)) {
        console.error(`❌ Directory '${PROJECT_DIR}' not found!`);
        console.log('Please update the PROJECT_DIR constant to point to a valid directory.');
        return;
    }

    // Find all matching files
    const files = findFiles(PROJECT_DIR);

    if (files.length === 0) {
        console.log('❌ No matching files found in project folder.');
        console.log(`Try updating FILE_EXTENSIONS constant or check if ${PROJECT_DIR} contains files with these extensions.`);
        return;
    }

    console.log(`📁 Found ${files.length} matching files:`);
    files.forEach((file, index) => {
        console.log(`   ${index + 1}. ${file.relativePath} (${file.extension || 'no extension'})`);
    });

    // Generate combined code
    console.log('\n📝 Generating combined code...');
    const combinedCode = generateCombinedCode(files);

    // Write to output file
    console.log(`💾 Writing to ${OUTPUT_FILE}...`);
    writeOutputFile(combinedCode);

    console.log('\n🎉 Extraction complete!');
    console.log(`📄 Output file: ${OUTPUT_FILE}`);
    console.log(`📊 Total files processed: ${files.length}`);
}

// Run the script
if (require.main === module) {
    main();
}

module.exports = {
    findFiles,
    readFileContent,
    generateCombinedCode,
    writeOutputFile,
    getLanguageIdentifier,
    shouldIncludeFile
};
